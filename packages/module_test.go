package packages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageModule_Run(t *testing.T) {
	system := &testPackageSystem{
		info: packageInfo{
			Installed: false,
		},
	}

	module := PackageModule{
		Package: "htop",
		State:   Present,
		system:  system,
	}

	changed, err := module.Run()
	assert.Nil(t, err)
	assert.True(t, changed)
	assert.Equal(t, "install", system.action)
	assert.Equal(t, "htop", system.pkg)
}

func TestPackageModule_RunAlreadyInstalled(t *testing.T) {
	system := &testPackageSystem{
		info: packageInfo{
			Installed: true,
		},
	}

	module := PackageModule{
		Package: "htop",
		State:   Present,
		system:  system,
	}

	changed, err := module.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

func TestPackageModule_RunAlreadyUninstall(t *testing.T) {
	system := &testPackageSystem{
		info: packageInfo{
			Installed: true,
		},
	}

	module := PackageModule{
		Package: "htop",
		State:   Absent,
		system:  system,
	}

	changed, err := module.Run()
	assert.Nil(t, err)
	assert.True(t, changed)
	assert.Equal(t, "uninstall", system.action)
	assert.Equal(t, "htop", system.pkg)
}

func TestPackageModule_RunAlreadyUninstalled(t *testing.T) {
	system := &testPackageSystem{
		info: packageInfo{
			Installed: false,
		},
	}

	module := PackageModule{
		Package: "htop",
		State:   Absent,
		system:  system,
	}

	changed, err := module.Run()
	assert.Nil(t, err)
	assert.False(t, changed)
}

type testPackageSystem struct {
	info   packageInfo
	action string
	pkg    string
}

func (system *testPackageSystem) GetInfo(pkg string) packageInfo {
	return system.info
}

func (system *testPackageSystem) Install(pkg string) error {
	system.action = "install"
	system.pkg = pkg
	return nil
}

func (system *testPackageSystem) Uninstall(pkg string) error {
	system.action = "uninstall"
	system.pkg = pkg
	return nil
}
