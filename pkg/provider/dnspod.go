package provider

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/guguducken/ddns-go/pkg/config"
	tcam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

const (
	dnsPodAPIEndpoint     = `dnspod.tencentcloudapi.com`
	tencentCAMAPIEndpoint = "cam.tencentcloudapi.com"
)

type DNSPod struct {
	accessKey  string
	secretKey  string
	credential *tcommon.Credential
}

// NewDNSPodProvider will return nil pointer if accessKey is empty or secretKey is empty
func NewDNSPodProvider(accessKey string, secretKey string) *DNSPod {
	if accessKey == "" || secretKey == "" {
		return nil
	}
	return &DNSPod{
		accessKey:  accessKey,
		secretKey:  secretKey,
		credential: tcommon.NewCredential(accessKey, secretKey),
	}
}

func (d *DNSPod) NewDNSPodClient() (*dnspod.Client, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = dnsPodAPIEndpoint
	cpf.NetworkFailureMaxRetries = 10
	cpf.UnsafeRetryOnConnectionFailure = true
	return dnspod.NewClient(d.credential, "", cpf)
}

func (d *DNSPod) NewTencentCAMClient() (*tcam.Client, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = tencentCAMAPIEndpoint
	cpf.NetworkFailureMaxRetries = 10
	cpf.UnsafeRetryOnConnectionFailure = true
	return tcam.NewClient(d.credential, "", cpf)
}

func (d *DNSPod) GetType() string {
	return DNSPodProvider
}

func (d *DNSPod) CheckPermission() error {
	userUID, err := d.getUserUid()
	if err != nil {
		return err
	}
	allowedPermissions := []string{"AdministratorAccess", "QCloudResourceFullAccess", "QcloudDNSPodFullAccess", "CustomerDNSPodFullAccess"}

	request := tcam.NewListPoliciesGrantingServiceAccessRequest()
	request.TargetUin = tcommon.Uint64Ptr(userUID)

	client, err := d.NewTencentCAMClient()

	response, err := client.ListPoliciesGrantingServiceAccess(request)
	if err != nil {
		return err
	}

	checkDetailPermission := func(node *tcam.ListGrantServiceAccessNode) bool {
		if *node.Service.ServiceType != "dnspod" {
			return false
		}
		for _, policy := range node.Policy {
			if slices.Contains(allowedPermissions, *policy.PolicyName) {
				return true
			}
		}
		return false
	}

	for _, node := range response.Response.List {
		if checkDetailPermission(node) {
			return nil
		}
	}
	return ErrPermissionInvalid
}

func (d *DNSPod) GetDNSRecord(domain string, subDomain string) (config.DNSRecord, error) {
	var record config.DNSRecord

	// init new request
	request := dnspod.NewDescribeRecordFilterListRequest()

	request.Domain = tcommon.StringPtr(domain)
	request.SubDomain = tcommon.StringPtr(subDomain)
	request.IsExactSubDomain = tcommon.BoolPtr(true)
	request.Limit = tcommon.Uint64Ptr(100)

	// init client
	client, err := d.NewDNSPodClient()
	if err != nil {
		return record, err
	}

	response, err := client.DescribeRecordFilterList(request)
	if err != nil {
		return record, err
	}

	if *response.Response.RecordCountInfo.ListCount == 0 || len(response.Response.RecordList) == 0 {
		return record, errors.Join(ErrNoDNSRecord,
			errors.New(fmt.Sprintf("provider: %s, domain: %s, subDomain: %s", DNSPodProvider, domain, subDomain)))
	}

	// we can get only one record
	// parse to DNSRecord
	record = d.ParseToDNSRecord(domain, response.Response.RecordList[0])
	return record, nil
}

func (d *DNSPod) ListDNSRecords(domain string) (config.DNSRecords, error) {
	records := make(config.DNSRecords, 0, 30)

	var offset, limit uint64 = 0, 200
	for {
		recordsTemp, err := d.listDNSRecordByPage(domain, offset, limit)
		if err != nil {
			return nil, err
		}
		records = append(records, recordsTemp...)
		offset += limit
		if len(recordsTemp) < int(limit) {
			break
		}
	}

	return records, nil
}

func (d *DNSPod) listDNSRecordByPage(domain string, offset uint64, limit uint64) (config.DNSRecords, error) {
	records := make(config.DNSRecords, 0, 30)

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = tcommon.StringPtr(domain)
	request.Offset = tcommon.Uint64Ptr(offset)
	request.Limit = tcommon.Uint64Ptr(limit)

	client, err := d.NewDNSPodClient()
	if err != nil {
		return nil, err
	}

	response, err := client.DescribeRecordList(request)
	if err != nil {
		return nil, err
	}

	for _, item := range response.Response.RecordList {
		records = append(records, d.ParseToDNSRecord(domain, item))
	}

	return records, nil
}

func (d *DNSPod) CreateDNSRecord(record config.DNSRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}

	request := dnspod.NewCreateRecordRequest()
	request.Domain = tcommon.StringPtr(record.Domain)
	request.RecordType = tcommon.StringPtr(CheckIPDNSType(record.Value))
	request.RecordLine = tcommon.StringPtr(record.Line)
	request.Value = tcommon.StringPtr(record.Value)
	request.SubDomain = tcommon.StringPtr(record.SubDomain)
	request.TTL = tcommon.Uint64Ptr(record.TTL)
	request.Weight = tcommon.Uint64Ptr(record.Weight)
	request.Status = tcommon.StringPtr(record.Status)
	request.Remark = tcommon.StringPtr(record.Remark)

	client, err := d.NewDNSPodClient()
	if err != nil {
		return err
	}
	_, err = client.CreateRecord(request)
	return err
}

func (d *DNSPod) UpdateDNSRecord(record config.DNSRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	dr, err := d.GetDNSRecord(record.Domain, record.SubDomain)
	if err != nil {
		return err
	}

	request := dnspod.NewModifyRecordRequest()
	request.RecordId = tcommon.Uint64Ptr(dr.ID)

	request.Domain = tcommon.StringPtr(record.Domain)
	request.SubDomain = tcommon.StringPtr(record.SubDomain)
	request.RecordType = tcommon.StringPtr(record.Type)
	request.RecordLine = tcommon.StringPtr(record.Line)
	request.Value = tcommon.StringPtr(record.Value)
	request.TTL = tcommon.Uint64Ptr(record.TTL)

	// init client
	client, err := d.NewDNSPodClient()
	if err != nil {
		return err
	}
	_, err = client.ModifyRecord(request)
	return err
}

func (d *DNSPod) DeleteDNSRecord(record config.DNSRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return nil
}

func (d *DNSPod) ParseToDNSRecord(domain string, dnsPodRecord *dnspod.RecordListItem) config.DNSRecord {
	return config.DNSRecord{
		ID:         *dnsPodRecord.RecordId,
		Domain:     domain,
		SubDomain:  *dnsPodRecord.Name,
		Value:      *dnsPodRecord.Value,
		Status:     *dnsPodRecord.Status,
		Type:       *dnsPodRecord.Type,
		MX:         *dnsPodRecord.MX,
		Line:       *dnsPodRecord.Line,
		Remark:     *dnsPodRecord.Remark,
		TTL:        *dnsPodRecord.TTL,
		Weight:     *dnsPodRecord.Weight,
		UpdateTime: *dnsPodRecord.UpdatedOn,
	}
}

func (d *DNSPod) getUserUid() (uint64, error) {
	// get user uid
	request := tcam.NewGetUserAppIdRequest()

	client, err := d.NewTencentCAMClient()
	if err != nil {
		return 0, err
	}

	response, err := client.GetUserAppId(request)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(*response.Response.Uin, 10, 64)
}
