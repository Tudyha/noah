package request

type PageQuery struct {
	Page  int `json:"page" form:"page" binding:"required,min=1"`
	Limit int `json:"limit" form:"limit" binding:"required,min=1,max=100"`
}
