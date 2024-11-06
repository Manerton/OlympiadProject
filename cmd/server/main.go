package main

import (
	"main/internal/config"
	"main/internal/lib/liblogger"
	"main/internal/lib/midlogger"
	"main/internal/storage/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const localConfigPath = "config-yaml/local.yaml"

func main() {
	cfg := config.GetConfig(localConfigPath)

	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start jure-assignments-service")
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := postgresql.NewPosgreSQL(cfg.GetDataSourceName())
	if err != nil {
		log.Error("faile to init storage", liblogger.Err(err))
	} else {
		log.Info("storage are enabled")
	}

	// init chi router
	router := chi.NewRouter()
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)

	_ = storage
}
