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

// SendClientSpecs Report Client information
func (h *Handler) SendClientSpecs() (id uint, err error) {
	ClientSpecs, err := h.Services.Information.LoadClientSpecs()
	if err != nil {
		return 0, err
	}
	body, err := json.Marshal(ClientSpecs)
	if err != nil {
		return 0, err
	}
	res, err := h.Gateway.NewRequest(http.MethodPost, "/client", body)
	if err != nil {
		return 0, err
	}
	if res.Code != 0 {
		return 0, fmt.Errorf("error with status code: %d, msg: %s", res.Code, res.Msg)
	}
	return uint(math.Ceil(res.Data.(float64))), nil
}

func (h *Handler) ServerIsAvailable() error {
	res, err := h.Gateway.NewRequest(http.MethodGet, fmt.Sprintf("/client/%d/health", uint64(h.ClientID)), nil)
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
		conn, err := ws.NewConnection(h.Configuration, fmt.Sprintf("/client/%d/ws", h.ClientID))
		if err != nil {
			h.Log("[!] Error Reconnect on WS: ", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		h.Connection = conn
		h.Connected = true
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
			ClientSpecs, err := h.Services.Information.LoadClientSpecs()
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
			} else {
				response = encode.StringToByte(encode.PrettyJson(ClientSpecs))
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
			} else {
				// 下载文件
				err = h.Services.Download.DownloadFile(m["filename"], m["path"])

				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
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

			//删除服务器文件
			h.Gateway.NewRequest(http.MethodDelete, "/file/"+filename, nil)

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
				break
			}

			// 等待一段时间以确保新进程已经启动
			time.Sleep(1 * time.Second)

			// 确保新进程已经启动后再退出当前进程
			os.Exit(0)
		case "exit":
			os.Exit(0)
		case "explorer":
			p := request.Parameter
			var fileExplorerQuery entitie.FileExplorerQuery
			err := json.Unmarshal([]byte(p), &fileExplorerQuery)
			if err != nil {
				hasError = true
				response = encode.StringToByte(err.Error())
				break
			}
			op := fileExplorerQuery.Op
			path := fileExplorerQuery.Path
			switch op {
			case "list":
				res, err := h.Services.FileExplorer.GetFileExplorer(path)
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
				response = encode.StringToByte(encode.PrettyJson(res))
			case "cat":
				res, err := h.Services.FileExplorer.ReadFile(path)
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
				response = res
			case "rename":
				newFilename := fileExplorerQuery.Filename
				err := h.Services.FileExplorer.Rename(path, newFilename)
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
			case "remove":
				err := h.Services.FileExplorer.Remove(path)
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
			case "edit":
				fileContent := fileExplorerQuery.FileContent
				err := h.Services.FileExplorer.WriteFile(path, []byte(fileContent))
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
			case "mkdir":
				err := h.Services.FileExplorer.MkDir(path)
				if err != nil {
					hasError = true
					response = encode.StringToByte(err.Error())
				}
			}
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
