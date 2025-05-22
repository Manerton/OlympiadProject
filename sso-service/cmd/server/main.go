package main

import "main/internal/config"

const configPath = "config-yaml/local.yaml"

func main() {

	cfg := config.MustConfig(configPath)

	_ = cfg
}
