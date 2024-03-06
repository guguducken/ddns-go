package provider

import "net"

const (
	IPv4 = "A"
	IPv6 = "AAAA"
)

func CheckIPDNSType(ip string) string {
	i := net.ParseIP(ip)
	if i.To4() != nil {
		return IPv4
	}
	return IPv6
}
