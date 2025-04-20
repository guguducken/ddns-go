package config

import (
	"github.com/guguducken/ddns-go/pkg/cons"
	"gopkg.in/yaml.v3"
)

type Config struct {
	V4 Providers `yaml:"v4"`
	V6 Providers `yaml:"v6"`
}

type Providers struct {
	Providers []Provider `yaml:"providers"`
	Recorders []Recorder `yaml:"recorders"`
}

type Provider struct {
	Type   cons.ProviderType `yaml:"type"`
	Config yaml.Node         `yaml:"config"`
}

type Recorder struct {
	Type   cons.RecorderType `yaml:"type"`
	Config yaml.Node         `yaml:"config"`
}
