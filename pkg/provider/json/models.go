package json

import (
	"errors"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
)

type Config struct {
	RequestURL        string            `yaml:"request_url"`
	AdditionalParams  map[string]string `yaml:"additional_params"`
	AdditionalHeaders map[string]string `yaml:"additional_headers"`
	Path              string            `yaml:"path"`

	splitPath        []string
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
		config.Path = DefaultJsonPath
	}
	if config.Path == "" {
		return errors.New("path is required")
	}

	// internal config
	config.splitPath = strings.Split(config.Path, ".")
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
