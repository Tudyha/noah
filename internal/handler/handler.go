package handler

import (
	"noah/pkg/conn"
	"noah/pkg/constant"
)

func getSessionID(ctx conn.Context) string {
	return ctx.Value(constant.SESSION_ID_KEY).(string)
}
