package cons

type RecorderType string

const (
	RecorderTypeDNSPod     RecorderType = "dnspod"
	RecorderTypeAliyun     RecorderType = "aliyun"
	RecorderTypeCloudflare RecorderType = "cloudflare"
)

func (p RecorderType) String() string {
	return string(p)
}

type RecordType string

const (
	RecordTypeA    RecordType = "A"
	RecordTypeAAAA RecordType = "AAAA"
)

type ProviderType string

func (p ProviderType) String() string {
	return string(p)
}

const (
	ProviderTypePlain ProviderType = "plain"
	ProviderTypeJson  ProviderType = "json"
)

var (
	SupportedProviders = []ProviderType{
		ProviderTypePlain,
		ProviderTypeJson,
	}
	SupportedRecordTypes = []RecordType{
		RecordTypeA,
		RecordTypeAAAA,
	}
	SupportedRecorders = []RecorderType{
		RecorderTypeDNSPod,
		// RecorderTypeAliyun,
		// RecorderTypeCloudflare,
	}
)
