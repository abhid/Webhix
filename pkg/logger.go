package pkg

import (
	"log/slog"
	"os"
	"strings"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func init() {
	NewLogger()
}

func NewLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("ENV")
	}

	return NewLoggerWithEnv(env)
}

func NewLoggerWithEnv(env string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler

	switch normalizeEnv(env) {
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case EnvDev:
		opts.Level = slog.LevelDebug
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		opts.Level = slog.LevelDebug
		opts.AddSource = true
		handler = newPrettyHandler(os.Stdout, opts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}

func normalizeEnv(env string) string {
	env = strings.ToLower(strings.TrimSpace(env))
	if env == "" {
		return EnvLocal
	}

	switch env {
	case "development":
		return EnvDev
	case "production":
		return EnvProd
	default:
		return env
	}
}
