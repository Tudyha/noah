package request

type ClientQueryRequest struct {
	PageQuery
}

type ClientGenerateV2raySubscribeRequest struct {
	Ids []uint64 `form:"ids" json:"ids" binding:"required"`
}
