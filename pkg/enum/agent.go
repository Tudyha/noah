package enum

type AgentStatus int

const (
	AgentStatusOnline  = iota + 1 // 在线
	AgentStatusOffline            // 离线
)

type AgentOsType int

const (
	AgentOsTypeWindows = iota + 1 // Windows
	AgentOsTypeMac                // Mac
	AgentOsTypeLinux              // Linux
)

var AgentOsNameToOsTypeMap = map[string]AgentOsType{
	"Windows": AgentOsTypeWindows,
	"darwin":  AgentOsTypeMac,
	"Linux":   AgentOsTypeLinux,
}
