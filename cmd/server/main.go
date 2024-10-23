package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/lib/liblogger"
	"main/internal/middleware/midlogger"
	"main/internal/storage/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const LocalFilePath = "config-yaml/local.yaml"

func main() {
	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger
	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start event-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := postgresql.NewPosgreSQL(cfg.GetDataSourceName())
	if err != nil {
		log.Error("faile to init storage", liblogger.Err(err))
	} else {
		log.Info("storage are enabled")
	}

	// for test run)
	_ = storage

	//init chi router
	router := chi.NewRouter()
	//init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)

	log.Error("server stopped")
}
