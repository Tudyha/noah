package enum

import "runtime"

type OSType int

const (
	Unknown OSType = iota
	Windows
	Linux
	Darwin
)

var OSTargetMap = map[OSType]string{
	Windows: "Windows",
	Linux:   "Linux",
	Darwin:  "Mac OS",
}

var OSTargetIntMap = map[int]OSType{
	1: Windows,
	2: Linux,
	3: Darwin,
}

// DetectOS return an int which represent an OS type
func DetectOS(osName string) OSType {
	if osName == "" {
		osName = runtime.GOOS
	}
	switch osName {
	case `windows`:
		return Windows
	case `linux`:
		return Linux
	case `darwin`:
		return Darwin
	default:
		return Unknown
	}
}
