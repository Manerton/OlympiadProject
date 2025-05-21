package main

import "main/internal/config"

const configPath = ""

func main() {

	cfg := config.MustConfig(configPath)

	_ = cfg
}
