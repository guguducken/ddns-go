package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/guguducken/ddns-go/pkg/ipcheck"
	"github.com/guguducken/ddns-go/pkg/log"
	"github.com/guguducken/ddns-go/pkg/provider"
	"gopkg.in/yaml.v3"
)

type Config struct {
	IPGettersInput []IPGettersConfig `yaml:"ip_getters,omitempty" json:"ip-getters,omitempty"`
	CheckInterval  *time.Duration    `yaml:"check_interval,omitempty" json:"check-interval,omitempty"`
	Providers      []ProvidersConfig `yaml:"providers,omitempty" json:"providers,omitempty"`
	LogLevel       string            `yaml:"log_level,omitempty" json:"log-level,omitempty"`

	IPGetters   ipcheck.IPGetters
	DnsAppliers map[string]*DNSApplier
}

type IPGettersConfig struct {
	Type  string `yaml:"type,omitempty" json:"type,omitempty"`
	URL   string `yaml:"url,omitempty" json:"url,omitempty"`
	Token string `yaml:"token,omitempty" json:"token,omitempty"`
}

type ProvidersConfig struct {
	Type      string   `yaml:"type,omitempty" json:"type,omitempty"`
	AccessKey string   `yaml:"access_key,omitempty" json:"access-key,omitempty"`
	SecretKey string   `yaml:"secret_key,omitempty" json:"secret-key,omitempty"`
	Domains   []string `yaml:"domains,omitempty" json:"domains,omitempty"`
}

func (cfg Config) AddApplier(applierType string, provider provider.DNSProvider, domains []string) {
	if _, ok := cfg.DnsAppliers[applierType]; !ok {
		cfg.DnsAppliers[applierType] = NewDNSApplier(provider, domains)
		return
	}

	cfg.DnsAppliers[applierType].AddDomains(domains)
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
	log.Init(cfg.LogLevel)
	log.Debug(fmt.Sprintf("config file path is: %s", path))
	log.Debug(fmt.Sprintf("config file content is: %s", string(content)))

	// parse to dns applier
	log.Debug("start init dns appliers")
	cfg.DnsAppliers = make(map[string]*DNSApplier, 20)
	for _, p := range cfg.Providers {
		switch p.Type {
		case provider.DNSPodProvider:
			log.Debug(fmt.Sprintf("add one dnspod applier to config"))
			cfg.AddApplier(provider.DNSPodProvider, provider.NewDNSPodProvider(p.AccessKey, p.SecretKey), p.Domains)
		default:
			err = errors.Join(provider.ErrUnsupportedProvider, errors.New(fmt.Sprintf("invalid dns provider is: %s", p.Type)))
			return nil, err
		}
	}

	// parse ip_getters input to ipcheck.IPGetters
	log.Debug("start init ip getters")
	cfg.IPGetters = make(ipcheck.IPGetters, 0, 10)
	for _, getter := range cfg.IPGettersInput {
		switch getter.Type {
		case ipcheck.HttpbinGetter:
			log.Debug(fmt.Sprintf("add one httpbin style ip_getter to config"))
			cfg.IPGetters = append(cfg.IPGetters, ipcheck.NewHttpbinGetter(getter.URL, getter.Token))
		case ipcheck.IpInfoGetter:
			log.Debug(fmt.Sprintf("add one ipinfo style ip_getter to config"))
			cfg.IPGetters = append(cfg.IPGetters, ipcheck.NewIPInfoGetter(getter.URL, getter.Token))
		default:
			err = errors.Join(ipcheck.ErrUnsupportedIPGetter, errors.New(fmt.Sprintf("invalid ip getter is: %s", getter.Type)))
			return nil, err
		}
	}
	return &cfg, nil
}
