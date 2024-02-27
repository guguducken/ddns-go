package ddns

import (
	"context"
	"time"

	"github.com/guguducken/ddns-go/pkg/config"
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
			case t := <-ticker.C:
				log.Info().Msg("start ddns check and upgrade round")
				if err := clientRoundRun(cfg, t); err != nil {
					log.Error().Err(err)
				}
			case <-ctx.Done():
				// TODO: some stop work
				for i := 0; i < 20; i++ {
					log.Info().Int("times", i).Msg("simulate stop some work")
					time.Sleep(1 * time.Second)
				}
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

func clientRoundRun(cfg *config.Config, t time.Time) error {
	ip, err := cfg.IPGetters.GetIP()
	if err != nil {
		return err
	}
	log.Info().Msgf("obtain ip success, the result is: %s", ip)

	errs := cfg.DNSAppliers.Apply(ip)
	// log errors if len(errs) != 0
	if len(errs) != 0 {
		log.Error().Errs("errors", errs).Msg("some applier report upgrade dns record to provider failed")
	}
	if len(errs) == cfg.GetTotalDomains() {
		return config.ErrAllApplierFailed
	}
	return nil
}
