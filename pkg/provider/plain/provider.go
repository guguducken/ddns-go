package plain

import (
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/utils/iputil"
	"github.com/guguducken/ddns-go/pkg/utils/poolutils"
	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
)

type Provider struct {
	config *Config

	name string
}

func NewProvider(config yaml.Node, name string, isV4 bool) (*Provider, error) {
	p := &Provider{
		name:   name,
		config: &Config{},
	}

	if err := p.config.init(config, isV4); err != nil {
		return nil, err
	}
	if len(p.name) == 0 {
		p.name = poolutils.GenString(uuid.Must(uuid.NewV7()).String(), "(", p.config.RequestURL, ")")
	}
	return p, nil
}

func (p *Provider) GetType() cons.ProviderType {
	return cons.ProviderTypePlain
}

func (p *Provider) GetName() string {
	return p.name
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
