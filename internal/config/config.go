package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string           `yaml:"env" env-default:"local"`
	HTTP HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set.")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist: %s.", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Can not read configuration by %s", configPath)
	}

	return &cfg
}
