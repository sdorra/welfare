package packages_test

import (
	"io/ioutil"
	"os"
	"testing"

	"path"

	"github.com/sdorra/welfare/packages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAptRepositoryModule_Run(t *testing.T) {
	dir, err := ioutil.TempDir("", "apt_repository")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	repo := packages.NewAptRepositoryModule("scm-manager", packages.Present)
	repo.Directory = dir
	repo.Repository = "deb http://maven.scm-manager.org/nexus/content/repositories/releases ./"

	changed, err := repo.Run()

	assert.Nil(t, err)
	assert.True(t, changed)
}

func TestAptRepositoryModule_RunWithPresentRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "apt_repository")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	sourcesDir := path.Join(dir, "sources.list.d")
	err = os.MkdirAll(sourcesDir, 0755)
	require.Nil(t, err)

	content := "deb http://maven.scm-manager.org/nexus/content/repositories/releases ./\n"
	repoFile := path.Join(sourcesDir, "scm.list")
	err = ioutil.WriteFile(repoFile, []byte(content), 0644)
	require.Nil(t, err)

	repo := packages.NewAptRepositoryModule("scm-manager", packages.Present)
	repo.Directory = dir
	repo.Repository = "deb http://maven.scm-manager.org/nexus/content/repositories/releases ./"

	changed, err := repo.Run()

	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestAptRepositoryModule_RunWithPresentRepositoryInSourceList(t *testing.T) {
	dir, err := ioutil.TempDir("", "apt_repository")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	content := ` 
		# See http://help.ubuntu.com/community/UpgradeNotes for how to upgrade to
		# newer versions of the distribution.
		deb http://archive.ubuntu.com/ubuntu/ xenial main restricted
		# deb-src http://archive.ubuntu.com/ubuntu/ xenial main restricted
		## Major bug fix updates produced after the final release of the
		## distribution.
		
		deb http://archive.ubuntu.com/ubuntu/ xenial-updates main restricted
		# deb-src http://archive.ubuntu.com/ubuntu/ xenial-updates main restricted
		
		## N.B. software from this repository is ENTIRELY UNSUPPORTED by the Ubuntu
		## team. Also, please note that software in universe WILL NOT receive any
		## review or updates from the Ubuntu security team.
		deb http://archive.ubuntu.com/ubuntu/ xenial universe
		deb-src http://archive.ubuntu.com/ubuntu/ xenial universe
		deb http://archive.ubuntu.com/ubuntu/ xenial-updates universe
		deb-src http://archive.ubuntu.com/ubuntu/ xenial-updates universe`

	sourcesFile := path.Join(dir, "sources.list")
	err = ioutil.WriteFile(sourcesFile, []byte(content), 0644)
	require.Nil(t, err)

	repo := packages.NewAptRepositoryModule("archive", packages.Present)
	repo.Directory = dir
	repo.Repository = "deb http://archive.ubuntu.com/ubuntu/ xenial universe"

	changed, err := repo.Run()

	assert.Nil(t, err)
	assert.False(t, changed)
}
