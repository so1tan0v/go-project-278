package config

import (
	configDomain "link-service/src/domain/config"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init(envPath string) (configDomain.Config, error) {
	if envPath != "" {
		if err := load(envPath); err != nil {
			return configDomain.Config{}, err
		}
	}

	var development bool
	var host string
	var port int
	var loggingIO bool

	development = os.Getenv(configDomain.Development) == "true"
	host = os.Getenv(configDomain.Host)
	if host == "" {
		host = "localhost"
	}

	port, _ = strconv.Atoi(os.Getenv(configDomain.Port))
	loggingIO = os.Getenv(configDomain.LoggingIo) == "true"

	return configDomain.Config{
		Development: development,
		Host:        host,
		Port:        port,
		LoggingIO:   loggingIO,
	}, nil
}

func load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
