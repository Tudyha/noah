package handler

import (
	"encoding/json"
	"fmt"
	"github.com/creack/pty"
	"io"
	"math"
	"net"
	"net/http"
	"noah/client/app/environment"
	"noah/client/app/service"
	"noah/pkg/conn"
	"noah/pkg/utils"
	"os"
	"os/exec"
	"runtime"
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
func (h *Handler) Ping() {
	sleepTime := 30 * time.Second

	for {
		time.Sleep(sleepTime)

		h.Connected = false

		err := h.ServerIsAvailable()
		if err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
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
	systemInfo, err := h.Services.GetSystemInfo()
	if err != nil {
		h.Log("[!] Error getting system info:", err.Error())
		return err
	}
	body, err := json.Marshal(systemInfo)
	if err != nil {
		return err
	}
	res, err := h.Gateway.NewRequest(http.MethodPost, fmt.Sprintf("/client/%d/health", uint64(h.ClientID)), body)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("error with status code: %d, msg: %s", res.Code, res.Msg)
	}
	return nil
}

func (h *Handler) WebsocketConnection() {
retry:
	time.Sleep(time.Second * 3)
	wsconn, err := ws.NewConnection(h.Configuration, fmt.Sprintf("/client/%d/ws", h.ClientID))
	if err != nil {
		h.Log("[!] Error connecting to server: ", err.Error())
		goto retry
	}
	listen := conn.NewMux(wsconn)
	for {
		srcConn, err := listen.Accept()
		if err != nil {
			h.Log("[!] Error accepting connection: ", err.Error())
			goto retry
		}
		if c, ok := srcConn.(*conn.Conn); ok {
			go h.handleConn(c)
		}
	}
}

func (h *Handler) handleConn(srcConn *conn.Conn) {
	defer srcConn.Close()
	lk := srcConn.GetLk()
	if lk.Network == "" {
		return
	}
	var target io.ReadWriteCloser
	switch lk.Network {
	case "pty":
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case `linux`:
			cmd = exec.Command("bash")
		case `darwin`:
			cmd = exec.Command("zsh")
		default:
			return
		}

		// 打开 pty
		ptmx, err := pty.Start(cmd)
		if err != nil {
			h.Log("Error starting PTY:", err)
			return
		}

		target = &conn.PtyReaderWriterCloser{IO: ptmx}
	case "tcp":
		t, err := net.DialTimeout(lk.Network, lk.Addr, time.Second*5)
		if err != nil {
			return
		}
		target = t
	case "cmd":
		var cmd entitie.MessageType
		// 字符串转MessageType
		cmd = entitie.MessageTypeFromString(lk.Addr)
		data := make([]byte, 1024)
		n, err := srcConn.Read(data)
		if err != nil {
			return
		}
		res, err := h.handleMessage(cmd, data[:n])
		if err != nil {
			res = err.Error()
		}
		b, err := utils.AnyToBytes(res)
		if err != nil {
			b = []byte(err.Error())
		}
		srcConn.Write(b)
	}
	if target != nil {
		defer target.Close()
		srcConn.Copy(target)
	}
}

func (h *Handler) handleMessage(messageType entitie.MessageType, data []byte) (response any, err error) {
	switch messageType {
	case entitie.MessageTypeCommand:
		var commandRequest entitie.CommandReq
		if err := json.Unmarshal(data, &commandRequest); err != nil {
			return nil, err
		}
		return h.Services.Command.Run(commandRequest.Command)
	case entitie.MessageTypeDownload:
		var downloadParams entitie.DownloadReq
		err := json.Unmarshal(data, &downloadParams)
		if err != nil {
			return nil, err
		}
		// 下载文件
		err = h.Services.Download.DownloadFile(downloadParams.Filename, downloadParams.Path)

		if err != nil {
			return nil, err
		}
	case entitie.MessageTypeUpdate:
		filename := string(data)

		// 下载文件
		filepath := "/tmp/" + filename
		err = h.Services.Download.DownloadFile(filename, filepath)
		if err != nil {
			return nil, err
		}

		//删除服务器文件
		h.Gateway.NewRequest(http.MethodDelete, "/file/"+filename, nil)

		// 设置新版本文件的执行权限
		err = os.Chmod(filepath, 0755)
		if err != nil {
			return nil, err
		}
		// 使用 nohup 命令启动新进程
		cmd := exec.Command("nohup", filepath, "&")

		err = cmd.Start()
		if err != nil {
			return nil, err
		}

		// 等待一段时间以确保新进程已经启动
		time.Sleep(1 * time.Second)

		// 确保新进程已经启动后再退出当前进程
		os.Exit(0)
	case entitie.MessageTypeExit:
		os.Exit(0)
	case entitie.MessageTypeFileExplorer:
		var fileExplorerQuery entitie.FileExplorerQuery
		err := json.Unmarshal(data, &fileExplorerQuery)
		if err != nil {
			return nil, err
		}
		op := fileExplorerQuery.Op
		path := fileExplorerQuery.Path
		switch op {
		case "list":
			res, err := h.Services.FileExplorer.GetFileExplorer(path)
			if err != nil {
				return nil, err
			}
			response = res
		case "cat":
			res, err := h.Services.FileExplorer.ReadFile(path)
			if err != nil {
				return nil, err
			}
			response = res
		case "rename":
			newFilename := fileExplorerQuery.Filename
			err := h.Services.FileExplorer.Rename(path, newFilename)
			if err != nil {
				return nil, err
			}
		case "remove":
			err := h.Services.FileExplorer.Remove(path)
			if err != nil {
				return nil, err
			}
		case "edit":
			fileContent := fileExplorerQuery.FileContent
			err := h.Services.FileExplorer.WriteFile(path, []byte(fileContent))
			if err != nil {
				return nil, err
			}
		case "mkdir":
			err := h.Services.FileExplorer.MkDir(path)
			if err != nil {
				return nil, err
			}
		}
	case entitie.MessageTypeSystemInfo:
		var systemInfoReq entitie.SystemInfoReq
		err := json.Unmarshal(data, &systemInfoReq)
		if err != nil {
			return nil, err
		}
		switch systemInfoType := systemInfoReq.SystemInfoType; systemInfoType {
		case "process":
			switch action := systemInfoReq.Action; action {
			case "list":
				process, err := h.Services.Information.GetProcessList()
				if err != nil {
					return nil, err
				}
				response = process
			case "kill":
				pid, _ := strconv.Atoi(systemInfoReq.Params)
				err := h.Services.Information.KillProcess(int32(pid))
				if err != nil {
					return nil, err
				}
			}
		case "net":
			switch action := systemInfoReq.Action; action {
			case "list":
				networkInfo, err := h.Services.Information.GetNetworkInfo()
				if err != nil {
					return nil, err
				}
				response = networkInfo
			}
		case "docker":
			switch action := systemInfoReq.Action; action {
			case "containerList":
				dockerContainerList, err := h.Services.Information.GetDockerContainerList()
				if err != nil {
					return nil, err
				}
				response = dockerContainerList
			}

		default:
		}
	}
	return response, nil
}
