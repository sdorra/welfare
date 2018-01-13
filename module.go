package welfare

// Module represents a single declarative module
type Module interface {
	Run() (changed bool, err error)
}
