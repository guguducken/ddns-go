package log

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warning"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

var (
	logger       zerolog.Logger
	zerologLevel zerolog.Level = zerolog.ErrorLevel
	once         sync.Once
)

func initLogger() {
	// set log time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = log.Level(zerologLevel).Output(os.Stdout)
}

func Init(logLevel string) {
	switch logLevel {
	case DebugLevel:
		zerologLevel = zerolog.DebugLevel
	case InfoLevel:
		zerologLevel = zerolog.InfoLevel
	case WarnLevel:
		zerologLevel = zerolog.WarnLevel
	case ErrorLevel:
		zerologLevel = zerolog.ErrorLevel
	case FatalLevel:
		zerologLevel = zerolog.FatalLevel
	case PanicLevel:
		zerologLevel = zerolog.PanicLevel
	default:
		zerologLevel = zerolog.ErrorLevel
	}
	once.Do(initLogger)
}

func Debug(message string) {
	logger.Debug().Msg(message)
}

func Info(message string) {
	logger.Info().Msg(message)
}

func Warn(message string) {
	logger.Warn().Msg(message)
}

func Error(err error) {
	logger.Error().Err(err).Stack()
}

func Fatal(message string) {
	logger.Fatal().Msg(message)
}

func Panic(message string) {
	logger.Panic().Msg(message)
}
