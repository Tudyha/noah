package entitie

type Device struct {
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	UserID      string `json:"userId"`
	OSName      string `json:"osName"`
	OSArch      string `json:"osArch"`
	MacAddress  string `json:"macAddress"`
	IPAddress   string `json:"ipAddress"`
	Port        string `json:"port"`
	FetchedUnix int64  `json:"fetchedUnix"`
}
