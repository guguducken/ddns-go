package utils

import (
	"os"
	"sync"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
)

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logutil.Fatal(
			nil,
			"environment variable is not found",
			logutil.Field{Key: "name", Value: key},
		)
	}
	return value
}

func Must(fn func() error) {
	if err := fn(); err != nil {
		logutil.Fatal(
			err,
			"failed to execute function",
		)
	}
}

var uniqueDomainMap = sync.Map{}

func CheckUniqueDomain(domain string) bool {
	if _, ok := uniqueDomainMap.Load(domain); ok {
		return false
	}
	uniqueDomainMap.Store(domain, struct{}{})
	return true
}
