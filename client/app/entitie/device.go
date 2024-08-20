package entitie

type Device struct {
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	UserID      string `json:"user_id"`
	OSName      string `json:"os_name"`
	OSArch      string `json:"os_arch"`
	MacAddress  string `json:"mac_address"`
	IPAddress   string `json:"ip_address"`
	Port        string `json:"port"`
	FetchedUnix int64  `json:"fetched_unix"`
}
