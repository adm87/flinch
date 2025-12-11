package manifest

type Model struct {
	Package     string
	Embedded    []string
	Directories []Directory
}

type Directory struct {
	Name       string
	Path       string
	IsEmbedded bool
	Files      []File
}

type File struct {
	Path string
	Name string
	Hash string
}
