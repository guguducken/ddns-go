package ddns

import (
	"context"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func Run(cfg *config.Config) (stopper Stopper) {
	ctx, cancel := context.WithCancel(context.Background())
	switch cfg.Type {
	case "client":
		stopper = runClient(ctx, cfg)
	case "server":
	default:
		log.Fatal().Msg("unsupported run type")
	}
	stopper.SetCancelFunc(cancel)
	return stopper
}
