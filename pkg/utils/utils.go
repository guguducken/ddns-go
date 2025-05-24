package utils

import (
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
)

func MustGetEnv(key string) string {
	// 检查环境变量名称的合法性
	validEnvName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !validEnvName.MatchString(key) {
		logutil.Fatal(
			nil,
			"invalid environment variable name",
			logutil.Field{Key: "name", Value: key},
		)
	}

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

func ParseEnv(input string) string {
	// 匹配条件：
	// 1. 以 ${{ 开头
	// 2. 以 }} 结尾
	// 3. 内容不能为空，或只包含空白字符
	reg := regexp.MustCompile(`\${{([^}]+?\S[^}]*?)}}`)
	return reg.ReplaceAllStringFunc(input, func(match string) string {
		return MustGetEnv(strings.TrimSpace(reg.FindStringSubmatch(match)[1]))
	})
}

var uniqueDomainMap = sync.Map{}

func CheckUniqueDomain(domain string) bool {
	if _, ok := uniqueDomainMap.Load(domain); ok {
		return false
	}
	uniqueDomainMap.Store(domain, struct{}{})
	return true
}

func MustWriteBytesTo(out io.Writer, message []byte) {
	if _, err := out.Write(message); err != nil {
		logutil.Fatal(err, "failed to write to output")
	}
}

func MustWriteStringTo(out io.Writer, message string) {
	if _, err := out.Write(UnsafeToByteSlice(message)); err != nil {
		logutil.Fatal(err, "failed write message to output")
	}
}

// UnsafeToString warp unsafe package functions
//
// NOTE: must not modify result string
func UnsafeToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// UnsafeToByteSlice warp unsafe package functions
//
// NOTE: must not modify result []byte
func UnsafeToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
