package ip

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"

	"github.com/oschwald/maxminddb-golang"
)

var db *maxminddb.Reader

func init() {
	// 获取当前文件的绝对路径
	_, filename, _, _ := runtime.Caller(0)
	// 获取当前文件所在的目录
	dir := filepath.Dir(filename)
	// 构建数据库文件的绝对路径
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")

	if d, err := maxminddb.Open(dbPath); err != nil {
		fmt.Println(err)
	} else {
		db = d
	}
}

func GetIPCountry(ip string) string {
	if db == nil {
		fmt.Println("db is nil")
		return ""
	}

	addr := net.ParseIP(ip)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err := db.Lookup(addr, &record)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return record.Country.ISOCode
}
