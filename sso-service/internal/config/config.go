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
	JwtConfig                 `yaml:"JWT_INFO"`
	RedisConfig               `yaml:"REDIS_INFO"`
	RabbitConfig              `yaml:"RABBIT_INFO"`
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

type JwtConfig struct {
	Key             string `yaml:"key"`
	AccessDuration  int    `yaml:"access_duration"`
	RefreshDuration int    `yaml:"refresh_duration"`
}

type RedisConfig struct {
	AddressRedisPath string `yaml:"address_path"`
}

type RabbitConfig struct {
	AddressRabbitPath string `yaml:"address_path"`
	QueueName         string `yaml:"queue_name"`
}

type AdditionalAddressesConfig struct {
	ReactVision            string `yaml:"react"`
	JureAssignmentsService string `yaml:"jure-assignments"`
	Event                  string `yaml:"event"`
	OlympiadResult         string `yaml:"olympiad"`
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

func MustConfig(configPath string) *Config {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatalf("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Can not reaf config file %s: %s", configPath, err)
	}
	return &config
}
