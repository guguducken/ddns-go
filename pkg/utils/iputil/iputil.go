package iputil

import (
	"net"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
)

func CheckIPType(ip string) (cons.RecordType, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return "", errno.OverrideError(
			errno.ErrInvalidIPAddress,
			errno.AppendAdditionalMessage("Content", ip),
		)
	}
	if i.To4() != nil {
		return cons.RecordTypeA, nil
	}
	return cons.RecordTypeAAAA, nil
}
