package cmd

import (
	"fmt"
	"os"

	"github.com/guguducken/ddns-go/cmd/run"
	"github.com/spf13/cobra"
)

func initRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ddns-go",
		Short: "ddns client which written by golang",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// add flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "specify the path of config file")
	return rootCmd
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
