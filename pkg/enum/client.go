package enum

type ClientStatus int

const (
	ClientStatusOnline  = iota + 1 // 在线
	ClientStatusOffline            // 离线
)

type ClientOsType int

const (
	ClientOsTypeWindows = iota + 1 // Windows
	ClientOsTypeMac                // Mac
	ClientOsTypeLinux              // Linux
)

var ClientOsNameToOsTypeMap = map[string]ClientOsType{
	"Windows": ClientOsTypeWindows,
	"darwin":  ClientOsTypeMac,
	"Linux":   ClientOsTypeLinux,
}
