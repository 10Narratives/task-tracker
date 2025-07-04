package sl

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogdiscard"
	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogpretty"
	"github.com/natefinch/lumberjack"
)

type LoggerOptions struct {
	level  slog.Level
	format string
	output string
}

func defaultOptions() *LoggerOptions {
	return &LoggerOptions{
		level:  slog.LevelError,
		format: "json",
		output: "stdout",
	}
}

type LoggerOption func(*LoggerOptions)

func WithLevel(level string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.level = parseLevel(level)
	}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("unsupported log level: " + level)
	}
}

func WithFormat(format string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.format = format
	}
}

func WithOutput(output string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.output = output
	}
}

func MustLogger(opts ...LoggerOption) *slog.Logger {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	output := createOutput(options)
	handler := createHandler(options.format, output, options.level)
	return slog.New(handler)
}

func createOutput(opts *LoggerOptions) io.Writer {
	if opts.output == "stdout" {
		return os.Stdout
	}

	dir := filepath.Dir(opts.output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	return &lumberjack.Logger{
		Filename:  opts.output,
		LocalTime: true,
	}
}

func createHandler(format string, output io.Writer, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch format {
	case "json":
		return slog.NewJSONHandler(output, opts)
	case "pretty":
		return slogpretty.NewPrettyLogger(&slogpretty.PrettyHandlerOptions{
			SlogOpts: opts,
		}, output).Handler()
	case "plain":
		return slog.NewTextHandler(output, opts)
	case "discard":
		return slogdiscard.NewDiscardLogger().Handler()
	default:
		panic("unsupported log format: " + format)
	}
}
