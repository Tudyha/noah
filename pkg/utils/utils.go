package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MD5 计算MD5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// GenerateRandomNumber 生成随机数字
func GenerateRandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsValidUsername 验证用户名格式
func IsValidUsername(username string) bool {
	pattern := `^[a-zA-Z][a-zA-Z0-9_]{3,15}$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// IsValidPassword 验证密码强度
func IsValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	// 至少包含字母和数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

// GetLocalIP 获取本地IP地址
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}

// StringToInt 字符串转整数
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringToInt64 字符串转int64
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// StringToUint64 字符串转uint64
func StringToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// IntToString 整数转字符串
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// Int64ToString int64转字符串
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// Float64ToString float64转字符串
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// StringToFloat64 字符串转float64
func StringToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// GetFileExt 获取文件扩展名
func GetFileExt(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetFileNameWithoutExt 获取不带扩展名的文件名
func GetFileNameWithoutExt(filename string) string {
	name := filepath.Base(filename)
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}

// TruncateString 截断字符串
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// Contains 检查字符串是否包含子串
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// StartsWith 检查字符串是否以指定前缀开始
func StartsWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// EndsWith 检查字符串是否以指定后缀结束
func EndsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// JoinStrings 连接字符串数组
func JoinStrings(strs []string, separator string) string {
	return strings.Join(strs, separator)
}

// SplitString 分割字符串
func SplitString(str, separator string) []string {
	return strings.Split(str, separator)
}

// ReplaceString 替换字符串
func ReplaceString(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}

// ToLower 转换为小写
func ToLower(str string) string {
	return strings.ToLower(str)
}

// ToUpper 转换为大写
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// TrimSpace 去除首尾空格
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// FormatBytes 格式化字节数
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var address []string
	for _, i := range interfaces {
		a := i.HardwareAddr.String()
		if a != "" {
			address = append(address, a)
		}
	}
	if len(address) == 0 {
		return ""
	}
	return address[0]
}
