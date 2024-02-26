package run

import (
	"errors"
	"fmt"
	"net"

	"github.com/guguducken/ddns-go/pkg/config"
	derrors "github.com/guguducken/ddns-go/pkg/errors"
	"github.com/guguducken/ddns-go/pkg/ipcheck"
	"github.com/guguducken/ddns-go/pkg/log"
	"github.com/spf13/cobra"
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
	log.Info("config init success")

	// get ip
	log.Info("start get ip")
	ip, err := cfg.IPGetters.GetIP()
	if err != nil {
		return err
	}

	if ip == "" {
		return ipcheck.ErrAllGetterFailed
	}
	if net.ParseIP(ip) == nil {
		return errors.Join(ipcheck.ErrInvalidResponseIP, errors.New(fmt.Sprintf("response ip is %s", ip)))
	}
	log.Info(fmt.Sprintf("get ip success, ip is: %s", ip))

	return nil
}
