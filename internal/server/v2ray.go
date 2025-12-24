package server

import (
	"context"
	"noah/pkg/app"
	"noah/pkg/constant"
	"noah/pkg/proto"
	"noah/pkg/utils"

	"noah/internal/mq"
	_ "noah/internal/server/proxy"
	"noah/pkg/config"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/outbound"
	"github.com/v2fly/v2ray-core/v5/app/router"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	_ "github.com/v2fly/v2ray-core/v5/proxy/vmess"
	vmessIn "github.com/v2fly/v2ray-core/v5/proxy/vmess/inbound"
	"google.golang.org/protobuf/types/known/anypb"

	gonet "net"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	mapset "github.com/deckarep/golang-set/v2"
)

type V2rayServer struct {
	addr    string
	logPath string

	config   *core.Config
	instance *core.Instance

	sub     *gochannel.GoChannel
	clients mapset.Set[string]
}

func NewV2rayServer() app.Server {
	cfg := config.Get()

	s := new(V2rayServer)
	s.addr = cfg.Server.V2ray.Addr
	s.logPath = cfg.Server.V2ray.LogPath
	s.clients = mapset.NewSet[string]()
	s.initV2rayInstanceConfig()
	s.sub = mq.GetPubSub()
	return s
}

func (s *V2rayServer) subscribe() {
	go func() {
		messages, err := s.sub.Subscribe(context.Background(), constant.MQ_TOPIC_CLIENT_OFFLINE)
		if err != nil {
			return
		}
		for msg := range messages {
			s.clients.Remove(string(msg.Payload))
			s.restart()

			// we need to Acknowledge that we received and processed the message,
			// otherwise, it will be resent over and over again.
			msg.Ack()
		}
	}()
	messages, err := s.sub.Subscribe(context.Background(), constant.MQ_TOPIC_CLIENT_ONLINE)
	if err != nil {
		return
	}
	for msg := range messages {
		s.clients.Add(string(msg.Payload))
		s.restart()

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

func (s *V2rayServer) Start(ctx context.Context) (err error) {
	go s.subscribe()

	return s.start()
}

func (s *V2rayServer) start() (err error) {
	s.instance, err = core.New(s.config)
	if err != nil {
		return err
	}

	if err := s.instance.Start(); err != nil {
		return err
	}

	return nil
}

func (s *V2rayServer) restart() {
	if s.instance != nil {
		s.instance.Close()
	}
	s.initV2rayInstanceConfig()
	s.start()
}

func (s *V2rayServer) Stop(ctx context.Context) error {
	if s.instance != nil {
		s.instance.Close()
	}
	if s.sub != nil {
		s.sub.Close()
	}
	return nil
}

func (s *V2rayServer) String() string {
	return "v2ray service address: " + s.addr
}

func (s *V2rayServer) initV2rayInstanceConfig() {
	host, port, _ := gonet.SplitHostPort(s.addr)
	p, _ := utils.StringToUint64(port)
	config := &core.Config{
		App: []*anypb.Any{
			// serial.ToTypedMessage(&log.Config{
			// 	Access: &log.LogSpecification{
			// 		Type:  log.LogType_File,
			// 		Level: 1,
			// 		Path:  s.logPath,
			// 	},
			// 	Error: &log.LogSpecification{
			// 		Type:  log.LogType_File,
			// 		Level: 1,
			// 		Path:  s.logPath,
			// 	},
			// }),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(&router.Config{}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(net.Port(p)),
					Listen:    net.NewIPOrDomain(net.ParseAddress(host)),
				}),
				ProxySettings: serial.ToTypedMessage(&vmessIn.SimplifiedConfig{
					Users: s.clients.ToSlice(),
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&proto.SimplifiedConfig{}),
			},
		},
	}
	s.config = config

}
