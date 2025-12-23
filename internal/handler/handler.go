package handler

import (
	"noah/pkg/conn"
	"noah/pkg/constant"
)

func getSessionID(ctx conn.Context) uint64 {
	return ctx.Value(constant.SESSION_ID_KEY).(uint64)
}
