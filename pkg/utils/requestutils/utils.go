package requestutils

import (
	"fmt"
	"strings"
)

func GenParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	p := make([]string, 0, len(params))
	for k, v := range params {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, "&")
}
