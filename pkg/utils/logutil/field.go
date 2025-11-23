package logutil

import "github.com/rs/zerolog"

type Field func(event *zerolog.Event)

func Str(key string, value string) Field {
	return func(event *zerolog.Event) {
		event.Str(key, value)
	}
}

func Int(key string, value int) Field {
	return func(event *zerolog.Event) {
		event.Int(key, value)
	}
}

func Bool(key string, value bool) Field {
	return func(event *zerolog.Event) {
		event.Bool(key, value)
	}
}
