package environment

type Environment struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	HttpAddr string `yaml:"httpPort"`
	TcpAddr  string `yaml:"tcpPort"`
}
