package ipgetter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/rs/zerolog/log"
)

var (
	ErrAllGetterFailed     = errors.New("all ip getter response fail")
	ErrUnsupportedIPGetter = errors.New("unsupported ipp getter")
	ErrInvalidResponseIP   = errors.New("ip getter response an invalid ip")
)

var once sync.Once

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

var ipGetters IPGetters

func InitIPGetters(cfg *config.Config) IPGetters {
	if ipGetters != nil {
		return ipGetters
	}
	getters := make(IPGetters, 0, 10)
	log.Debug().Msg("start init ip getters")
	for _, getter := range cfg.IPGettersInput {
		switch getter.Type {
		case HttpbinGetter:
			log.Debug().Msg(fmt.Sprintf("add one httpbin style ip_getter to config"))
			getters = append(getters, NewHttpbinGetter(getter.URL, getter.Token))
		case IpInfoGetter:
			log.Debug().Msg(fmt.Sprintf("add one ipinfo style ip_getter to config"))
			getters = append(getters, NewIPInfoGetter(getter.URL, getter.Token))
		default:
			log.Error().Err(errors.Join(ErrUnsupportedIPGetter, errors.New(fmt.Sprintf("invalid ip getter is: %s, so skip it", getter.Type))))
		}
	}
	ipGetters = getters
	return ipGetters
}

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
