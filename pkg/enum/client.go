package enum

type ClientStatus int

const (
	ClientStatusOnline  = iota + 1 // 在线
	ClientStatusOffline            // 离线
)
