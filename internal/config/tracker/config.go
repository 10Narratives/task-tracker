package trackercfg

import (
	"time"

	"github.com/10Narratives/task-tracker/internal/config"
)

// Config represents the root configuration structure for the task tracker application.
// It contains nested configurations for storage, HTTP server, and logging components.
// Fields are loaded from YAML configuration files and can be overridden by environment variables.
type Config struct {
	Storage StorageConfig    `yaml:"storage"`     // Database storage configuration
	HTTP    HTTPServerConfig `yaml:"http_server"` // HTTP server configuration
	Logger  LoggerConfig     `yaml:"logging"`     // Logging system configuration
}

// StorageConfig defines parameters for database connection and operation.
type StorageConfig struct {
	DriverName      string `yaml:"driver" env-default:"sqlite3"`   // Database driver name (e.g., mysql, postgres)
	DataSourceName  string `yaml:"dsn" env-default:"scheduler.db"` // Data Source Name (connection string)
	PaginationLimit uint   `yaml:"limit" env-default:"10"`         // Maximum records per paginated response
}

// HTTPServerConfig contains settings for the HTTP web server.
type HTTPServerConfig struct {
	Address        string        `yaml:"address" env-default:"localhost"`      // IP address or hostname to bind to
	Port           string        `yaml:"port" env-default:"8000"`              // TCP port to listen on
	Timeout        time.Duration `yaml:"timeout" env-default:"4s"`             // Request timeout duration
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`       // Keep-alive connection timeout
	FileServerPath string        `yaml:"file_server_path" env-default:"./web"` // Path to static web assets directory
}

// LoggerConfig defines parameters for log output and formatting.
type LoggerConfig struct {
	Level  string `yaml:"level" env-default:"error"` // Log verbosity level (debug, info, warn, error)
	Format string `yaml:"format" env-default:"json"` // Log output format (text, json)
	Output string `yaml:"output" env-default:"file"` // Log destination (stdout, stderr, file)
}

var loader = config.ConfigLoader[Config]{}

// MustLoad loads configuration using the default loader instance.
// Terminates the application on any configuration error.
// Returns a pointer to the initialized Config struct.
func MustLoad() *Config {
	return loader.MustLoad()
}

// MustLoadFromFile loads configuration from a specific file path.
// Terminates the application on any configuration error.
// path: Absolute or relative path to the configuration file
// Returns a pointer to the initialized Config struct.
func MustLoadFromFile(path string) *Config {
	return loader.MustLoadFromFile(path)
}
