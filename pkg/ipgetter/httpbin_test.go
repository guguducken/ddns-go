package ipgetter

import (
	"fmt"
	"testing"
)

func TestHttpbin_GetIP(t *testing.T) {
	ipGetter := NewHttpbinGetter("https://httpbin.org/get?test=1", "")
	ip, err := ipGetter.GetIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ip: %v\n", ip)
}
