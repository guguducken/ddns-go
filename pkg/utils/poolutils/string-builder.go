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

func GenString(strs ...string) string {
	builder := GetStringBuilder()
	defer PutStringBuilder(builder)

	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}
