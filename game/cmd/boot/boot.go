package boot

import (
	"os"
	"path/filepath"

	"github.com/adm87/flinch/data"
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
			rootAbs, err := filepath.Abs(rootPath)
			if err != nil {
				return err
			}

			data.Assets.UseFilesystem(os.DirFS(filepath.Join(rootAbs, "data", "assets")))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.PersistentFlags().StringVar(&rootPath, "root-path", "", "Path to the root directory")

	return command
}
