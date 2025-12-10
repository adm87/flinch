package manifest

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	model := &Model{}

	command := &cobra.Command{
		Use: "manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.Flags().StringVarP(&model.Package, "package", "p", "", "Package name for the generated manifest.go file")
	command.Flags().StringArrayVarP(&model.Embedded, "embed", "e", []string{}, "Directories to embed in the manifest")

	command.MarkFlagRequired("package")

	return command
}
