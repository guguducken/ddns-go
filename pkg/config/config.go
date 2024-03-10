package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	dlog "github.com/guguducken/ddns-go/pkg/log"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyProviderDomains = errors.New("empty provider domains")
)

const (
	defaultCheckInterval int64 = 10

	EnvInputPrefix = "ENV_"
)

type Config struct {
	IPGettersInput []IPGettersConfig `yaml:"ip_getters,omitempty" json:"ip-getters,omitempty"`
	CheckInterval  int64             `yaml:"check_interval,omitempty" json:"check-interval,omitempty"`
	Providers      []ProvidersConfig `yaml:"providers,omitempty" json:"providers,omitempty"`
	LogLevel       string            `yaml:"log_level,omitempty" json:"log-level,omitempty"`
	Type           string            `yaml:"type,omitempty" json:"type,omitempty"`
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

	// validate the provider's domains
	for _, pro := range cfg.Providers {
		if err = pro.Domains.Validate(); err != nil {
			return nil, err
		}
	}

	InitEnvInputs(&cfg)

	return &cfg, nil
}

func InitEnvInputs(cfg *Config) {
	// get provider config from env
	for i := 0; i < len(cfg.Providers); i++ {
		if strings.HasPrefix(cfg.Providers[i].Type, EnvInputPrefix) {
			cfg.Providers[i].Type = MustGetEnv(cfg.Providers[i].Type[len(EnvInputPrefix):])
		}
		if strings.HasPrefix(cfg.Providers[i].AccessKey, EnvInputPrefix) {
			cfg.Providers[i].AccessKey = MustGetEnv(cfg.Providers[i].AccessKey[len(EnvInputPrefix):])
		}
		if strings.HasPrefix(cfg.Providers[i].SecretKey, EnvInputPrefix) {
			cfg.Providers[i].SecretKey = MustGetEnv(cfg.Providers[i].SecretKey[len(EnvInputPrefix):])
		}
	}
}
