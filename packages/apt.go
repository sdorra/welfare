package packages

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// NewAptModule creates a new PackageModule for debian based operating systems
func NewAptModule(pkg string, state State) *PackageModule {
	apt := &aptPackageSystem{}
	return &PackageModule{
		Package: pkg,
		State:   state,
		system:  apt,
	}
}

type aptPackageSystem struct {
}

func (apt *aptPackageSystem) GetInfo(pkg string) packageInfo {
	packageInfo := packageInfo{}

	_, err := exec.Command("dpkg", "-s", pkg).Output()
	if err != nil {
		packageInfo.Installed = false
	} else {
		// TODO read version from output
		packageInfo.Installed = true
	}

	return packageInfo
}

func (apt *aptPackageSystem) Install(pkg string) error {
	env := append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")

	cmd := exec.Command("apt-get", "-y", "update")
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute package update command")
	}

	cmd = exec.Command("apt-get", "-y", "install", pkg)
	cmd.Env = env
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute package update command")
	}

	return nil
}

func (apt *aptPackageSystem) Uninstall(pkg string) error {
	cmd := exec.Command("apt-get", "-y", "remove", pkg)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute package update command")
	}

	return nil
}
