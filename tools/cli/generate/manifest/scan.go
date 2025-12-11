package manifest

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

func Scan(model *Model, directory string, output string) error {
	directories := make(map[string]Directory)

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		parts := strings.SplitN(relPath, string(filepath.Separator), 2)
		if parts[0] == "." || parts[0] == ".." {
			return nil
		}

		if len(parts) < 2 {
			return nil
		}

		dirName := parts[0]
		fileName := d.Name()

		if fileName == output {
			return nil
		}

		dir, exists := directories[dirName]
		if !exists {
			dir = Directory{
				Name:       dirName,
				Path:       dirName,
				IsEmbedded: slices.Contains(model.Embedded, dirName),
				Files:      []File{},
			}
		}

		file := File{
			Path: relPath,
			Name: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
			Hash: fmt.Sprintf("0x%x", HashFNV(fileName)),
		}

		dir.Files = append(dir.Files, file)
		directories[dirName] = dir

		return nil
	})

	if err != nil {
		return err
	}

	model.Directories = make([]Directory, 0, len(directories))
	for _, dir := range directories {
		slices.SortFunc(dir.Files, func(f1, f2 File) int {
			return strings.Compare(f1.Path, f2.Path)
		})
		model.Directories = append(model.Directories, dir)
	}

	slices.SortFunc(model.Directories, func(d1, d2 Directory) int {
		return strings.Compare(d1.Name, d2.Name)
	})

	return nil
}
