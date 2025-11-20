package provider

import (
	"context"
	"errors"

	"github.com/guguducken/ddns-go/pkg/utils"
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

func (ps Providers) ProviderIP(ctx context.Context) (string, error) {
	for _, p := range ps {
		ip, err := utils.RunWithContext(ctx, p.ProviderIP)
		if err != nil {
			// if error is context related error, need direct return empty ip and err
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return ip, err
			}
			logutil.Error(
				err,
				"failed to provide ip",
				logutil.NewField("provider", p.GetName()),
				logutil.NewField("type", p.GetType().String()),
			)
			continue
		}
		return ip, nil
	}
	return "", errno.ErrCanNotProvideIP
}
