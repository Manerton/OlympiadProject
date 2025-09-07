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
	JwtTemp                   `yaml:"JWT_TEMP_INFO"`
	RedisConfig               `yaml:"REDIS_INFO"`
	RabbitConfig              `yaml:"RABBIT_INFO"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	SslModel string `yaml:"ssl_model"`
}

type AdditionalAddressesConfig struct {
	ReactVision string `yaml:"react"`
}

type HTTPServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type JwtTemp struct {
	Key string `yaml:"key"`
}

type RabbitConfig struct {
	AddressRabbitPath string `yaml:"address_path"`
	QueueName         string `yaml:"queue_name"`
}

type RedisConfig struct {
	AddressRedisPath string `yaml:"address_path"`
}

func (cfg *Config) GetDataSourceName() string {
	// dsn := "host=localhost user=user dbname=db password=password sslmode=disable"
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.SslModel,
	)
}

func (cfg *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", cfg.HTTPServerConfig.Host, cfg.HTTPServerConfig.Port)
}

func GetConfig(configPath string) *Config {
	envConfigPath := os.Getenv("CONFIG_PATH")
	if envConfigPath == "" {
		envConfigPath = configPath
		if envConfigPath == "" {
			log.Fatalf("CONFIG PATH is not set")
		}
	}

	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", envConfigPath)
	}

	var myconfig Config
	if err := cleanenv.ReadConfig(envConfigPath, &myconfig); err != nil {
		log.Fatalf("Can not read config file %s: %s", envConfigPath, err)
	}

	return &myconfig
}
