package provider

import (
	"fmt"
	"os"
	"testing"
)

func TestDNSPod_CheckPermission(t *testing.T) {
	d := NewDNSPodProvider(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"))
	fmt.Printf("d.CheckPermission(): %v\n", d.CheckPermission())
}

func TestDNSPod_CreateDNSRecord(t *testing.T) {
	d := NewDNSPodProvider(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"))
}
