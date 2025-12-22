package boot

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		rootPath string
	)

	command := &cobra.Command{
		Use:   "flinch",
		Short: "Launch Flinch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.PersistentFlags().StringVar(&rootPath, "root-path", "", "Path to the root directory")

	return command
}
