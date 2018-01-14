package files_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/sdorra/welfare/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileModule_RunWithStateFileAndContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")

	file := files.NewFileModule(target, files.File)
	file.Content = "Hello My Name is"

	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	bytes, err := ioutil.ReadFile(target)
	assert.Nil(t, err)
	assert.Equal(t, "Hello My Name is", string(bytes))
}

func TestFileModule_RunWithStateAbsentExistingFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hello My Name is"), 0644)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.Absent)
	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	_, err = os.Stat(target)
	assert.True(t, os.IsNotExist(err))
}

func TestFileModule_RunWithStateAbsentExistingDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = os.Mkdir(target, 0777)
	require.Nil(t, err)

	one := path.Join(target, "one")
	err = ioutil.WriteFile(one, []byte("Hello My Name is"), 0644)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.Absent)
	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	_, err = os.Stat(target)
	assert.True(t, os.IsNotExist(err))
}

func TestFileModule_RunWithStateAbsentNonExistingFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")

	file := files.NewFileModule(target, files.Absent)
	changed, err := file.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestFileModule_RunWithStateFileAndContentExistingFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hello My Name is"), 0644)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.File)
	file.Content = "Hello My Name is"

	changed, err := file.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestFileModule_RunWithStateFileAndOtherContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hi My Name is."), 0644)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.File)
	file.Content = "Hello My Name is"

	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	bytes, err := ioutil.ReadFile(target)
	assert.Nil(t, err)
	assert.Equal(t, "Hello My Name is", string(bytes))
}

func TestFileContentModule_RunWithStateFileAndWrongPermissions(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hello My Name is"), 0777)
	require.Nil(t, err)

	content := files.NewFileModule(target, files.File)
	content.Content = "Hello My Name is"

	changed, err := content.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)
	assert.Equal(t, content.FileMode, stat.Mode())
}

func TestFileModule_RunWithStateDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = os.Mkdir(target, 0755)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.Directory)
	changed, err := file.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestFileModule_RunWithStateDirectoryNonExisting(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")

	file := files.NewFileModule(target, files.Directory)
	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())
}

func TestFileModule_RunWithStateDirectoryChangePermissions(t *testing.T) {
	dir, err := ioutil.TempDir("", "file")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")
	err = os.Mkdir(target, 0777)
	require.Nil(t, err)

	err = os.Chmod(target, 0777)
	require.Nil(t, err)

	file := files.NewFileModule(target, files.Directory)
	file.FileMode = 0755

	changed, err := file.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)

	assert.Equal(t, os.FileMode(0755), stat.Mode().Perm())
}
