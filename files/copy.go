package files

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// NewCopyModule creates a new CopyModule with the given source and target file
func NewCopyModule(source string, target string) *CopyModule {
	module := &CopyModule{
		Source: source,
		Target: target,
	}
	module.FileMode = os.FileMode(0)
	module.UID = -1
	module.GID = -1
	return module
}

// CopyModule ensures that the target is an exact copy of the source file
type CopyModule struct {
	permissions
	Source string
	Target string
}

func (module *CopyModule) Run() (bool, error) {
	expected, err := collectAndMergeFileInfo(module.Source, module.permissions)
	if err != nil {
		return false, err
	}

	if expected.State != File {
		return false, errors.Errorf("expected file %s seams to be not a file", module.Source)
	}

	target, err := collectFileInfo(module.Target)
	if err != nil {
		return false, err
	}

	return ensureCopy(expected, target)
}

func ensureCopy(expected, target fileInfo) (bool, error) {
	contentChanged := false
	if target.State == Absent || expected.Checksum != target.Checksum {
		err := copy(expected, target)
		if err != nil {
			return false, err
		}
		contentChanged = true
	}

	permissionsChanged, err := ensurePermissions(expected.permissions, target)
	if err != nil {
		return false, err
	}

	return contentChanged || permissionsChanged, nil
}

func copy(expected, target fileInfo) error {
	sourcePath := expected.Path
	targetPath := target.Path

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open sourceFile %s", sourcePath)
	}

	var targetFile *os.File
	if target.State == Absent {
		targetFile, err = os.Create(targetPath)
		if err != nil {
			return errors.Wrapf(err, "failed to create file at %s", targetPath)
		}
	} else {
		targetFile, err = os.OpenFile(targetPath, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return errors.Wrapf(err, "failed to open file at %s", targetPath)
		}
	}

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", sourcePath, targetPath)
	}

	return nil
}
