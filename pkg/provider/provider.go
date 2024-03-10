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
	Do(ip string) error
	FillUpDefaultValue()
	GetType() string
	CheckPermission() error
	GetDNSRecord(domain string, subDomain string) (config.DNSRecord, error)
	ListDNSRecords(domain string) (config.DNSRecords, error)
	CreateDNSRecord(record config.DNSRecord) (id uint64, err error)
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
		p := cfg.Providers[i]
		switch p.Type {
		case DNSPodProvider:
			providers = append(providers, NewDNSPodProvider(p.AccessKey, p.SecretKey, p.Domains))
		default:
			log.Error().Err(ErrUnsupportedProvider)
		}
	}
	dnsProviders = providers
	return dnsProviders
}
