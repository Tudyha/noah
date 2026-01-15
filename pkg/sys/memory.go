package sys

import "github.com/shirou/gopsutil/v4/mem"

func GetTotalMemory() uint64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return memInfo.Total
}
