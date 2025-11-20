package plain

import (
	"fmt"
	"net/http"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
)

type Config struct {
	RequestURL        string            `yaml:"request_url"`
	AdditionalParams  map[string]string `yaml:"additional_params"`
	AdditionalHeaders map[string]string `yaml:"additional_headers"`
	SuccessStatusCode int               `yaml:"success_status_code"`

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
	if config.SuccessStatusCode == 0 {
		logutil.Warn(fmt.Sprintf("success_status_code not set, will use default %d", DefaultSuccessStatusCode))
		config.SuccessStatusCode = DefaultSuccessStatusCode
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
