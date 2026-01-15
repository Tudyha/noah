package sys

import (
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type BasicInfo struct {
	Hostname   string
	Username   string
	Uid        string
	Gid        string
	Os         string
	OsArch     string
	KernelArch string
}

type HostInfo struct {
	OsName          string
	Uptime          uint64
	BootTime        uint64
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	HostId          string
}

func GetBasicInfo() (*BasicInfo, error) {
	info := &BasicInfo{}
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return info, err
	}
	info.Hostname = hostname

	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		return info, err
	}
	info.Username = currentUser.Username
	info.Uid = currentUser.Uid
	info.Gid = currentUser.Gid

	// 操作系统和架构
	info.Os = runtime.GOOS
	info.OsArch = runtime.GOARCH
	info.KernelArch = runtime.GOARCH
	return info, nil
}

func GetHostInfo() (*HostInfo, error) {
	info := &HostInfo{}

	hostInfo, err := host.Info()
	if err != nil {
		return info, err
	}
	// 填充主机信息
	info.OsName = hostInfo.OS
	info.Platform = hostInfo.Platform
	info.PlatformFamily = hostInfo.PlatformFamily
	info.PlatformVersion = hostInfo.PlatformVersion
	info.KernelVersion = hostInfo.KernelVersion
	info.HostId = hostInfo.HostID

	// 启动时间相关
	info.BootTime = uint64(hostInfo.BootTime)
	if hostInfo.Uptime > 0 {
		info.Uptime = uint64(hostInfo.Uptime)
	} else {
		info.Uptime = uint64(time.Now().Unix()) - hostInfo.BootTime
	}
	return info, nil
}
