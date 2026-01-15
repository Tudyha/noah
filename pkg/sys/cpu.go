package sys

import (
	"runtime"

	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v4/cpu"
)

type CpuInfo struct {
	CpuNum  int32
	CpuInfo []InfoStat
}

type InfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

func GetCpuInfo() (*CpuInfo, error) {
	info := &CpuInfo{}
	info.CpuNum = int32(runtime.NumCPU())

	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	copier.Copy(info.CpuInfo, cpuInfos)

	return info, nil
}
