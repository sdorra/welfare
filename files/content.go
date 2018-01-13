package files

import (
	"os"

	"io/ioutil"

	"github.com/pkg/errors"
)

// NewContentModule creates a new ContentModule for the given target and content
func NewContentModule(target string, content string) *ContentModule {
	module := &ContentModule{
		Target:  target,
		Content: content,
	}
	module.FileMode = os.FileMode(0644)
	module.UID = os.Getuid()
	module.GID = os.Getegid()
	return module
}

// ContentModule ensures that the target file exists with the specified content
type ContentModule struct {
	permissions
	Target  string
	Content string
}

func (module *ContentModule) Run() (bool, error) {
	target, err := collectFileInfo(module.Target)
	if err != nil {
		return false, err
	}

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

func ensureContent(target fileInfo, content string, mode os.FileMode) (bool, error) {
	bytes := []byte(content)
	if target.Exists {
		hashAlg := createHashAlg()
		_, err := hashAlg.Write(bytes)
		if err != nil {
			return false, errors.Wrap(err, "failed to create checksum for content")
		}
		hash := hashToString(hashAlg)
		if hash != target.Checksum {
			err := ioutil.WriteFile(target.Path, bytes, mode)
			if err != nil {
				return false, errors.Wrapf(err, "failed to overwrite content of %s", target.Path)
			}
			return true, nil
		}
	} else {
		err := ioutil.WriteFile(target.Path, bytes, mode)
		if err != nil {
			return false, errors.Wrapf(err, "failed to write content to %s", target.Path)
		}
		return true, nil
	}
	return false, nil
}
