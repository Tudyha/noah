package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"noah/client/app/environment"
	"noah/client/app/gateway"
	"noah/client/app/service"
	"noah/pkg/enum"
	myio "noah/pkg/io"
	"noah/pkg/mux/message"
	"noah/pkg/response"
	"noah/pkg/utils"
	"noah/pkg/utils/network"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"noah/pkg/mux"

	"github.com/creack/pty"
	"github.com/samber/do/v2"
)

var (
	interval = 30 * time.Second
)

type Handler struct {
	gateway             gateway.Gateway
	informationService  service.IInformationService
	downloadService     service.IDownloadService
	fileExplorerService service.IFileExplorerService
	Env                 *environment.Environment
}

func NewHandler(i do.Injector) (Handler, error) {
	return Handler{
		gateway:             do.MustInvoke[gateway.Gateway](i),
		informationService:  do.MustInvoke[service.IInformationService](i),
		downloadService:     do.MustInvoke[service.IDownloadService](i),
		fileExplorerService: do.MustInvoke[service.IFileExplorerService](i),
		Env:                 do.MustInvoke[*environment.Environment](i),
	}, nil
}

func (h *Handler) Connect() {
retry:
	//tcp 长连接
	err := h.connect()
	if err != nil {
		log.Printf("Failed to connect to server: %v", err)
		time.Sleep(interval)
		goto retry
	}

}

func (h *Handler) connect() error {
	addr := fmt.Sprintf("%s:%d", h.Env.Server.Host, h.Env.Server.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := h.informationService.LoadClientSpecs()
	if err != nil {
		return err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/api/client/connect", addr)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))

	_, err = conn.Write([]byte(network.ConvertRequestToString(req)))
	if err != nil {
		log.Printf("Failed to write data to server: %v", err)
		return err
	}

	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to connect to server: %s", res.Status)
	}

	contentLength := res.ContentLength

	body := make([]byte, contentLength)

	_, err = io.ReadFull(res.Body, body)
	if err != nil {
		return err
	}
	var response response.Response
	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("failed to connect to server: %v", response.Msg)
	}

	m := mux.NewMux(conn, conn)

	m.SetPingHandler(func() []byte {
		systemInfo, err := h.informationService.GetSystemInfo()
		if err != nil {
			return nil
		}
		data, err := utils.AnyToJsonBytes(systemInfo)
		if err != nil {
			return nil
		}
		return data
	})

	err = m.Start()
	if err != nil {
		return err
	}
	defer m.Close()

	log.Println("Connected to server")

	for {
		conn, err := m.Accept()
		if err != nil {
			return err
		}
		go h.handlerConnection(conn.(*mux.Conn))
	}
}

func (h *Handler) handlerConnection(conn *mux.Conn) {
	defer conn.Close()

	var target io.ReadWriteCloser
	switch conn.GetNetwork() {
	case "tcp":
		tc, err := net.Dial("tcp", conn.GetAddr())
		if err != nil {
			return
		}

		target = tc
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
			return
		}

		target = &myio.PtyReaderWriterCloser{IO: ptmx}
	case "cmd":
		addr := conn.GetAddr()
		var cmdInfo message.CmdInfo
		if err := json.Unmarshal([]byte(addr), &cmdInfo); err != nil {
			return
		}
		result, err := h.handlerCommand(cmdInfo.Cmd, cmdInfo.Data)
		if err != nil {
			result = err.Error()
		}
		b, err := utils.AnyToJsonBytes(result)
		if err != nil {
			b = []byte(err.Error())
		}
		conn.Write(b)
	}
	if target != nil {
		defer target.Close()
		if err := myio.Copy(conn, target); err != nil {
			return
		}
	}
}

func (h *Handler) handlerCommand(cmd string, data []byte) (response any, err error) {
	switch cmd {
	case enum.Download:
		var downloadParams message.DownloadReq
		err := json.Unmarshal(data, &downloadParams)
		if err != nil {
			return nil, err
		}
		// 下载文件
		err = h.downloadService.DownloadFile(downloadParams.Filename, downloadParams.Path)

		if err != nil {
			return nil, err
		}
	case enum.Update:
		filename := string(data)

		// 下载文件
		filepath := "/tmp/" + filename
		err = h.downloadService.DownloadFile(filename, filepath)
		if err != nil {
			return nil, err
		}

		//删除服务器文件
		h.gateway.NewRequest(http.MethodDelete, "/file/"+filename, nil)

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
	case enum.Exit:
		os.Exit(0)
	case enum.FileExplorer:
		var fileExplorerQuery message.FileExplorerReq
		err := json.Unmarshal(data, &fileExplorerQuery)
		if err != nil {
			return nil, err
		}
		op := fileExplorerQuery.Op
		path := fileExplorerQuery.Path
		switch op {
		case "list":
			res, err := h.fileExplorerService.GetFileExplorer(path)
			if err != nil {
				return nil, err
			}
			response = res
		case "cat":
			res, err := h.fileExplorerService.ReadFile(path)
			if err != nil {
				return nil, err
			}
			response = res
		case "rename":
			newFilename := fileExplorerQuery.Filename
			err := h.fileExplorerService.Rename(path, newFilename)
			if err != nil {
				return nil, err
			}
		case "remove":
			err := h.fileExplorerService.Remove(path)
			if err != nil {
				return nil, err
			}
		case "edit":
			fileContent := fileExplorerQuery.FileContent
			err := h.fileExplorerService.WriteFile(path, []byte(fileContent))
			if err != nil {
				return nil, err
			}
		case "mkdir":
			err := h.fileExplorerService.MkDir(path)
			if err != nil {
				return nil, err
			}
		}
	case enum.SystemInfo:
		var systemInfoReq message.SystemInfoReq
		err := json.Unmarshal(data, &systemInfoReq)
		if err != nil {
			return nil, err
		}
		switch systemInfoType := systemInfoReq.SystemInfoType; systemInfoType {
		case "process":
			switch action := systemInfoReq.Action; action {
			case "list":
				process, err := h.informationService.GetProcessList()
				if err != nil {
					return nil, err
				}
				response = process
			case "kill":
				pid, _ := strconv.Atoi(systemInfoReq.Params)
				err := h.informationService.KillProcess(int32(pid))
				if err != nil {
					return nil, err
				}
			}
		case "net":
			switch action := systemInfoReq.Action; action {
			case "list":
				networkInfo, err := h.informationService.GetNetworkInfo()
				if err != nil {
					return nil, err
				}
				response = networkInfo
			}
		case "docker":
			switch action := systemInfoReq.Action; action {
			case "containerList":
				dockerContainerList, err := h.informationService.GetDockerContainerList()
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
