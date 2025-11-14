package main

import (
	"fmt"
	"log/slog"
	"movielight/internal/db"
	"movielight/internal/config"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// структура приложения
type application struct {
	config config.Config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// загружаем конфиг
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Error("config load error", "err", err)
		os.Exit(1)
	}

	// подключение к БД
	database, err := db.Open(cfg.DB.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
	}

	// сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}