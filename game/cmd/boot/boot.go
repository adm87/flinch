package boot

import (
	"os"
	"path/filepath"

	"github.com/adm87/flinch/data"
	"github.com/adm87/flinch/game/src/game"
	"github.com/hajimehoshi/ebiten/v2"
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
			absRoot, err := filepath.Abs(rootPath)
			if err != nil {
				return err
			}

			// Link the asset resource system to its filesystem on disk.
			data.Assets.SetFileSystem(os.DirFS(filepath.Join(absRoot, "data", "assets")))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ebiten.RunGame(game.NewGame())
		},
	}

	command.PersistentFlags().StringVar(&rootPath, "root-path", "", "Path to the root directory")

	return command
}
