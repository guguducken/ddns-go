package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	dlog "github.com/guguducken/ddns-go/pkg/log"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyProviderDomains = errors.New("empty provider domains")
)

const (
	defaultCheckInterval int64 = 10
)

type Config struct {
	IPGettersInput []IPGettersConfig `yaml:"ip_getters,omitempty" json:"ip-getters,omitempty"`
	CheckInterval  int64             `yaml:"check_interval,omitempty" json:"check-interval,omitempty"`
	Providers      []ProvidersConfig `yaml:"providers,omitempty" json:"providers,omitempty"`
	LogLevel       string            `yaml:"log_level,omitempty" json:"log-level,omitempty"`
	Type           string            `yaml:"type,omitempty" json:"type,omitempty"`

	totalDomains int
}

type IPGettersConfig struct {
	Type  string `yaml:"type,omitempty" json:"type,omitempty"`
	URL   string `yaml:"url,omitempty" json:"url,omitempty"`
	Token string `yaml:"token,omitempty" json:"token,omitempty"`
}

type ProvidersConfig struct {
	Type      string     `yaml:"type,omitempty" json:"type,omitempty"`
	AccessKey string     `yaml:"access_key,omitempty" json:"access-key,omitempty"`
	SecretKey string     `yaml:"secret_key,omitempty" json:"secret-key,omitempty"`
	Domains   DNSRecords `yaml:"domains,omitempty" json:"domains,omitempty"`
}

func (cfg *Config) GetTotalDomains() int {
	return cfg.totalDomains
}

func NewConfig(path string) (*Config, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// read from config file
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// parse to config struct
	cfg := Config{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	// init logger
	dlog.Init(cfg.LogLevel)
	log.Debug().Msg(fmt.Sprintf("config file path is: %s", path))
	log.Debug().Msg(fmt.Sprintf("config file content is: %s", string(content)))

	// init run type
	if cfg.Type == "" {
		cfg.Type = "client"
	}

	// init check interval
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = defaultCheckInterval
	}

	// parse to dns applier
	log.Debug().Msg("start init dns appliers")
	for _, p := range cfg.Providers {

		// check the number of p.domains
		if len(p.Domains) == 0 {
			log.Error().Err(errors.Join(ErrEmptyProviderDomains, errors.New(fmt.Sprintf("provider is %s", p.Type))))
			continue
		}

		// calculate total domains
		cfg.totalDomains += len(p.Domains)

	}

	// must return err if no domains to create dns record
	if cfg.totalDomains == 0 {
		return nil, ErrEmptyProviderDomains
	}

	return &cfg, nil
}
