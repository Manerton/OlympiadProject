package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                       string `yaml:"env"`
	DatabaseConfig            `yaml:"DB_INFO"`
	HTTPServerConfig          `yaml:"HTTP_SERVER_INFO"`
	AdditionalAddressesConfig `yaml:"ADDITIONAL_ADDRESSES"`
	JWTConfig                 `yaml:"JWT_TEMP_INFO"`
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

type JWTConfig struct {
	Key string `yaml:"key"`
}

type AdditionalAddressesConfig struct {
	ReactVision  string `yaml:"react"`
	EventService string `yaml:"event"`
	JuryService  string `yaml:"jury"`
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

func GetConfig(configPath string) *Config {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatalf("CONFIG PATH is not set")
		}
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
