package generate

import (
	"github.com/adm87/flinch/tools/cli/generate/manifest"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.AddCommand(
		manifest.Command(),
	)

	return command
}
