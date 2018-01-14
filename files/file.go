package files

import (
	"os"

	"github.com/pkg/errors"
)

// NewFileModule creates new FileModule for the given path
func NewFileModule(path string, state State) *FileModule {
	module := &FileModule{
		Path:  path,
		State: state,
	}

	mode := os.FileMode(0644)
	if state == Directory {
		mode = os.FileMode(0755)
	}

	module.FileMode = mode
	module.UID = os.Getuid()
	module.GID = os.Getegid()
	return module
}

// FileModule ensures the state of a path
type FileModule struct {
	permissions
	Path    string
	Content string
	State   State
}

func (module *FileModule) Run() (bool, error) {
	target, err := collectFileInfo(module.Path)
	if err != nil {
		return false, err
	}

	switch module.State {
	case File:
		return module.file(target)
	case Directory:
		return module.directory(target)
	case Absent:
		return module.absent(target)
	default:
		return false, errors.New("not yet implemented")
	}
}

func (module *FileModule) file(target fileInfo) (bool, error) {
	contentChanged, err := ensureContent(target, module.Content, module.FileMode)
	if err != nil {
		return false, err
	}

	permissionsChanged, err := ensurePermissions(module.permissions, target)
	if err != nil {
		return false, err
	}

	return contentChanged || permissionsChanged, nil
}

func (module *FileModule) directory(target fileInfo) (bool, error) {
	directoryChanged := false

	switch target.State {
	case Absent:
		err := os.MkdirAll(target.Path, module.FileMode)
		if err != nil {
			return false, errors.Wrapf(err, "failed to create directory %s", target.Path)
		}

		directoryChanged = true
	}

	permissionsChanged, err := ensurePermissions(module.permissions, target)
	if err != nil {
		return false, err
	}

	return directoryChanged || permissionsChanged, nil
}

func (module *FileModule) absent(target fileInfo) (bool, error) {
	switch target.State {
	case Absent:
		return false, nil
	case File:
		err := os.Remove(target.Path)
		if err != nil {
			return false, errors.Wrapf(err, "failed to remove file %s", target.Path)
		}
		return true, nil
	case Directory:
		err := os.RemoveAll(target.Path)
		if err != nil {
			return false, errors.Wrapf(err, "failed to remove directory %s", target.Path)
		}
		return true, nil
	default:
		return false, errors.New("not yet implemented")
	}
}
