package utils

import (
	"encoding/json"
	"math"
)

func AnyToJsonBytes(data any) (b []byte, err error) {
	switch d := data.(type) {
	case []byte:
		b = d
	case string:
		b = []byte(d)
	default:
		b, err = json.Marshal(data)
	}
	return b, err
}

func CoverToGb(b uint64) float64 {
	f := float64(b) / (1024 * 1024 * 1024)
	return RoundToTwoDecimals(f)
}

// CoverToKb 将字节数转换为GB, 保留两位小数
func CoverToKb(b float64) float64 {
	f := float64(b) / (1024)
	return RoundToTwoDecimals(f)
}

// RoundToTwoDecimals 将浮点数保留两位小数
func RoundToTwoDecimals(f float64) float64 {
	return math.Round(f*100) / 100
}
