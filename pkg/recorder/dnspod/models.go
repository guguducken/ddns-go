package dnspod

import (
	"errors"

	"github.com/guguducken/ddns-go/pkg/cons"
	"gopkg.in/yaml.v3"
)

const (
	DefaultDnsPodAPIEndpoint = "dnspod.tencentcloudapi.com"
)

type Config struct {
	DnsPodAPIEndpoint string   `yaml:"dnspod_api_endpoint"`
	AccessKey         string   `yaml:"access_key"`
	SecretKey         string   `yaml:"secret_key"`
	Domains           []Domain `yaml:"domains"`
}

type Domain struct {
	Domain     string   `yaml:"domain"`
	SubDomains []string `yaml:"sub_domains"`
	RecordLine string   `yaml:"record_line"`
	TTL        uint64   `yaml:"ttl"`
	// RecordType cons.RecordType `yaml:"record_type"`
	Weight uint64 `yaml:"weight"`
	Status string `yaml:"status"`
	Remark string `yaml:"remark"`
}

type domainDetail struct {
	// id is the id of the record which is returned by the dnspod,
	// if the record is not found, it is 0
	id         uint64
	domain     string
	subDomain  string
	recordLine string
	recordType cons.RecordType
	ttl        uint64
	weight     uint64
	status     string
	remark     string
	value      string
}

func (d domainDetail) Copy() domainDetail {
	return domainDetail{
		id:         d.id,
		domain:     d.domain,
		subDomain:  d.subDomain,
		recordLine: d.recordLine,
		recordType: d.recordType,
		ttl:        d.ttl,
		weight:     d.weight,
		status:     d.status,
		remark:     d.remark,
		value:      d.value,
	}
}

func (c *Config) Init(data yaml.Node) error {
	if err := data.Decode(c); err != nil {
		return err
	}

	if c.AccessKey == "" || c.SecretKey == "" {
		return errors.New("access_key and secret_key are required")
	}
	if c.DnsPodAPIEndpoint == "" {
		c.DnsPodAPIEndpoint = DefaultDnsPodAPIEndpoint
	}
	for i := range c.Domains {
		if err := c.Domains[i].Validate(); err != nil {
			return err
		}
		setDomainDefault(&c.Domains[i])
	}
	return nil
}

func (d Domain) Validate() error {
	if d.TTL != 0 && (d.TTL < 1 || d.TTL > 604800) {
		return errors.New("ttl must be between 1 and 604800")
	}
	if d.Weight != 0 && (d.Weight < 1 || d.Weight > 100) {
		return errors.New("weight must be between 1 and 100")
	}
	if d.Status != "" && d.Status != "ENABLE" && d.Status != "DISABLE" {
		return errors.New("status must be ENABLE or DISABLE")
	}

	return nil
}

func setDomainDefault(domain *Domain) {
	if domain.RecordLine == "" {
		domain.RecordLine = DefaultRecordLine
	}
	if domain.TTL == 0 {
		domain.TTL = DefaultTTL
	}
	if domain.Weight == 0 {
		domain.Weight = DefaultWeight
	}
	if domain.Status == "" {
		domain.Status = DefaultStatus
	}
	if domain.Remark == "" {
		domain.Remark = DefaultRemark
	}
}
