package ipcheck

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrAllGetterFailed     = errors.New("all ip getter response fail")
	ErrUnsupportedIPGetter = errors.New("unsupported ipp getter")
	ErrInvalidResponseIP   = errors.New("ip getter response an invalid ip")
)

const (
	HttpbinGetter = "httpbin"
	IpInfoGetter  = "ipinfo"
)

type IPGetter interface {
	GetIP() (ip string, err error)
	GetIPWithContext(ctx context.Context) (ip string, err error)
	GetType() string
	GetURL() string
	GetToken() string
}

type IPGetters []IPGetter

func (i IPGetters) GetIP() (ip string, err error) {
	return i.GetIPWithContext(context.Background())
}

func (i IPGetters) GetIPWithContext(ctx context.Context) (ip string, err error) {
	for _, getter := range i {
		ip, err = getter.GetIPWithContext(ctx)
		if err != nil {
			fmt.Printf("ip getter error: %s", getter.GetURL())
			continue
		}
		return ip, err
	}
	return "", ErrAllGetterFailed
}

type SimpleChecker struct {
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

func (s SimpleChecker) GetToken() string {
	return s.Token
}

func (s SimpleChecker) GetURL() string {
	return s.URL
}

func (s SimpleChecker) GetType() string {
	return s.Type
}
