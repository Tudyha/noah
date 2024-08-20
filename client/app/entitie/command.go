package entitie

type Command struct {
	ClientID  uint   `json:"client_id,omitempty"`
	Command   string `json:"command,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Response  []byte `json:"response,omitempty"`
	HasError  bool   `json:"has_error,omitempty"`
}
