package plain

import (
	"net/http"

	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RequestURL        string            `yaml:"request_url"`
	AdditionalParams  map[string]string `yaml:"additional_params"`
	AdditionalHeaders map[string]string `yaml:"additional_headers"`

	headers          http.Header
	targetRequestURL string
}

func (config *Config) init(node yaml.Node, isV4 bool) error {
	if err := node.Decode(config); err != nil {
		return err
	}
	if config.RequestURL == "" {
		if isV4 {
			config.RequestURL = DefaultIPInfoAPIEndpoint
		} else {
			config.RequestURL = DefaultIPInfoV6APIEndpoint
		}
	}

	// internal config
	config.targetRequestURL = config.RequestURL

	if params := requestutils.GenParams(config.AdditionalParams); params != "" {
		config.targetRequestURL += "?" + params
	}

	config.headers = make(http.Header, len(config.AdditionalHeaders))
	for k, v := range config.AdditionalHeaders {
		config.headers.Set(k, v)
	}
	return nil
}
