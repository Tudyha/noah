package environment

type Environment struct {
	Server ServerConfig `json:"server"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
