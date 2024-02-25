package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/guguducken/ddns-go/pkg/provider"
)

type Config struct {
	IPGetters     *IPGettersConfig  `yaml:"ip_getters,omitempty" json:"ip-getters,omitempty"`
	CheckInterval *time.Duration    `yaml:"check_interval,omitempty" json:"check-interval,omitempty"`
	Providers     []ProvidersConfig `yaml:"providers,omitempty" json:"providers,omitempty"`
	dnsAppliers   map[string]*DNSApplier
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
	if _, ok := cfg.dnsAppliers[applierType]; !ok {
		cfg.dnsAppliers[applierType] = NewDNSApplier(provider, domains)
		return
	}

	cfg.dnsAppliers[applierType].AddDomains(domains)
}

type DNSApplier struct {
	provider provider.DNSProvider
	domains  map[string]struct{}
}

func NewDNSApplier(provider provider.DNSProvider, domains []string) *DNSApplier {
	m := make(map[string]struct{})
	for i := 0; i < len(domains); i++ {
		m[domains[i]] = struct{}{}
	}
	return &DNSApplier{
		provider: provider,
		domains:  m,
	}
}

func (da *DNSApplier) Apply(ip string) []error {
	errs := make([]error, 0, len(da.domains))
	for i := 0; i < len(da.domains); i++ {

	}
	return errs
}

func (da *DNSApplier) AddDomains(domains []string) {
	for i := 0; i < len(domains); i++ {
		if _, ok := da.domains[domains[i]]; !ok {
			da.domains[domains[i]] = struct{}{}
		}
	}
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
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	// parse to dns applier
	for _, p := range cfg.Providers {
		switch p.Type {
		case provider.DNSPodProvider:
			cfg.AddApplier(provider.DNSPodProvider, provider.NewDNSPodProvider(p.AccessKey, p.SecretKey), p.Domains)
		default:
			err = errors.Join(provider.ErrUnsupportedProvider, errors.New(fmt.Sprintf("invalid dns provider is: %s", p.Type)))
			return nil, err
		}
	}
	return &cfg, nil
}
