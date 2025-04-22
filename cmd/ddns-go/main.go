package main

import (
	"flag"

	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"github.com/guguducken/ddns-go/pkg/version"
)

type Options struct {
	LogLevel   *string
	ConfigFile *string
	Version    *bool
}

func main() {
	opts := NewOptions()
	if *opts.Version {
		version.Print()
		return
	}
	logutil.Init(*opts.LogLevel, nil)
}

func NewOptions() *Options {
	opts := &Options{
		LogLevel:   flag.String("log-level", "info", "log level, valid values: [debug|info|warn|error|panic|fatal]"),
		ConfigFile: flag.String("config", "config.yaml", "config file path"),
		Version:    flag.Bool("version", false, "show version information"),
	}
	flag.Usage = Usage()
	flag.Parse()
	return opts
}
