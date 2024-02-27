package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guguducken/ddns-go/pkg/config"
	"github.com/guguducken/ddns-go/pkg/ddns"
	derrors "github.com/guguducken/ddns-go/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	mostShutdownTime = 10 * time.Second
)

func InitRunCommand() *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "start running ddns client or ip getter server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get config path
			configPath := cmd.Flag("config").Value.String()
			if configPath == "" {
				return derrors.ErrEmptyConfigPath
			}
			return run(configPath)
		},
	}
	return runCommand
}

func run(configPath string) error {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}
	log.Info().Msg("config init success")

	// run ddns client or ip getter server
	stopper := ddns.Run(cfg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	sig := <-signalChan
	log.Warn().Msg("start shutdown")
	log.Warn().Msg(fmt.Sprintf("shutdown signal is %s", sig.String()))

	shutdownFinished := stopper.Stop()

	log.Warn().Msg(fmt.Sprintf("most shutdown time is %s", mostShutdownTime.String()))
	timer := time.NewTimer(mostShutdownTime)
	defer timer.Stop()

	for {
		select {
		case <-shutdownFinished:
			log.Warn().Msg("shutdown finished")
			return nil
		case <-timer.C:
			log.Error().Msg("maximum shutdown time exceeded")
			return nil
		case <-signalChan:
			log.Error().Msg(fmt.Sprintf("receive signal again, will exit now"))
			return nil
		}
	}
}
