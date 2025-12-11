package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		return
	}
	node = n
}

func GenID() int64 {
	return node.Generate().Int64()
}
