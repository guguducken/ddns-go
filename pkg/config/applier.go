package config

import (
	"errors"

	"github.com/guguducken/ddns-go/pkg/provider"
)

type DNSApplier struct {
	provider provider.DNSProvider
	domains  map[string]struct{}
}

type DNSAppliers map[string]*DNSApplier

var (
	ErrAllApplierFailed = errors.New("all provider and all domains report error")
)

func (das DNSAppliers) Apply(ip string) []error {
	return nil
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
