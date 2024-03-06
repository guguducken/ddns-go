package provider

import (
	"errors"
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
	InitDNSRecord(domain, subDomain, value string) DNSRecord
	CheckPermission() error
	GetDNSRecord(domain string, subDomain string) (DNSRecord, error)
	ListDNSRecords(domain string) (DNSRecords, error)
	CreateDNSRecord(domain string, record DNSRecord) error
	UpdateDNSRecord(domain string, record DNSRecord) error
	DeleteDNSRecord(domain string, record DNSRecord) error
}

type DNSRecord struct {
	Domain     string
	Name       string
	Value      string
	Type       string
	Status     string
	Line       string
	Weight     uint64
	Remark     string
	TTL        uint64
	MX         uint64
	UpdateTime string
}

type DNSRecords []DNSRecord

type DNSProviders []DNSProvider
