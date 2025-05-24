package provider

import (
	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/provider/json"
	"github.com/guguducken/ddns-go/pkg/provider/plain"
	"github.com/guguducken/ddns-go/pkg/utils/logutil"
)

type Provider interface {
	GetType() cons.ProviderType
	GetName() string
	ProviderIP() (string, error)
}

func NewProvider(t cons.ProviderType, config yaml.Node, name string, isV4 bool) (Provider, error) {
	switch t {
	case cons.ProviderTypePlain:
		return plain.NewProvider(config, name, isV4)
	case cons.ProviderTypeJson:
		return json.NewProvider(config, name, isV4)
	default:
		return nil, errno.OverrideError(
			errno.ErrInvalidProviderType,
			errno.AppendAdditionalMessage("ProviderType", string(t)),
		)
	}
}

type Providers []Provider

func (p Providers) ProviderIP() (string, error) {
	for _, p := range p {
		ip, err := p.ProviderIP()
		if err != nil {
			logutil.Error(err, "failed to provide ip", logutil.NewField("provider", p.GetName()))
			continue
		}
		return ip, nil
	}
	return "", errno.ErrCanNotProvideIP
}
