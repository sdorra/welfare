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

func TestTemplateModule_Run(t *testing.T) {
	dir, err := ioutil.TempDir("", "template")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	target := path.Join(dir, "target")

	tpl := files.NewTemplateModule(target, "Hello My Name is {{.Name}}", &Context{"sorbot"})

	changed, err := tpl.Run()
	assert.Nil(t, err)
	assert.True(t, changed)

	bytes, err := ioutil.ReadFile(target)
	assert.Nil(t, err)
	assert.Equal(t, "Hello My Name is sorbot", string(bytes))
}

type Context struct {
	Name string
}
