package environment

type Environment struct {
	Server   ServerConfig `json:"server"`
	Compress uint8        `json:"compress"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
