package manifest

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var output string

	model := &Model{}

	command := &cobra.Command{
		Use: "manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			workingDir, err := cmd.Flags().GetString("working-dir")
			if err != nil {
				return err
			}

			absPath, err := filepath.Abs(workingDir)
			if err != nil {
				return err
			}

			if err := Scan(model, absPath, output); err != nil {
				return err
			}

			content, err := GenerateFromTemplate(model)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filepath.Join(absPath, output), []byte(content), 0644); err != nil {
				return err
			}

			return exec.Command("go", "fmt", filepath.Join(absPath, output)).Run()
		},
	}

	command.Flags().StringVarP(&model.Package, "package", "p", model.Package, "Package name for the generated manifest.go file")
	command.Flags().StringArrayVarP(&model.Embedded, "embed", "e", model.Embedded, "Directories to embed in the manifest")
	command.Flags().StringVarP(&output, "output", "o", output, "Output path for the generated manifest.go file")

	command.MarkFlagRequired("package")
	command.MarkFlagRequired("output")

	return command
}
