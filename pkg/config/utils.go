package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func MustGetEnv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic(fmt.Sprintf("can not found env %s", name))
	}
	log.Error().Msgf("get input %s from env", name)
	return env
}
