package config

import (
	"fmt"
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

	appConfig := initAppConfig()
	dbConfig, err := initDatabaseConfig()
	if err != nil {
		return configDomain.Config{}, err
	}

	return configDomain.Config{
		App:      appConfig,
		Database: dbConfig,
	}, nil
}

func initAppConfig() configDomain.AppConfig {
	var development bool
	var host string
	var port int
	var loggingIO bool

	development = os.Getenv("DEVELOPMENT") == "true"
	host = os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port, _ = strconv.Atoi(os.Getenv("PORT"))
	loggingIO = os.Getenv("LOGGING_IO") == "true"

	return configDomain.AppConfig{
		Development: development,
		Port:        port,
		Host:        host,
		LoggingIO:   loggingIO,
	}
}

func initDatabaseConfig() (configDomain.DatabaseConfig, error) {
	var host string
	var port int
	var user string
	var password string
	var database string

	host = os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	if port == 0 {
		port = 5432
	}
	user = os.Getenv("DB_USER")
	if user == "" {
		return configDomain.DatabaseConfig{}, fmt.Errorf("DB_USER environment variable not set")
	}
	password = os.Getenv("DB_PASSWORD")
	if password == "" {
		return configDomain.DatabaseConfig{}, fmt.Errorf("DB_PASSWORD environment variable not set")
	}

	return configDomain.DatabaseConfig{
		Database: database,
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
	}, nil
}

func load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
