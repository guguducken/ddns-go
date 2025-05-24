package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/provider"
	"github.com/guguducken/ddns-go/pkg/recorder"
	"github.com/guguducken/ddns-go/pkg/utils/logutil"
	"github.com/guguducken/ddns-go/pkg/version"
)

type Options struct {
	LogLevel    *string
	ConfigFile  *string
	Version     *bool
	ExitTimeOut *int
}

func main() {
	opts := NewOptions()
	if *opts.Version {
		version.Print()
		return
	}
	// load log config and init log util
	logutil.Init(*opts.LogLevel, nil)
	// load config file to config.Config
	cfg, err := config.LoadConfig(*opts.ConfigFile)
	if err != nil {
		logutil.Fatal(err, "failed to load config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	finishedChan := make(chan struct{}, 1)
	ExitTimeOut := time.Duration(*opts.ExitTimeOut) * time.Second
	start(ctx, cfg, finishedChan, ExitTimeOut)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
	sig := <-signalChan
	logutil.Info(
		fmt.Sprintf("receive system signal, ddns-go will exit after clean with timeout %d seconds", *opts.ExitTimeOut),
		logutil.NewField("signal", sig.String()),
	)

	// cancel all process
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ExitTimeOut)
	defer shutdownCancel()
	for {
		select {
		case <-shutdownCtx.Done():
			logutil.Error(nil, "shutdown timed out, will exit now")
			os.Exit(1)
		case <-finishedChan:
			logutil.Info("shutdown finished, will exit now")
			os.Exit(0)
		case s := <-signalChan:
			logutil.Info(
				"receive shutdown signal again, will exit now",
				logutil.NewField("signal", s.String()),
			)
			os.Exit(1)
		}
	}

}

func NewOptions() *Options {
	opts := &Options{
		LogLevel:    flag.String("log-level", "info", "log level, valid values: [debug|info|warn|error|panic|fatal]"),
		ConfigFile:  flag.String("config", "config.yaml", "config file path"),
		Version:     flag.Bool("version", false, "show version information"),
		ExitTimeOut: flag.Int("exit-timeout", 30, "exit timeout in seconds"),
	}
	flag.Usage = Usage()
	flag.Parse()
	return opts
}

func start(ctx context.Context, cfg *config.Config, finishedChan chan struct{}, exitTimeout time.Duration) {
	eg := errgroup.Group{}
	if cfg.V4 != nil {
		eg.Go(func() error {
			return startCycle(ctx, cfg.V4, cfg.CheckInterval, true, exitTimeout)
		})
	}
	if cfg.V6 != nil {
		eg.Go(func() error {
			return startCycle(ctx, cfg.V6, cfg.CheckInterval, true, exitTimeout)
		})
	}
	go func() {
		if err := eg.Wait(); err != nil {
			logutil.Fatal(err, "some error occurred")
		}
		finishedChan <- struct{}{}
	}()
}

func startCycle(
	ctx context.Context,
	cfg *config.Pairs,
	checkInterval int,
	isV4 bool,
	exitTimeout time.Duration,
) error {
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	if len(cfg.Providers) == 0 {
		return errno.ErrNoConfiguredProvider
	}
	if len(cfg.Recorders) == 0 {
		logutil.Error(errno.ErrNoConfiguredRecorder, "no configured recorders, will return nil and not apply record")
	}
	// init providers and recorders
	providers, recorders := make(provider.Providers, 0, len(cfg.Providers)), make(recorder.Recorders, 0, len(cfg.Recorders))
	// init providers first
	for _, p := range cfg.Providers {
		pp, err := provider.NewProvider(p.Type, p.Config, p.Name, isV4)
		if err != nil {
			return err
		}
		providers = append(providers, pp)
	}

	// init recorders
	for _, r := range cfg.Recorders {
		rr, err := recorder.NewRecorder(ctx, r.Type, r.Config)
		if err != nil {
			return err
		}
		recorders = append(recorders, rr)
	}
	for {
		select {
		case <-ticker.C:
			ip, err := providers.ProviderIP()
			if err != nil {
				logutil.Error(err, "failed to provide ip, will skip apply record")
				continue
			}
			if err = recorders.ApplyValue(ctx, ip); err != nil {
				logutil.Error(err, "failed to apply record")
			}
		case <-ctx.Done():
			exitCtx, _ := context.WithTimeout(context.Background(), exitTimeout)
			return recorders.Exit(exitCtx)
		}
	}
}
