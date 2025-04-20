package test

import (
	"os"
	"testing"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/guguducken/ddns-go/pkg/provider"
	"github.com/guguducken/ddns-go/pkg/recorder"
	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"gopkg.in/yaml.v3"
)

func TestProvider(t *testing.T) {
	logutil.Init("debug", nil)
	data, err := os.ReadFile("provider.yaml")
	if err != nil {
		logutil.Fatal(err, "failed to read provider.yaml")
	}

	var configs config.Config
	if err = yaml.Unmarshal(data, &configs); err != nil {
		t.Fatal(err)
	}

	type Pairs struct {
		Providers provider.Providers
		Recorders []recorder.Recorder
	}

	v4s := Pairs{
		Providers: make(provider.Providers, 0, len(configs.V4.Providers)),
		Recorders: make([]recorder.Recorder, 0, len(configs.V4.Recorders)),
	}
	// v6s := Pairs{
	// 	Providers: make([]provider.Provider, 0, len(configs.V6.Providers)),
	// 	Recorders: make([]recorder.Recorder, 0, len(configs.V6.Recorders)),
	// }
	for _, cfg := range configs.V4.Providers {
		p, err := provider.NewProvider(cfg.Type, cfg.Config, true)
		if err != nil {
			logutil.Fatal(err, "failed to create provider")
		}
		v4s.Providers = append(v4s.Providers, p)
	}
	// for _, cfg := range configs.V6.Recorders {
	// 	r, err := recorder.NewRecorder(context.Background(), cfg.Type, cfg.Config)
	// 	if err != nil {
	// 		logutil.Fatal(err, "failed to create recorder")
	// 	}
	// 	v4s.Recorders = append(v4s.Recorders, r)
	// }
	ip, err := v4s.Providers[1].ProviderIP()
	if err != nil {
		logutil.Fatal(err, "failed to get provider ip")
	}
	logutil.Info("ip", logutil.NewField("ip", ip))
	// err = v4s.Recorders[0].ApplyValue(context.Background(), ip)
	// if err != nil {
	// 	logutil.Fatal(err, "failed to apply value")
	// }
	// err = v4s.Recorders[0].Exit(context.Background())
	// if err != nil {
	// 	logutil.Fatal(err, "failed to exit")
	// }
}
