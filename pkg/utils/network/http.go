package network

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ConvertRequestToString(req *http.Request) string {
	var buf bytes.Buffer

	// 写入请求行
	buf.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL.RequestURI()))
	// 写入请求host
	buf.WriteString(fmt.Sprintf("Host: %s\r\n", req.Host))

	// 写入请求头
	for key, values := range req.Header {
		for _, value := range values {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}

	// 写入空行
	buf.WriteString("\r\n")

	// 写入请求体
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		buf.Write(bodyBytes)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置请求体
	}

	return buf.String()
}
