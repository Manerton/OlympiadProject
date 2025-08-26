package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Route struct {
	Prefix    string   `yaml:"prefix"`
	Target    string   `yaml:"target"`
	SkipAuth  bool     `yaml:"skip_auth"`
	Aggregate bool     `yaml:"aggregate"`
	Services  []string `yaml:"services"` // список URL-ов для агрегации
	Roles     []int    `yaml:"roles"`    // ← ДОБАВИЛИ
}

type HTTPServers struct {
	Routes []Route `yaml:"routes"`
}

type HTTPServerMainConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

func (s *HTTPServerMainConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", s.Address, s.Port)
}

type JwtTemp struct {
	Key string `yaml:"key"`
}

type RabbitConfig struct {
	AddressPath string `yaml:"address_path"`
	QueueName   string `yaml:"queue_name"`
}

type RedisConfig struct {
	AddressPath string `yaml:"address_path"`
}

type AdditionalAddressesConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
}

type Config struct {
	Env                       string               `yaml:"env"`
	HTTPServers               HTTPServers          `yaml:"HTTP_SERVERS_INFO"`
	HTTPServerMain            HTTPServerMainConfig `yaml:"HTTP_SERVER_MAIN_INFO"`
	JwtTemp                   JwtTemp              `yaml:"JWT_TEMP_INFO"`
	RedisInfo                 RedisConfig          `yaml:"REDIS_INFO,omitempty"`
	RabbitInfo                RabbitConfig         `yaml:"RABBIT_INFO,omitempty"`
	AdditionalAddressesConfig `yaml:"CORS_INFO"`
}

func GetConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
