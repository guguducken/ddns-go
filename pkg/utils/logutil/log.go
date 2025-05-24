package logutil

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/guguducken/ddns-go/pkg/errno"
)

var (
	logOutputs zerolog.LevelWriter
)

func Init(logLevel string, outputs []io.Writer) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(parseLogLevel(logLevel))
	if len(outputs) == 0 {
		outputs = []io.Writer{os.Stdout}
	}
	logOutputs = zerolog.MultiLevelWriter(outputs...)
}

func Debug(message string, fields ...Field) {
	event := GetSkipOneLogger().Debug().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	event.Msg(message)
}

func Info(message string, fields ...Field) {
	event := GetSkipOneLogger().Info().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	event.Msg(message)
}

func Warn(message string, fields ...Field) {
	event := GetSkipOneLogger().Warn().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	event.Msg(message)
}

func Error(err error, message string, fields ...Field) {
	event := GetSkipOneLogger().Error().Stack().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	AddFieldFromErrorAdditionalInfo(event, err)
	event.Err(err).Msg(message)
}

func Panic(err error, message string, fields ...Field) {
	event := GetSkipOneLogger().Panic().Stack().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	AddFieldFromErrorAdditionalInfo(event, err)
	event.Err(err).Msg(message)
}

func Fatal(err error, message string, fields ...Field) {
	event := GetSkipOneLogger().Fatal().Stack().Timestamp()
	for _, field := range fields {
		event.Str(field.Key, field.Value)
	}
	AddFieldFromErrorAdditionalInfo(event, err)
	event.Err(err).Msg(message)
}

func AddFieldFromErrorAdditionalInfo(event *zerolog.Event, err error) {
	additionalInfo := errno.GetAdditionalInfo(err)
	if additionalInfo == nil {
		return
	}
	for key, value := range additionalInfo {
		event.Str(key, value)
	}
}
