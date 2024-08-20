package network

import (
	"io"
	"log"
	"net"
	"net/http"
)

func GetMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var address []string
	for _, i := range interfaces {
		a := i.HardwareAddr.String()
		if a != "" {
			address = append(address, a)
		}
	}
	if len(address) == 0 {
		return "", nil
	}
	return address[0], nil
}

func GetLocalIP() string {
	// conn, err := net.Dial(`udp`, `8.8.8.8:80`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	// localAddr := conn.LocalAddr().(*net.UDPAddr)
	// return localAddr.IP
	// 使用http.Get函数请求外部服务
	resp, err := http.Get("http://ip.plus/myip")
	if err != nil {
		log.Fatal("请求失败: ", err)
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("读取响应体失败: ", err)
	}

	return string(body)
}
