package plain

import (
	"testing"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"gopkg.in/yaml.v3"
)

func TestProvider(t *testing.T) {
	logutil.Init("debug", nil)
	config := `config:
  request_url: ""
`
	type Temp struct {
		Config yaml.Node `yaml:"config"`
	}
	var temp Temp
	err := yaml.Unmarshal([]byte(config), &temp)
	if err != nil {
		t.Fatal(err)
	}
	provider, err := NewProvider(temp.Config, true)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}
	ip, err := provider.ProviderIP()
	if err != nil {
		logutil.Fatal(err, "failed to get ip")
	}
	logutil.Info("ip", logutil.NewField("ip", ip))
}
