package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string `yaml:"env"`
	DatabaseConfig   `yaml:"DB_INFO"`
	HTTPServerConfig `yaml:"HTTP_SERVER_INFO"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	SslModel string `yaml:"ssl_model"`
}

type HTTPServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func (cfg *Config) GetDataSourceName() string {
	// dsn := "host=localhost user=user dbname=db password=password sslmode=disable"
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.SslModel,
	)
}

func (cfg *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", cfg.HTTPServerConfig.Host, cfg.HTTPServerConfig.Port)
}

func GetConfig(config_path string) *Config {
	configPath := config_path
	if configPath == "" {
		log.Fatal("CONFIG PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var myconfig Config
	if err := cleanenv.ReadConfig(configPath, &myconfig); err != nil {
		log.Fatalf("Can not read config file %s: %s", configPath, err)
	}

	return &myconfig
}
