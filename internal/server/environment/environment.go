package environment

type Environment struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	Host     string `yaml:"host"`
	HttpPort int    `yaml:"http-port"`
	TcpPort  int    `yaml:"tcp-port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AdminConfig struct {
	Password string `yaml:"password"`
}
