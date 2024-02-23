package provider

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	tcam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

const (
	dnsPodAPIEndpoint     = "dnspod.tencentcloudapi.com"
	tencentCAMAPIEndpoint = "cam.tencentcloudapi.com"
)

type DNSPod struct {
	accessKey  string
	secretKey  string
	credential *tcommon.Credential
}

func NewDNSPodProvider(accessKey string, secretKey string) DNSPod {
	return DNSPod{
		accessKey:  accessKey,
		secretKey:  secretKey,
		credential: tcommon.NewCredential(accessKey, secretKey),
	}
}

func (d DNSPod) NewDNSPodClient() (*dnspod.Client, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = dnsPodAPIEndpoint
	return dnspod.NewClient(d.credential, "", cpf)
}

func (d DNSPod) NewTencentCAMClient() (*tcam.Client, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = dnsPodAPIEndpoint
	return tcam.NewClient(d.credential, "", cpf)
}

func (d DNSPod) GetType() string {
	return DNSPodProvider
}

func (d DNSPod) CheckPermission() error {
	userUID, err := d.getUserUid()
	if err != nil {
		return err
	}
	allowedPermissions := []string{"AdministratorAccess", "QCloudResourceFullAccess", "QcloudDNSPodFullAccess"}

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
			if slices.Contains(allowedPermissions, *policy.PolicyType) {
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

func (d DNSPod) GetDNSRecord(domain string, subDomain string) (DNSRecord, error) {
	var record DNSRecord

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
	record = d.ParseToDNSRecord(response.Response.RecordList[0])
	return record, nil
}

func (d DNSPod) ListDNSRecords(domain string) (DNSRecords, error) {
	records := make(DNSRecords, 0, 30)

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = tcommon.StringPtr(domain)

	client, err := d.NewDNSPodClient()
	if err != nil {
		return nil, err
	}

	response, err := client.DescribeRecordList(request)
	if err != nil {
		return nil, err
	}
	for _, item := range response.Response.RecordList {
		records = append(records, d.ParseToDNSRecord(item))
	}

	return records, nil
}

func (d DNSPod) CreateDNSRecord(record DNSRecord) error {
	return nil
}

func (d DNSPod) UpdateDNSRecord(record DNSRecord) error {
	return nil
}

func (d DNSPod) DeleteDNSRecord(record DNSRecord) error {
	return nil
}

func (d DNSPod) ParseToDNSRecord(dnsPodRecord *dnspod.RecordListItem) DNSRecord {
	return DNSRecord{
		Name:       *dnsPodRecord.Name,
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

func (d DNSPod) getUserUid() (uint64, error) {
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
