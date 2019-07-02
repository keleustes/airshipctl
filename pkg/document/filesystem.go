package document

import (
	"io/ioutil"

	"sigs.k8s.io/kustomize/api/filesys"
)

// File extends kustomize File and provide abstraction to creating temporary files
type File interface {
	filesys.File
	Name() string
}

// FileSystem extends kustomize FileSystem and provide abstraction to creating temporary files
type FileSystem interface {
	filesys.FileSystem
	TempFile(string, string) (File, error)
}

// DocumentFs is adaptor to TempFile
type DocumentFs struct {
	filesys.FileSystem
}

// NewDocumentFs returns an instalce of DocumentFs
func NewDocumentFs() FileSystem {
	return &DocumentFs{FileSystem: filesys.MakeFsOnDisk()}
}

// TempFile creates file in temporary filesystem, at default os.TempDir
func (dfs DocumentFs) TempFile(tmpDir string, prefix string) (File, error) {
	return ioutil.TempFile(tmpDir, prefix)
}
