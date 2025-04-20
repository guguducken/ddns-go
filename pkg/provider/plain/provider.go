package plain

import (
	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/utils/iputil"
	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
	"gopkg.in/yaml.v3"
)

type Provider struct {
	config *Config
}

func NewProvider(config yaml.Node, isV4 bool) (*Provider, error) {
	p := &Provider{
		config: &Config{},
	}

	if err := p.config.init(config, isV4); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Provider) GetType() cons.ProviderType {
	return cons.ProviderTypePlain
}

func (p *Provider) ProviderIP() (string, error) {
	response, err := requestutils.Get(p.config.targetRequestURL, p.config.headers)
	if err != nil {
		return "", errno.OverrideError(
			errno.ErrFailedProvideIP,
			errno.AppendAdditionalMessage("RequestError", err.Error()),
		)
	}
	defer requestutils.ReleaseResponse(response)
	body := string(response.Body())
	if _, err := iputil.CheckIPType(body); err != nil {
		return "", err
	}
	return body, nil
}
