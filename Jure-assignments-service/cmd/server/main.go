package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/lib/liblogger"
	"os"
	"os/signal"
	"syscall"

	_ "main/docs"
)

const localConfigPath = "config-yaml/local.yaml"
const DebugLocalFilePath = "C:/code_folder/OlympiadProject/Jure-assignments-service/config-yaml/local.yaml"
const DebugDevlFilePath = "C:/Users/Admin/goProject/OlympiadProject/Jure-assignments-service/config-yaml/dev.yaml"
const devConfigPath = "config-yaml/dev.yaml"

// @title Jury-assigments Service API
// @version 1.0
// @description Документация к микросервису назначения жюри на этапы
// @host 172.16.1.39:8090
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Jury-assigments
func main() {
	cfg := config.GetConfig(devConfigPath)

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
