package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/lib/liblogger"
	"os"
	"os/signal"
	"syscall"
)

const localConfigPath = "config-yaml/local.yaml"
const DebugLocalFilePath = "C:/code_folder/OlympiadProject/Jure-assignments-service/config-yaml/local.yaml"

func main() {
	cfg := config.GetConfig(DebugLocalFilePath)

	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start jure-assignments-service")
	log.Debug("debug messages are enabled")

	// init app
	app := app.New(log, cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// start app
	go app.MustRun()
	log.Info("server started")
	// wait
	<-done
	// stop app
	app.Stop()
	log.Info("server stopped")
}
