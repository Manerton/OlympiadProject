package main

import (
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"main/internal/lib/liblogger"
	"os"
	"os/signal"
	"syscall"
)

const LocalFilePath = "config-yaml/dev.yaml"

// @title Event Service API
// @version 1.0
// @description Документация к микросервису событий
// @host 172.16.1.39:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Events
var DockerFilePath string = os.Getenv("CONFIG_PATH")

func main() {
	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger
	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start event-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init app
	app := app.New(log, cfg)
	// init gorutine for start server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// start app
	go app.MustRun()
	log.Info("server started")
	<-done
	// stop app
	app.Stop()
	log.Info("server stopped")
}
