package configDomain

const (
	LoggingIo   = "LOGGING_IO"
	Development = "DEVELOPMENT"
	Port        = "PORT"
	Host        = "HOST"
)

type IConfig interface {
	Init(envPath string) (Config, error)
	load(envPath string)
}

type Config struct {
	Development bool
	Port        int
	Host        string
	LoggingIO   bool
}
