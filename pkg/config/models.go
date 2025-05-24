package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/utils"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err = yaml.Unmarshal(utils.UnsafeToByteSlice(utils.ParseEnv(string(data))), c); err != nil {
		return nil, err
	}
	if c.CheckInterval == 0 {
		c.CheckInterval = 10
	}
	return c, nil
}

type Config struct {
	V4 *Pairs `yaml:"v4"`
	V6 *Pairs `yaml:"v6"`

	CheckInterval int `yaml:"check_interval"`
}

type Pairs struct {
	Providers []Provider `yaml:"providers"`
	Recorders []Recorder `yaml:"recorders"`
}

type Provider struct {
	Type   cons.ProviderType `yaml:"type"`
	Name   string            `yaml:"name"`
	Config yaml.Node         `yaml:"config"`
}

type Recorder struct {
	Type   cons.RecorderType `yaml:"type"`
	Config yaml.Node         `yaml:"config"`
}
