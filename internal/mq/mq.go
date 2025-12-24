package mq

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var pubSub *gochannel.GoChannel

func Init() error {
	pubSub = gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
	return nil
}

func GetPubSub() *gochannel.GoChannel {
	return pubSub
}
