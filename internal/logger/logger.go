package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/config"
)

// Init sets up the global logger
func Init() {
	// Pretty console output in development
	if config.C.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// JSON output in production
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().
		Str("env", config.C.Env).
		Msg("📝 Logger initialized")
}

// Convenience helpers (shortcuts to global logger)
func Info(msg string) *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}