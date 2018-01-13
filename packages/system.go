package packages

type packageSystem interface {
	GetInfo(pkg string) packageInfo
	Install(pkg string) error
	Uninstall(pkg string) error
}

type packageInfo struct {
	Installed bool
}
