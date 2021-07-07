package logconfig

import (
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	EnvLogMode = "MYSTERIUM_LOG_MODE"
	ModeJSON   = "json"
)

type Config struct {
	LogFilePath string
	LogFileName string
}

func BootstrapDefaultLogger() *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	writers := []io.Writer{}

	if isJSONMode() {
		writers = append(writers, os.Stderr)
	} else {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := log.Output(zerolog.MultiLevelWriter(writers...)).
		Level(zerolog.DebugLevel).
		With().
		Caller().
		Timestamp().
		Logger()

	setGlobalLogger(&logger)

	return &logger
}

func setGlobalLogger(logger *zerolog.Logger) {
	log.Logger = *logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

func isJSONMode() bool {
	v, ok := os.LookupEnv(EnvLogMode)
	if !ok {
		return false
	}
	return v == ModeJSON
}
