package modules_test

import (
	"io/ioutil"
	"path"
	"testing"

	"os"

	"github.com/sdorra/welfare/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentModule_Run(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)

	target := path.Join(dir, "target")

	content := modules.NewContentModule(target, "Hello My Name is")

	changed, err := content.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	bytes, err := ioutil.ReadFile(target)
	assert.Nil(t, err)
	assert.Equal(t, "Hello My Name is", string(bytes))
}

func TestContentModule_RunWithEqualContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hello My Name is"), 0644)
	require.Nil(t, err)

	content := modules.NewContentModule(target, "Hello My Name is")

	changed, err := content.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestContentModule_RunWithOtherContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hi My Name is."), 0644)
	require.Nil(t, err)

	content := modules.NewContentModule(target, "Hello My Name is")

	changed, err := content.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	bytes, err := ioutil.ReadFile(target)
	assert.Nil(t, err)
	assert.Equal(t, "Hello My Name is", string(bytes))
}

func TestContentModule_RunWithWrongPermissions(t *testing.T) {
	dir, err := ioutil.TempDir("", "copy")
	require.Nil(t, err)

	target := path.Join(dir, "target")
	err = ioutil.WriteFile(target, []byte("Hello My Name is"), 0777)
	require.Nil(t, err)

	content := modules.NewContentModule(target, "Hello My Name is")

	changed, err := content.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	stat, err := os.Stat(target)
	assert.Nil(t, err)
	assert.Equal(t, content.FileMode, stat.Mode())
}
