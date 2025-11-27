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
		Dsn         string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	CORS struct {
		TrustedOrigins []string
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

	// SMTP Configuration
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	if cfg.SMTP.Host == "" {
		cfg.SMTP.Host = "localhost" // значение по умолчанию
	}

	// SMTP Port
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		smtpPortStr = "587" // стандартный порт для SMTP
	}
	cfg.SMTP.Port, err = strconv.Atoi(smtpPortStr)
	if err != nil {
		return cfg, fmt.Errorf("invalid SMTP port: %w", err)
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.SMTP.Sender = os.Getenv("SMTP_SENDER")

	return cfg, nil
}