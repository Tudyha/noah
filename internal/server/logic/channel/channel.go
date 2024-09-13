package channel

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"noah/internal/server/config"
	"noah/internal/server/enum"
	"noah/internal/server/middleware/log"
	"noah/internal/server/utils"
	"sync"
	"time"
)

type channelService struct {
	channels map[string]*Channel
}

func NewChannelService() *channelService {
	return &channelService{
		channels: make(map[string]*Channel),
	}
}

type Channel struct {
	ChannelId   string           // 通道id
	ChannelType enum.ChannelType // 通道类型
	conn        *websocket.Conn  // 前端websocket连接
	closeSignal chan struct{}    // 用于关闭连接的信号
	isClosed    bool             // 标记是否已经关闭
	mu          sync.Mutex       // 用于保护对channelId的访问
	clientConn  *websocket.Conn  // 客户端连接
}

func (c channelService) NewChannel(channelType enum.ChannelType) (channel *Channel) {
	channelId := utils.RandString(16)
	channel = &Channel{
		ChannelId:   channelId,
		ChannelType: channelType,
		isClosed:    false,
		closeSignal: make(chan struct{}),
	}

	c.channels[channelId] = channel
	return channel
}

func (c channelService) ClientConnect(channelId string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("connection is nil")
	}

	channel, ok := c.channels[channelId]
	if !ok {
		return errors.New("channel is not found")
	}

	channel.clientConn = conn

	go channel.Listen()

	return nil
}

func (channel *Channel) Listen() {
	log.Info("starting server on port 3322", nil)
	//监听指定端口，同时把请求转发到客户端
	listener, err := net.Listen("tcp", ":3322")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// 监听连接
	for {
		conn, err := listener.Accept()
		fmt.Println("accepted connection from:")
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 线程中处理连接
		go func() {
			for {
				b := make([]byte, 1024)
				n, err := conn.Read(b)
				if err != nil {
					fmt.Println("Error reading from connection:", err)
					return
				}
				if n > 0 {
					channel.clientConn.WriteMessage(websocket.BinaryMessage, b)
				}
			}
		}()

		go func() {
			for {
				_, b, err := channel.clientConn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading from client:", err)
					return
				}
				conn.Write(b)
			}
		}()
	}
}

func (channel *Channel) Start(conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("connection is nil")
	}

	// 是否需要判断客户端是否已连接？

	channel.conn = conn

	go channel.read()
	go channel.write()

	return nil
}

func (channel *Channel) write() {
	defer func() {
		if err := channel.close(); err != nil {
			log.Error("channel write goroutine closed error: "+err.Error(), nil)
		}
	}()

	for {
		select {
		case <-channel.closeSignal:
			return
		default:
			time.Sleep(10 * time.Millisecond)

			if channel.clientConn == nil {
				continue
			}

			msgType, data, err := channel.clientConn.ReadMessage()
			if err != nil {
				log.Error("channel write goroutine ReadMessage error: "+err.Error(), nil)
				return
			}

			channel.conn.SetWriteDeadline(time.Now().Add(config.MessageWait))
			if err := channel.conn.WriteMessage(msgType, data); err != nil {
				log.Error("channel write goroutine WriteMessage error: "+err.Error(), nil)
				return
			}
		}
	}
}

func (channel *Channel) read() {
	defer func() {
		if err := channel.close(); err != nil {
			log.Error("channel read goroutine closed error: "+err.Error(), nil)
		}
	}()

	for {
		select {
		case <-channel.closeSignal:
			return
		default:
			if channel.clientConn == nil {
				continue
			}
			msgType, data, err := channel.conn.ReadMessage()
			if err != nil {
				log.Error("channel read goroutine ReadMessage error: "+err.Error(), nil)
				return
			}

			err = channel.clientConn.WriteMessage(msgType, data)
			if err != nil {
				log.Error("channel read goroutine WriteMessage error: "+err.Error(), nil)
				return
			}
		}
	}
}

func (channel *Channel) close() error {
	channel.mu.Lock()
	defer channel.mu.Unlock()

	if channel.isClosed {
		return nil
	}

	if channel.conn != nil {
		channel.conn.Close()
	}

	if channel.clientConn != nil {
		channel.clientConn.Close()
	}

	close(channel.closeSignal)
	channel.isClosed = true

	return nil
}
