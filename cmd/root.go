package cmd

import (
	"fmt"
	"os"

	"github.com/guguducken/ddns-go/cmd/run"
	"github.com/guguducken/ddns-go/pkg/log"
	"github.com/spf13/cobra"
)

func initRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "ddns-go",
		Short:            "ddns client which written by golang",
		PersistentPreRun: preInitCommand,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// add flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "specify the path of config file")
	rootCmd.PersistentFlags().StringP("log-level", "l", "", "specify log level, it will override the all other settings")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode, it will override the settings in the configuration file")

	return rootCmd
}

func preInitCommand(cmd *cobra.Command, args []string) {
	if cmd.Flag("log-level").Changed {
		loglevel := cmd.Flag("log-level").Value.String()
		log.Init(loglevel)
	}
	// check whether debug mode is enabled, It will override the settings in the configuration file
	if cmd.Flag("debug").Changed {
		log.Init(log.DebugLevel)
	}
}

func Execute() {
	rootCommand := initRootCommand()
	// add sub command to root command
	rootCommand.AddCommand(run.InitRunCommand())

	// run
	err := rootCommand.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
