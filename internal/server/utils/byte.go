package utils

import "math"

func ByteToString(value []byte) string {
	return string(value)
}

func StringToByte(value string) []byte {
	return []byte(value)
}

// CoverToGb 将字节数转换为GB, 保留两位小数
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
