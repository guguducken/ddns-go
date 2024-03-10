package config

import (
	"errors"
	"fmt"
)

type DNSRecord struct {
	ID         uint64
	Type       string `yaml:"type,omitempty" json:"type,omitempty"`
	Domain     string `yaml:"domain,omitempty" json:"domain,omitempty"`
	SubDomain  string `yaml:"sub_domain,omitempty" json:"sub-domain,omitempty"`
	Status     string `yaml:"status,omitempty" json:"status,omitempty"`
	Line       string `yaml:"line,omitempty" json:"line,omitempty"`
	Weight     uint64 `yaml:"weight,omitempty" json:"weight,omitempty"`
	Remark     string `yaml:"remark,omitempty" json:"remark,omitempty"`
	TTL        uint64 `yaml:"ttl,omitempty" json:"ttl,omitempty"`
	MX         uint64 `yaml:"mx,omitempty" json:"mx,omitempty"`
	Value      string
	UpdateTime string
}

type DNSRecords []DNSRecord

func (ds DNSRecords) Validate() error {
	for _, d := range ds {
		err := d.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DNSRecord) Validate() error {
	errTpl := "%s of dns record must not be empty"

	if d.Type == "" {
		return errors.New(fmt.Sprintf(errTpl, "type"))
	}
	if d.Domain == "" {
		return errors.New(fmt.Sprintf(errTpl, "domain"))
	}
	if d.SubDomain == "" {
		return errors.New(fmt.Sprintf(errTpl, "subdomain"))
	}
	if d.Value == "" {
		return errors.New(fmt.Sprintf(errTpl, "value"))
	}
	if d.Type == "" {
		return errors.New(fmt.Sprintf(errTpl, "type"))
	}
	if d.Line == "" {
		return errors.New(fmt.Sprintf(errTpl, "line"))
	}
	if d.TTL == 0 {
		return errors.New(fmt.Sprintf(errTpl, "ttl"))
	}
	return nil
}
