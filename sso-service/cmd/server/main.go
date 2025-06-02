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

const configPath = "config-yaml/dev.yaml"

func main() {

	cfg := config.MustConfig(configPath)

	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start sso-service", slog.String("env", cfg.Env))
	log.Debug("debug message are enable")

	// init app
	app := app.New(log, cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// start app
	go app.MustRun()
	log.Info("server started")
	<-done
	// stop app
	app.Stop()
	log.Info("server stopped")
}
