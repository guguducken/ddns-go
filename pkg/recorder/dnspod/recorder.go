package dnspod

import (
	"context"
	"errors"
	"fmt"

	dcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	dprofile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/utils"
	"github.com/guguducken/ddns-go/pkg/utils/iputil"
	"github.com/guguducken/ddns-go/pkg/utils/logutil"
)

type Recorder struct {
	config *Config

	dnsPodClient *dsdk.Client
	credential   *dcommon.Credential

	domains []domainDetail
}

func NewRecorder(ctx context.Context, config yaml.Node) (*Recorder, error) {
	r := &Recorder{
		config:  &Config{},
		domains: make([]domainDetail, 0, 10),
	}

	if err := r.config.Init(config); err != nil {
		return nil, errno.OverrideError(
			err,
			errno.AppendAdditionalMessage("phase", "init dnspod config"),
		)
	}

	if err := r.initClient(); err != nil {
		return nil, errno.OverrideError(
			err,
			errno.AppendAdditionalMessage("phase", "init dnspod client"),
		)
	}

	if err := r.initDomainDetails(ctx); err != nil {
		return nil, errno.OverrideError(
			err,
			errno.AppendAdditionalMessage("phase", "init dnspod domain details"),
		)
	}

	return r, nil
}

func (r *Recorder) initClient() (err error) {
	r.credential = dcommon.NewCredential(
		r.config.AccessKey,
		r.config.SecretKey,
	)

	cpf := dprofile.NewClientProfile()
	cpf.NetworkFailureMaxRetries = 3
	cpf.NetworkFailureRetryDuration = dprofile.ExponentialBackoff
	cpf.RateLimitExceededMaxRetries = 3
	cpf.RateLimitExceededRetryDuration = dprofile.ExponentialBackoff

	cpf.HttpProfile.Endpoint = r.config.DnsPodAPIEndpoint

	r.dnsPodClient, err = dsdk.NewClient(
		r.credential,
		"",
		cpf,
	)
	return err
}

func (r *Recorder) initDomainDetails(ctx context.Context) (err error) {
	for _, domain := range r.config.Domains {
		for _, subDomain := range domain.SubDomains {
			record, err := r.describeRecord(ctx, domain.Domain, subDomain)
			if err != nil {
				return err
			}
			detail := domainDetail{
				domain:     domain.Domain,
				subDomain:  subDomain,
				recordLine: domain.RecordLine,
				ttl:        domain.TTL,
				weight:     domain.Weight,
				status:     domain.Status,
				remark:     domain.Remark,
			}
			if record != nil {
				detail.id = *record.RecordId
				detail.value = *record.Value
				detail.recordType = cons.RecordType(*record.Type)
			}
			if !utils.CheckUniqueDomain(fmt.Sprintf("%s.%s", detail.subDomain, detail.domain)) {
				return errno.OverrideError(
					errno.ErrDomainNotUnique,
					errno.AppendAdditionalMessage("domain", detail.domain),
					errno.AppendAdditionalMessage("sub_domain", detail.subDomain),
				)
			}
			r.domains = append(r.domains, detail)
		}
	}
	return nil
}

func (r *Recorder) GetType() cons.RecorderType {
	return cons.RecorderTypeDNSPod
}

func (r *Recorder) ApplyValue(ctx context.Context, value string) (err error) {
	recordType, err := iputil.CheckIPType(value)
	if err != nil {
		return err
	}

	for ind, domain := range r.domains {
		if domain.value == value {
			logutil.Info("the record value is the same as the new ip, do nothing",
				logutil.NewField("domain", domain.domain),
				logutil.NewField("sub_domain", domain.subDomain),
				logutil.NewField("value", value),
			)
			continue
		}
		domain.value = value
		// if the record is already created, check the record type and set value to the new ip, then update the record
		if domain.id != 0 {
			if domain.recordType != recordType {
				return errno.OverrideError(
					errno.ErrRecordTypeConflict,
					errno.AppendAdditionalMessage("domain", domain.domain),
					errno.AppendAdditionalMessage("sub_domain", domain.subDomain),
					errno.AppendAdditionalMessage("record_type", string(domain.recordType)),
					errno.AppendAdditionalMessage("new_record_type", string(recordType)),
				)
			}
			if err = r.updateRecordValue(ctx, domain); err != nil {
				logutil.Error(err, "failed to update record value")
				continue
			}
		} else {
			// set the record type for the new record
			domain.recordType = recordType
			// if the record is not created, create the record
			id, err := r.createRecord(ctx, domain)
			if err != nil {
				logutil.Error(err, "failed to create record",
					logutil.NewField("domain", domain.domain),
					logutil.NewField("sub_domain", domain.subDomain),
					logutil.NewField("value", value),
				)
				continue
			}
			domain.id = id
		}
		// update the record to the domains
		r.domains[ind] = domain
	}
	return nil
}

func (r *Recorder) Exit(ctx context.Context) (err error) {
	for _, domain := range r.domains {
		if domain.id == 0 {
			logutil.Info("the record is not created, do nothing",
				logutil.NewField("domain", domain.domain),
				logutil.NewField("sub_domain", domain.subDomain),
			)
			continue
		}
		if err = r.deleteRecord(ctx, domain); err != nil {
			logutil.Error(err, "failed to delete record",
				logutil.NewField("domain", domain.domain),
				logutil.NewField("sub_domain", domain.subDomain),
			)
		}
	}
	return nil
}

func (r *Recorder) describeRecord(
	ctx context.Context,
	domain string,
	subDomain string,
) (record *dsdk.RecordListItem, err error) {
	request := dsdk.NewDescribeRecordFilterListRequest()
	request.SetContext(ctx)

	// set the domain and subdomain
	request.Domain = dcommon.StringPtr(domain)
	request.SubDomain = dcommon.StringPtr(subDomain)
	// set the limit to 1
	request.Limit = dcommon.Uint64Ptr(1)

	response, err := r.dnsPodClient.DescribeRecordFilterList(request)
	if err != nil {
		return nil, errno.OverrideError(
			errno.ErrFailedDescribeRecord,
			errno.AppendAdditionalMessage("dsdk_error", err.Error()),
		)
	}
	if *response.Response.RecordCountInfo.ListCount == 0 {
		return nil, nil
	}

	return response.Response.RecordList[0], nil
}

func (r *Recorder) updateRecordValue(
	ctx context.Context,
	domain domainDetail,
) (err error) {

	request := dsdk.NewModifyRecordRequest()
	request.SetContext(ctx)

	request.Domain = dcommon.StringPtr(domain.domain)
	request.SubDomain = dcommon.StringPtr(domain.subDomain)
	request.RecordId = dcommon.Uint64Ptr(domain.id)
	request.RecordLine = dcommon.StringPtr(domain.recordLine)
	request.RecordType = dcommon.StringPtr(string(domain.recordType))
	request.Value = dcommon.StringPtr(domain.value)
	request.TTL = dcommon.Uint64Ptr(domain.ttl)
	request.Weight = dcommon.Uint64Ptr(domain.weight)
	request.Status = dcommon.StringPtr(domain.status)
	request.Remark = dcommon.StringPtr(domain.remark)

	_, err = r.dnsPodClient.ModifyRecord(request)
	if err != nil {
		return errno.OverrideError(
			errno.ErrFailedModifyRecord,
			errno.AppendAdditionalMessage("dsdk_error", err.Error()),
		)
	}
	return nil
}

func (r *Recorder) createRecord(
	ctx context.Context,
	domain domainDetail,
) (uint64, error) {
	request := dsdk.NewCreateRecordRequest()
	request.SetContext(ctx)

	request.Domain = dcommon.StringPtr(domain.domain)
	request.SubDomain = dcommon.StringPtr(domain.subDomain)
	request.RecordType = dcommon.StringPtr(string(domain.recordType))
	request.Value = dcommon.StringPtr(domain.value)
	request.RecordLine = dcommon.StringPtr(domain.recordLine)
	request.TTL = dcommon.Uint64Ptr(domain.ttl)
	request.Weight = dcommon.Uint64Ptr(domain.weight)
	request.Status = dcommon.StringPtr(domain.status)
	request.Remark = dcommon.StringPtr(domain.remark)

	response, err := r.dnsPodClient.CreateRecord(request)
	if err != nil {
		return 0, errno.OverrideError(
			errno.ErrFailedCreateRecord,
			errno.AppendAdditionalMessage("dsdk_error", err.Error()),
		)
	}
	return *response.Response.RecordId, nil
}

func (r *Recorder) deleteRecord(
	ctx context.Context,
	domain domainDetail,
) (err error) {
	if domain.id == 0 {
		return errors.New("record id is 0, please create the record first")
	}
	request := dsdk.NewDeleteRecordRequest()
	request.SetContext(ctx)

	request.Domain = dcommon.StringPtr(domain.domain)
	request.RecordId = dcommon.Uint64Ptr(domain.id)

	_, err = r.dnsPodClient.DeleteRecord(request)
	if err != nil {
		return errno.OverrideError(
			errno.ErrFailedDeleteRecord,
			errno.AppendAdditionalMessage("dsdk_error", err.Error()),
		)
	}
	return nil
}
