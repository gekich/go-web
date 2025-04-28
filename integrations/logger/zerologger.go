package logger

import (
	"os"

	intLogger "github.com/gekich/go-web/internal/logger"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	l zerolog.Logger
}

func New(env string) intLogger.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch env {
	case "prod":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "dev", "local":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	baseLogger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", "your-service-name").
		Str("env", env).
		Logger()

	return &ZeroLogger{l: baseLogger}
}

func (z *ZeroLogger) Info(msg string, fields ...intLogger.Field) {
	e := z.l.Info()
	for _, f := range fields {
		e = e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (z *ZeroLogger) Error(msg string, fields ...intLogger.Field) {
	e := z.l.Error()
	for _, f := range fields {
		e = e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (z *ZeroLogger) Debug(msg string, fields ...intLogger.Field) {
	e := z.l.Debug()
	for _, f := range fields {
		e = e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (z *ZeroLogger) Warn(msg string, fields ...intLogger.Field) {
	e := z.l.Warn()
	for _, f := range fields {
		e = e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}
