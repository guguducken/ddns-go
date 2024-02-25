package provider

import (
	"errors"
)

const (
	DNSPodProvider = "dnsPod"
	AliDNSProvider = "aliDNS"
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
	GetDNSRecord(domain string, subDomain string) (DNSRecord, error)
	ListDNSRecords(domain string) (DNSRecords, error)
	CreateDNSRecord(record DNSRecord) error
	UpdateDNSRecord(record DNSRecord) error
	DeleteDNSRecord(record DNSRecord) error
}

type DNSRecord struct {
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
