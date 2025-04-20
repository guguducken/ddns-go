package poolutils

import (
	"strings"
	"sync"
)

var stringBuilderPool *sync.Pool

func GetStringBuilder() *strings.Builder {
	return stringBuilderPool.Get().(*strings.Builder)
}

func PutStringBuilder(b *strings.Builder) {
	b.Reset()
	stringBuilderPool.Put(b)
}
