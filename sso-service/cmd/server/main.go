package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/lib/liblogger"
)

const configPath = "config-yaml/local.yaml"

func main() {

	cfg := config.MustConfig(configPath)

	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start sso-service", slog.String("env", cfg.Env))
	log.Debug("debug message are enable")

}
