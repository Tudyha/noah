package environment

type Environment struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
