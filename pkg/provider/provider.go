package provider

import (
	"errors"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/rs/zerolog/log"
)

const (
	DNSPodProvider = "dnspod"
	AliDNSProvider = "alidns"
)

var (
	AllowedProviderTypes = map[string]string{
		DNSPodProvider: DNSPodProvider,
		AliDNSProvider: AliDNSProvider,
	}
)

var (
	ErrUnsupportedProvider = errors.New("unsupported dns provider")
	ErrNoDNSRecord         = errors.New("no such dns record")
	ErrPermissionInvalid   = errors.New("user permission is invalid")
)

type DNSProvider interface {
	GetType() string
	CheckPermission() error
	GetDNSRecord(domain string, subDomain string) (config.DNSRecord, error)
	ListDNSRecords(domain string) (config.DNSRecords, error)
	CreateDNSRecord(record config.DNSRecord) error
	UpdateDNSRecord(record config.DNSRecord) error
	DeleteDNSRecord(record config.DNSRecord) error
}

type DNSProviders []DNSProvider

var dnsProviders DNSProviders

func InitDNSProviders(cfg *config.Config) DNSProviders {
	if dnsProviders != nil {
		return dnsProviders
	}
	providers := make(DNSProviders, 0, 10)
	for i := 0; i < len(cfg.Providers); i++ {
		provider := cfg.Providers[i]
		switch provider.Type {
		case DNSPodProvider:
			providers = append(providers, NewDNSPodProvider(provider.AccessKey, provider.SecretKey))
		default:
			log.Error().Err(ErrUnsupportedProvider)
		}
	}
	dnsProviders = providers
	return dnsProviders
}
