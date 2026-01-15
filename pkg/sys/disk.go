package sys

import (
	"fmt"
	"runtime"

	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiskTotal() uint64 {
	partitions, err := getDiskUsage()
	if err != nil {
		return 0
	}

	total := lo.Reduce(partitions, func(v uint64, item *disk.UsageStat, _ int) uint64 {
		return v + item.Total
	}, 0)
	return total
}

func getDiskUsage() ([]*disk.UsageStat, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		u, err := disk.Usage("/")
		return []*disk.UsageStat{u}, err
	case "windows":
		parts, err := disk.Partitions(false)
		if err != nil {
			return nil, err
		}
		var us []*disk.UsageStat
		for _, p := range parts {
			u, err := disk.Usage(p.Mountpoint)
			if err != nil {
				return nil, err
			}
			us = append(us, u)
		}
		return us, nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
