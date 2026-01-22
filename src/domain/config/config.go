package configDomain

type IConfig interface {
	Init(envPath string) (Config, error)
	load(envPath string)
}

type AppConfig struct {
	Development bool
	Port        int
	Host        string
	LoggingIO   bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}
