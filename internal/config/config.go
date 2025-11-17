package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DSN string
}

type Config struct {
	Port int
	Env  string
	DB   struct {
		Dsn string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime time.Duration
	}
}

// Load загружает конфигурацию из .env
func Load(logger *slog.Logger) (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using defaults")
	}

	// PORT
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "4550"
	}
	cfg.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	// ENV
	cfg.Env = os.Getenv("ENV")
	if cfg.Env == "" {
		cfg.Env = "development"
	}

	// DB_DSN
	cfg.DB.Dsn = os.Getenv("DB_DSN")
	if cfg.DB.Dsn == "" {
		fmt.Println("No database connection")
	}

	return cfg, nil
}