package ipgetter

import (
	"fmt"
	"testing"
)

func TestIPInfo_GetIP(t *testing.T) {
	ipGetter := NewIPInfoGetter("https://ipinfo.io/ip", "")
	ip, err := ipGetter.GetIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ip: %v\n", ip)
}
