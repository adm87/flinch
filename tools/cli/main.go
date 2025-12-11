package main

import (
	"os"

	"github.com/adm87/flinch/tools/cli/generate"
	"github.com/spf13/cobra"
)

func main() {
	var (
		workingDir string
	)

	command := &cobra.Command{
		Use: "flinch-cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.PersistentFlags().StringVarP(&workingDir, "working-dir", "C", "", "Set the working directory for the command")

	command.AddCommand(
		generate.Command(),
	)

	if err := command.Execute(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
