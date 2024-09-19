package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"noah/client/app/environment"
	"noah/client/app/service"
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

		wsMessageType, wsMessage, err := h.Connection.ReadMessage()
		if err != nil {
			h.Log("[!] Error reading from connection:", err)
			h.Reconnect()
			continue
		}

		var message entitie.Message
		if err := json.Unmarshal(wsMessage, &message); err != nil {
			continue
		}

		response, err := h.handleMessage(wsMessageType, message)
		errMsg := ""
		if err != nil {
			h.Log("[!] Error handling message:", err)
			errMsg = err.Error()
		}

		if message.MessageType == entitie.MessageTypeChannel {
			continue
		}
		ws.WriteMessage(h.Connection, message.MessageId, message.MessageType, response, errMsg)
	}
}

func (h *Handler) handleMessage(wsMessageType int, message entitie.Message) (response any, err error) {
	switch message.MessageType {
	case entitie.MessageTypeCommand:
		var commandRequest entitie.CommandReq
		if err := json.Unmarshal(message.Data, &commandRequest); err != nil {
			return nil, err
		}
		return h.Services.Command.Run(commandRequest.Command)
	case entitie.MessageTypeChannel:
		var channelRequest entitie.ChannelReq
		if err := json.Unmarshal(message.Data, &channelRequest); err != nil {
			fmt.Println("Error unmarshalling channel request:", err.Error())
			return nil, err
		}
		switch channelRequest.Action {
		case "open":
			addr := fmt.Sprintf("%s:%d", channelRequest.LocalIp, channelRequest.LocalPort)
			err := h.Services.Channel.NewChannel(channelRequest.ChannelId, channelRequest.ChannelType, h.Connection, addr)
			if err != nil {
				fmt.Println("Error opening channel:", err.Error())
				return nil, err
			}
		case "write":
			err := h.Services.Channel.Write(wsMessageType, channelRequest.ChannelId, channelRequest.ChannelData)
			if err != nil {
				fmt.Println("Error writing to channel:", err)
				return nil, err
			}
		}
	case entitie.MessageTypeDownload:
		var downloadParams entitie.DownloadReq
		err := json.Unmarshal(message.Data, &downloadParams)
		if err != nil {
			return nil, err
		}
		// 下载文件
		err = h.Services.Download.DownloadFile(downloadParams.Filename, downloadParams.Path)

		if err != nil {
			return nil, err
		}
	case entitie.MessageTypeUpdate:
		filename := string(message.Data)

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
		err := json.Unmarshal(message.Data, &fileExplorerQuery)
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
		err := json.Unmarshal(message.Data, &systemInfoReq)
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
