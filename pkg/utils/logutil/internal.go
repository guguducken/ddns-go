package logutil

import (
	"github.com/rs/zerolog"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

var (
	skipOneLogger zerolog.Logger
)

func initSkipOneLogger(writers zerolog.LevelWriter) {
	skipOneLogger = zerolog.New(writers).With().CallerWithSkipFrameCount(3).Logger()
}

func GetSkipOneLogger() *zerolog.Logger {
	return &skipOneLogger
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
