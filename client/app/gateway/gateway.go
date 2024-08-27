package gateway

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Gateway interface {
	NewRequest(method string, url string, body []byte) (*HttpResponse, error)
	NewFileDownloadRequest(method string, url string, body []byte) ([]byte, error)
}
