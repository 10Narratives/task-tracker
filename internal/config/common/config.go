package commoncfg

type LoggerConfig struct {
	Level  string `yaml:"level" env-default:"error"` // Log verbosity level (debug, info, warn, error)
	Format string `yaml:"format" env-default:"json"` // Log output format (text, json)
	Output string `yaml:"output" env-default:"file"` // Log destination (stdout, stderr, file)
}
