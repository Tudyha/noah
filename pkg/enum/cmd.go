package enum

type Cmd string

var (
	Command      Cmd = "command"
	Exit         Cmd = "exit"
	Update       Cmd = "update"
	Download     Cmd = "download"
	FileExplorer Cmd = "fileExplorer"
	SystemInfo   Cmd = "systemInfo"
	Unknown      Cmd = "unknown"
	Process      Cmd = "process"
)
