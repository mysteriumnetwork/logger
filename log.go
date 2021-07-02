package mlog

import (
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var DefaultConfig = Config{
	LogFilePath: "/var/log/containers/",
	LogFileName: "app",
}

type Config struct {
	LogFilePath string
	LogFileName string
}

func MustBootstrapDefaultLogger() *zerolog.Logger {
	logger, err := BootStrapLogger(DefaultConfig)
	if err != nil {
		panic(err)
	}
	return logger
}

func BootstrapDefaultLogger() (*zerolog.Logger, error) {
	return BootStrapLogger(DefaultConfig)
}

func BootStrapLogger(cfg Config) (*zerolog.Logger, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	writers := []io.Writer{}

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})

	fileWriter, err := NewRollingWriter(cfg.LogFilePath, cfg.LogFileName)
	if err == nil {
		writers = append(writers, fileWriter)
	} else {
		stdlog.Printf("Failed to create rolling file writer: %s", err)
	}

	logger := log.Output(zerolog.MultiLevelWriter(writers...)).
		Level(zerolog.DebugLevel).
		With().
		Caller().
		Timestamp().
		Logger()

	setGlobalLogger(&logger)

	return &logger, err
}

func setGlobalLogger(logger *zerolog.Logger) {
	log.Logger = *logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}
