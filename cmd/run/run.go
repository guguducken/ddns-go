package run

import (
	derrors "github.com/guguducken/ddns-go/pkg/errors"
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
	return nil
}
