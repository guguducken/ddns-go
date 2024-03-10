package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/guguducken/ddns-go/pkg/config"
)

func TestDNSPod_CheckPermission(t *testing.T) {
	d := NewDNSPodProvider(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"), []config.DNSRecord{})
	fmt.Printf("d.CheckPermission(): %v\n", d.CheckPermission())
}

func TestDNSPod_CreateDNSRecord(t *testing.T) {
	d := NewDNSPodProvider(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"), config.DNSRecords{
		config.DNSRecord{
			Domain: "1matrix.org",
			// SubDomain: "test",
			Value:  "124.77.159.88",
			Type:   IPv4,
			Line:   "默认",
			TTL:    600,
			Remark: "generated by ddns-go",
		},
	})
	if err := d.CheckPermission(); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		t.FailNow()
	}
	err := d.CreateDNSRecord(d.domains[0])
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		t.FailNow()
	}
}

func TestDNSPod_UpdateDNSRecord(t *testing.T) {
	d := NewDNSPodProvider(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"), config.DNSRecords{
		config.DNSRecord{
			Domain:    "1matrix.org",
			SubDomain: "test",
			Value:     "124.77.159.88",
			Type:      IPv4,
			Line:      "默认",
			TTL:       600,
			Remark:    "generated by ddns-go",
		},
	})
	if err := d.CheckPermission(); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		t.FailNow()
	}
	err := d.UpdateDNSRecord(d.domains[0])
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		t.FailNow()
	}
}
