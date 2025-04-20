package poolutils

import (
	"strings"
	"sync"
)

func init() {
	stringBuilderPool = &sync.Pool{
		New: func() any {
			builder := new(strings.Builder)
			builder.Grow(100)
			return builder
		},
	}

	stringChanPool = &sync.Pool{
		New: func() any {
			return make(chan string, 1)
		},
	}
}
