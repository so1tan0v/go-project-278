package config

import (
	"fmt"
	configDomain "link-service/src/domain/config"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init(envPath string) (*configDomain.Config, error) {
	// Для локальной разработки полезно подгружать переменные из .env.
	// В продакшене переменные окружения должны быть заданы извне.
	if envPath != "" {
		if err := load(envPath); err != nil {
			return nil, err
		}
	} else {
		// Игнорируем ошибку, если файла нет.
		_ = godotenv.Load()
	}

	appConfig := initAppConfig()
	dbConfig, err := initDatabaseConfig()
	if err != nil {
		return nil, err
	}

	return &configDomain.Config{
		App:      *appConfig,
		Database: *dbConfig,
	}, nil
}

func initAppConfig() *configDomain.AppConfig {
	var development bool
	var host string
	var port int
	var loggingIO bool
	var baseURL string

	development = os.Getenv("DEVELOPMENT") == "true"
	host = os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port, _ = strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}
	loggingIO = os.Getenv("LOGGING_IO") == "true"

	return &configDomain.AppConfig{
		Development: development,
		Port:        port,
		Host:        host,
		LoggingIO:   loggingIO,
		BaseURL: func() string {
			baseURL = os.Getenv("BASE_URL")
			if baseURL != "" {
				return baseURL
			}
			// Дефолт удобен для локальной разработки.
			return fmt.Sprintf("http://%s:%d", host, port)
		}(),
	}
}

func initDatabaseConfig() (*configDomain.DatabaseConfig, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	return &configDomain.DatabaseConfig{
		URL: databaseURL,
	}, nil
}

func load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
