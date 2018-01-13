package packages

// State of a package in the system
type State int

const (
	// Present ensures that the package is installed
	Present = iota
	// Absent ensures that the package is not installed
	Absent
)

// PackageModule ensures the state of a package
type PackageModule struct {
	Package string
	State   State
	system  packageSystem
}

func (module *PackageModule) Run() (bool, error) {
	pkgInfo := module.system.GetInfo(module.Package)

	changed := false
	if module.State == Present && !pkgInfo.Installed {
		err := module.system.Install(module.Package)
		if err != nil {
			return false, err
		}
		changed = true
	} else if module.State == Absent && pkgInfo.Installed {
		err := module.system.Uninstall(module.Package)
		if err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}
