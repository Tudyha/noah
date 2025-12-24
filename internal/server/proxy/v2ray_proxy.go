package proxy

import (
	"context"
	"fmt"
	"noah/pkg/packet"
	"noah/pkg/proto"

	mysession "noah/internal/session"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/common/task"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
	"github.com/v2fly/v2ray-core/v5/transport"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
)

func init() {
	common.Must(common.RegisterConfig((*proto.SimplifiedConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		h := new(handler)
		if err := core.RequireFeatures(ctx, func() error {
			return h.Init(config.(*proto.SimplifiedConfig))
		}); err != nil {
			return nil, err
		}
		return h, nil
	}))
}

type handler struct {
	config *proto.SimplifiedConfig
}

func (h *handler) Init(config *proto.SimplifiedConfig) error {
	h.config = config

	return nil
}

func (h *handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return fmt.Errorf("target not specified")
	}
	destination := outbound.Target

	input := link.Reader
	output := link.Writer

	account := session.InboundFromContext(ctx).User.Account.(*vmess.MemoryAccount)

	var tunnelType packet.OpenTunnel_TuunnelType
	switch destination.Network {
	case net.Network_TCP:
		tunnelType = packet.OpenTunnel_TCP
	case net.Network_UDP:
		tunnelType = packet.OpenTunnel_UDP
	default:
		return fmt.Errorf("unknown network: %s", destination.Network)
	}
	c1, err := mysession.GetSessionManager().OpenTunnel(account.ID.String(), tunnelType, destination.NetAddr())
	if err != nil {
		return err
	}
	defer c1.Close()

	requestDone := func() error {

		var writer buf.Writer
		if destination.Network == net.Network_TCP {
			writer = buf.NewWriter(c1)
		} else {
			writer = &buf.SequentialWriter{Writer: c1}
		}

		if err := buf.Copy(input, writer); err != nil {
			return err
		}

		return nil
	}

	responseDone := func() error {

		var reader buf.Reader
		if destination.Network == net.Network_TCP {
			reader = buf.NewReader(c1)
		} else {
			reader = buf.NewPacketReader(c1)
		}
		if err := buf.Copy(reader, output); err != nil {
			return err
		}

		return nil
	}

	if err := task.Run(ctx, requestDone, task.OnSuccess(responseDone, task.Close(output))); err != nil {
		return err
	}

	return nil
}
