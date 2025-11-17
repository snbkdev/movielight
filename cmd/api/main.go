package main

import (
	"fmt"
	"log/slog"
	"movielight/internal/config"
	"movielight/internal/db"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	database, err := db.Open(cfg.DB.Dsn)
	if err != nil {
    	logger.Error(err.Error())
    	os.Exit(1)
	}
	defer database.Close()

	logger.Info("database connection pool established")

	migrationDriver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file:///Users/asanbeksamudin/Documents/projects/mine/golang/movielight/migrations", "postgres", migrationDriver)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database migrations applied")

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