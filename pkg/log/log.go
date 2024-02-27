package log

import (
	"sync"

	"github.com/rs/zerolog"
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
	zerologLevel zerolog.Level = zerolog.ErrorLevel
	once         sync.Once
)

func initLogger() {
	// set log time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerologLevel)
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
