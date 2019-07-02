package testutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

        "github.com/stretchr/testify/require"
	"opendev.org/airship/airshipctl/pkg/document"

	"sigs.k8s.io/kustomize/api/filesys"
)

// SetupTestFs help manufacture a fake file system for testing purposes. It
// will iterate over the files in fixtureDir, which is a directory relative
// to the tests themselves, and will write each of those files (preserving
// names) to an in-memory file system and return that fs
func SetupTestFs(t *testing.T, fixtureDir string) filesys.FileSystem {

	x := filesys.MakeFsInMemory()

	files, err := ioutil.ReadDir(fixtureDir)
	require.NoError(t, err, "Failed to read fixture directory, setting up testfs failed")
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(fixtureDir, fileName)
		fileBytes, err := ioutil.ReadFile(filePath)
		require.NoError(t, err, "Failed to read file, setting up testfs failed")
		err = x.WriteFile(filepath.Join("/", file.Name()), fileBytes)
		require.NoError(t, err, "Failed to write file, setting up testfs failed")
	}
	return x

}

// NewTestBundle helps to create a new bundle with FakeFs containing documents from fixtureDir
func NewTestBundle(t *testing.T, fixtureDir string) document.Bundle {
	t.Helper()
	b, err := document.NewBundle(SetupTestFs(t, fixtureDir), "/", "/")
	require.NoError(t, err, "Failed to build a bundle, setting up TestBundle failed")
	return b
}
