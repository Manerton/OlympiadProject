package main

import (
	"main/internal/config"
	"main/internal/storage/postgresql"
)

const LocalFilePath = "config-yaml/local.yaml"

func main() {
	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger

	// Init storage
	storage, err := postgresql.NewPosgreSQL(cfg.GetDataSourceName())
	if err != nil {
		// use logger to error print
	}

	// for test run)
	_ = storage

}
