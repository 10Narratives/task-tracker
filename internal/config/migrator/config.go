package migratorcfg

import (
	"github.com/10Narratives/task-tracker/internal/config"
	commoncfg "github.com/10Narratives/task-tracker/internal/config/common"
)

type MigrateConfig struct {
	StoragePath     string `yaml:"storage_path" env-required:"true"`
	MigrationsPath  string `yaml:"migrations_path" env-required:"true"`
	MigrationsTable string `yaml:"migrations_table"`
}

type Config struct {
	Migrate MigrateConfig          `yaml:"migrate"`
	Logger  commoncfg.LoggerConfig `yaml:"logging"`
}

var loader = config.ConfigLoader[Config]{}

func MustLoad() *Config {
	return loader.MustLoad()
}

func MustLoadFromFile(path string) *Config {
	return loader.MustLoadFromFile(path)
}
