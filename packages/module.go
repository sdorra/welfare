package packages

// State of a package in the system
type State int

const (
	// StatePresent ensures that the package is installed
	StatePresent = iota
	// StateAbsent ensures that the package is not installed
	StateAbsent
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
	if module.State == StatePresent && !pkgInfo.Installed {
		err := module.system.Install(module.Package)
		if err != nil {
			return false, err
		}
		changed = true
	} else if module.State == StateAbsent && pkgInfo.Installed {
		err := module.system.Uninstall(module.Package)
		if err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}
