package response

// Response 统一响应格式
type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type Page[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}
