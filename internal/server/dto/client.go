package dto

type Command struct {
	Command   string `json:"command,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Response  []byte `json:"response,omitempty"`
	HasError  bool   `json:"has_error,omitempty"`
}

type RespondCommandRequestBody struct {
	ClientID uint   `json:"client_id,omitempty"`
	Response []byte `json:"response"`
	HasError bool   `json:"has_error,omitempty"`
}
