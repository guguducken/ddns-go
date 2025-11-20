package json

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/utils/requestutils"
)

type Config struct {
	RequestURL        string            `yaml:"request_url"`
	AdditionalParams  map[string]string `yaml:"additional_params"`
	AdditionalHeaders map[string]string `yaml:"additional_headers"`
	Path              string            `yaml:"path"`
	SuccessStatusCode int               `yaml:"success_status_code"`

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
		config.SuccessStatusCode = DefaultSuccessStatusCode
	}
	if config.Path == "" {
		return errors.New("path is required")
	}
	
	if config.SuccessStatusCode == 0 {
		logutil.Warn(fmt.Sprintf("success_status_code not set, will use default %d", DefaultSuccessStatusCode))
		config.SuccessStatusCode = DefaultSuccessStatusCode
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
