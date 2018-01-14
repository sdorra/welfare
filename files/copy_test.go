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

func TestCopyModule_Run(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")
	err = ioutil.WriteFile(source, []byte("a"), 0644)
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("b"), 0644)
	require.Nil(t, err)

	copy := files.NewCopyModule(source, target)

	changed, err := copy.Run()
	assert.Nil(t, err)
	assert.True(t, changed)
}

func TestCopyModule_RunWithEqualContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")
	err = ioutil.WriteFile(source, []byte("a"), 0644)
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("a"), 0644)
	require.Nil(t, err)

	copy := files.NewCopyModule(source, target)

	changed, err := copy.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestCopyModule_RunWithChangedPermissionsFromSource(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")
	err = ioutil.WriteFile(source, []byte("a"), 0644)
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("a"), 0755)
	require.Nil(t, err)

	copy := files.NewCopyModule(source, target)

	changed, err := copy.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)
	assert.Equal(t, os.FileMode(0644), stat.Mode())
}

func TestCopyModule_RunWithChangedPermissions(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")
	err = ioutil.WriteFile(source, []byte("a"), 0755)
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("a"), 0755)
	require.Nil(t, err)

	copy := files.NewCopyModule(source, target)
	copy.FileMode = 0644

	changed, err := copy.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)
	assert.Equal(t, os.FileMode(0644), stat.Mode())
}

func TestCopyModule_RunWithoutTarget(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")
	err = ioutil.WriteFile(source, []byte("a"), 0644)
	require.Nil(t, err)

	target := path.Join(dir, "target")

	copy := files.NewCopyModule(source, target)

	changed, err := copy.Run()
	assert.Nil(t, err)
	assert.True(t, changed)
}

func TestCopyModule_RunWithoutSource(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	source := path.Join(dir, "source")

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("a"), 0644)
	require.Nil(t, err)

	copy := files.NewCopyModule(source, target)

	_, err = copy.Run()
	assert.Error(t, err)
}
