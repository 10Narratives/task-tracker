package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Environment constants define the runtime environment modes for the application.
const (
	// EnvLocal represents the local development environment.
	EnvLocal = "local"

	// EnvDev represents the development environment.
	EnvDev = "dev"

	// EnvProd represents the production environment.
	EnvProd = "prod"
)

// Config represents the application's configuration settings.
type Config struct {
	// Env specifies the environment in which the application is running (e.g local, dev, prod).
	Env string `yaml:"env" env-default:"local"`

	// HTTP contains the configuration settings for the HTTP server.
	HTTP HTTPServerConfig `yaml:"http_server"`
}

// HTTPServerConfig holds configuration values for the HTTP server.
type HTTPServerConfig struct {
	// Address specifies the host address on which the server will listen.
	Address string `yaml:"address" env-default:"localhost"`

	// Port defines the port number on which the server will run.
	Port string `yaml:"port" env-default:"8000"`

	// Timeout sets the request timeout duration.
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`

	// IdleTimeout defines the maximum amount of time to wait for the next request.
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`

	// FileServerPath defines path to web pages
	FileServerPath string `yaml:"file_server_path" env-default:"./web"`
}

// MustConfig loads the application configuration from a YAML file specified by the CONFIG_PATH environment variable.
//
// This function ensures that:
//   - The CONFIG_PATH environment variable is set.
//   - The configuration file exists at the specified path.
//   - The configuration file can be successfully read and parsed.
//
// If any of these conditions are not met, the function logs a fatal error and terminates the application.
//
// Returns:
//   - *Config: A pointer to the fully populated Config struct.
func MustConfig() *Config {
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
