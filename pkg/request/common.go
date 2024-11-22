package request

type PageQuery struct {
	Page int `form:"page" default:"1"`
	Size int `form:"size" default:"10"`
}

type PageInfo struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
