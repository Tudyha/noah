package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"noah/internal/server/enum"
	"noah/internal/server/service"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type ShellController struct{}

func NewShellController() *ShellController {
	return &ShellController{}
}

const (
	// Time allowed to write or read a message.
	messageWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var terminalModes = ssh.TerminalModes{
	ssh.ECHO:          1,     // enable echoing (different from the example in docs)
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

type windowSize struct {
	High  int `json:"high"`
	Width int `json:"width"`
}

type sshClient struct {
	conn *websocket.Conn
	addr string
	user string
	// secret   string
	// keyfile  string
	client   *ssh.Client
	sess     *ssh.Session
	sessIn   io.WriteCloser
	sessOut  io.Reader
	closeSig chan struct{}
	clientId uint
}

func (c *sshClient) getWindowSize() (wdSize *windowSize, err error) {
	c.conn.SetReadDeadline(time.Now().Add(messageWait))
	msgType, msg, err := c.conn.ReadMessage()
	if msgType != websocket.BinaryMessage {
		err = fmt.Errorf("conn.ReadMessage: message type is not binary")
		return
	}
	if err != nil {
		err = fmt.Errorf("conn.ReadMessage: %w", err)
		return
	}

	wdSize = new(windowSize)
	if err = json.Unmarshal(msg, wdSize); err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return
	}
	return
}

func (c *sshClient) wsWrite() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	data := make([]byte, maxMessageSize)

	for {
		time.Sleep(10 * time.Millisecond)
		n, readErr := c.sessOut.Read(data)
		if n > 0 {
			c.conn.SetWriteDeadline(time.Now().Add(messageWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, data[:n]); err != nil {
				return fmt.Errorf("conn.WriteMessage: %w", err)
			}
		}
		if readErr != nil {
			return fmt.Errorf("sessOut.Read: %w", readErr)
		}
	}
}

func (c *sshClient) wsRead() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	var zeroTime time.Time
	c.conn.SetReadDeadline(zeroTime)

	for {
		msgType, connReader, err := c.conn.NextReader()
		if err != nil {
			return fmt.Errorf("conn.NextReader: %w", err)
		}
		if msgType != websocket.BinaryMessage {
			if _, err := io.Copy(c.sessIn, connReader); err != nil {
				return fmt.Errorf("io.Copy: %w", err)
			}
			continue
		}

		data := make([]byte, maxMessageSize)
		n, err := connReader.Read(data)
		if err != nil {
			return fmt.Errorf("connReader.Read: %w", err)
		}

		var wdSize windowSize
		if err := json.Unmarshal(data[:n], &wdSize); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}

		if err := c.sess.WindowChange(wdSize.High, wdSize.Width); err != nil {
			return fmt.Errorf("sess.WindowChange: %w", err)
		}
	}
}

func (c *sshClient) bridgeWSAndSSH() {
	defer c.conn.Close()

	wdSize, err := c.getWindowSize()
	if err != nil {
		log.Println("bridgeWSAndSSH: getWindowSize:", err)
		return
	}

	var auth ssh.AuthMethod
	// if c.secret != "" {
	// auth = ssh.Password(c.secret)
	// } else {
	key, err := os.ReadFile("server/config/private.key")
	if err != nil {
		log.Println("bridgeWSAndSSH: os.ReadFile:", err)
		return
	}
	privateKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("bridgeWSAndSSH: ssh.ParsePrivateKey:", err)
		return
	}
	auth = ssh.PublicKeys(privateKey)
	// }

	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{auth},
		// InsecureIgnoreHostKey returns a function
		// that can be used for ClientConfig.HostKeyCallback
		// to accept any host key.
		// It should not be used for production code.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	c.client, err = ssh.Dial("tcp", c.addr, config)
	if err != nil {
		//可能没上传公钥，执行命令上传以后再重新连接试试
		pubKey, err := os.ReadFile("server/config/public.key")
		if err != nil {
			log.Println("bridgeWSAndSSH: os.ReadFile:", err)
			return
		}
		cmd := "echo %s >> ~/.ssh/authorized_keys"
		fmt.Printf(cmd, strings.TrimSpace(string(pubKey)))
		service.GetChannelService().SendCommand(c.clientId, enum.MessageTypeCommand, fmt.Sprintf(cmd, strings.TrimSpace(string(pubKey))))

		//重新连接
		c.client, err = ssh.Dial("tcp", c.addr, config)

		if err != nil {
			log.Println("bridgeWSAndSSH:", err)
			return
		}
		// return
	}
	defer c.client.Close()

	c.sess, err = c.client.NewSession()
	if err != nil {
		log.Println("bridgeWSAndSSH: client.NewSession:", err)
		return
	}
	defer c.sess.Close()

	c.sess.Stderr = os.Stderr // TODO: check proper Stderr output
	c.sessOut, err = c.sess.StdoutPipe()
	if err != nil {
		log.Println("bridgeWSAndSSH: session.StdoutPipe:", err)
		return
	}

	c.sessIn, err = c.sess.StdinPipe()
	if err != nil {
		log.Println("bridgeWSAndSSH: session.StdinPipe:", err)
		return
	}
	defer c.sessIn.Close()

	if err := c.sess.RequestPty("xterm", wdSize.High, wdSize.Width, terminalModes); err != nil {
		log.Println("bridgeWSAndSSH: session.RequestPty:", err)
		return
	}
	if err := c.sess.Shell(); err != nil {
		log.Println("bridgeWSAndSSH: session.Shell:", err)
		return
	}

	log.Println("started a login shell on the remote host")
	defer log.Println("closed a login shell on the remote host")

	go func() {
		if err := c.wsRead(); err != nil {
			log.Println("bridgeWSAndSSH: wsRead:", err)
		}
	}()

	go func() {
		if err := c.wsWrite(); err != nil {
			log.Println("bridgeWSAndSSH: wsWrite:", err)
		}
	}()

	<-c.closeSig
}

// webSocket handles WebSocket requests for SSH from the clients.
func (h *ShellController) WebSocket(c *gin.Context) {
	// 获取路由参数clientId
	//id, _ := strconv.Atoi(c.Param("id"))
	//uintId := uint(id)

	// 获取WebSocket连接
	//conn, err := config.Upgrader.Upgrade(c.Writer, c.Request, nil)
	//if err != nil {
	//	log.Println("upgrader.Upgrade:", err)
	//	return
	//}

	//Client := dao.ClientDa.GetById(uintId)
	//
	//sshCli := &sshClient{
	//	conn: conn,
	//	// todo get port from sshd proses
	//	addr:     Client.IPAddress + ":22",
	//	user:     Client.Username,
	//	closeSig: make(chan struct{}, 1),
	//	clientId: uintId,
	//}
	//go sshCli.bridgeWSAndSSH()
}
