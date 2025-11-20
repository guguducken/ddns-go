package json

import (
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/buger/jsonparser"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/utils/iputil"
	"github.com/guguducken/ddns-go/pkg/utils/poolutils"
	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
)

type Provider struct {
	config *Config
	name   string
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
	return cons.ProviderTypeJson
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

	if response.StatusCode() != p.config.SuccessStatusCode {
		return "", errno.OverrideError(
			errno.ErrFailedProvideIP,
			errno.AppendAdditionalMessage("InvalidStatusCode", strconv.Itoa(response.StatusCode())),
		)
	}

	ip, err := jsonparser.GetString(response.Body(), p.config.splitPath...)
	if err != nil {
		return "", errno.OverrideError(
			errno.ErrFailedProvideIP,
			errno.AppendAdditionalMessage("GetJsonIPError", err.Error()),
		)
	}
	if _, err := iputil.CheckIPType(ip); err != nil {
		return "", err
	}
	return ip, nil
}
