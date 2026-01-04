package request

import "noah/pkg/enum"

type AgentQueryRequest struct {
	PageQuery

	Status   enum.AgentStatus `json:"status" form:"status"`
	Hostname string           `json:"hostname" form:"hostname"`
}

type AgentGenerateV2raySubscribeRequest struct {
	Ids []uint64 `form:"ids" json:"ids" binding:"required"`
}
