package ddns

import (
	"context"
	"errors"
	"time"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/guguducken/ddns-go/pkg/ipgetter"
	"github.com/guguducken/ddns-go/pkg/provider"
	"github.com/rs/zerolog/log"
)

func runClient(ctx context.Context, cfg *config.Config) Stopper {
	log.Warn().Msg("start ddns client")

	finishedChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.CheckInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Info().Msg("start ddns check and upgrade round")
				if err := clientRoundRun(cfg); err != nil {
					log.Error().Err(err)
				}
			case <-ctx.Done():
				// TODO: some stop work
				finishedChan <- struct{}{}
				return
			}
		}
	}()

	return &ClientStopper{
		shutdownFinished: finishedChan,
	}
}

type ClientStopper struct {
	cancel           context.CancelFunc
	shutdownFinished chan struct{}
}

func (cs *ClientStopper) Stop() chan struct{} {
	// cancel client goroutine
	go cs.cancel()
	return cs.shutdownFinished
}

func (cs *ClientStopper) SetCancelFunc(cancel context.CancelFunc) {
	cs.cancel = cancel
}

func clientRoundRun(cfg *config.Config) error {
	ip, err := ipgetter.InitIPGetters(cfg).GetIP()
	if err != nil {
		return err
	}
	log.Info().Msgf("obtain ip success, the result is: %s", ip)

	providers := provider.InitDNSProviders(cfg)

	errs := make([]error, 0, 10)
	errs = append(errs, errors.New("some dns provider create/update failed"))
	for _, p := range providers {
		if err = p.CheckPermission(); err != nil {
			return err
		}
		err = p.Do(ip)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 1 {
		return errors.Join(errs...)
	}
	return nil
}
