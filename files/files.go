package files

import (
	"io/ioutil"
	"os"
	"syscall"

	"fmt"
	"io"

	"hash"

	"crypto/sha256"

	"github.com/pkg/errors"
)

// State represents the state of a file
type State int

const (
	// File represents a normal file
	File = iota
	// Directory represents a directory
	Directory
	// Absent represents the absent of the path
	Absent
)

type fileInfo struct {
	permissions
	Path     string
	State    State
	Checksum string
}

type permissions struct {
	FileMode os.FileMode
	UID      int
	GID      int
}

func collectAndMergeFileInfo(path string, settings permissions) (fileInfo, error) {
	info, err := collectFileInfo(path)
	if err != nil {
		return info, err
	}
	return mergeFilePermissions(info, settings), nil
}

func collectFileInfo(path string) (fileInfo, error) {
	file := fileInfo{Path: path}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			file.State = Absent
			return file, nil
		}
		return file, errors.Wrapf(err, "failed to stat %s", path)
	}

	if stat.IsDir() {
		file.State = Directory
	} else {
		file.State = File

		hash, err := checksum(path)
		if err != nil {
			return file, err
		}

		file.Checksum = hash
	}

	file.FileMode = stat.Mode().Perm()
	sysStat, cast := stat.Sys().(*syscall.Stat_t)
	if !cast {
		return file, errors.New("stat not of type syscall.Stat_t")
	}

	file.UID = int(sysStat.Uid)
	file.GID = int(sysStat.Gid)

	return file, nil
}

func mergeFilePermissions(info fileInfo, settings permissions) fileInfo {
	if settings.FileMode > 0 {
		info.FileMode = settings.FileMode
	}

	if settings.UID >= 0 {
		info.UID = settings.UID
	}

	if settings.GID >= 0 {
		info.GID = settings.GID
	}
	return info
}

func checksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open file %s", path)
	}

	defer file.Close()

	hashAlg := createHashAlg()
	if _, err := io.Copy(hashAlg, file); err != nil {
		return "", errors.Wrapf(err, "failed to create hashAlg of %s", path)
	}

	return hashToString(hashAlg), nil
}

func createHashAlg() hash.Hash {
	return sha256.New()
}

func hashToString(hashAlg hash.Hash) string {
	return fmt.Sprintf("%x", hashAlg.Sum(nil))
}

func ensureContent(target fileInfo, content string, mode os.FileMode) (bool, error) {
	bytes := []byte(content)
	if target.State == File {
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
	} else if target.State == Absent {
		err := ioutil.WriteFile(target.Path, bytes, mode)
		if err != nil {
			return false, errors.Wrapf(err, "failed to write content to %s", target.Path)
		}
		return true, nil
	} else {
		return false, errors.Errorf("%s seams to be not a regular file", target.Path)
	}
	return false, nil
}

func ensurePermissions(expected permissions, target fileInfo) (bool, error) {
	modeChanged := false
	if expected.FileMode != target.FileMode {
		err := os.Chmod(target.Path, expected.FileMode)
		if err != nil {
			return false, errors.Wrapf(err, "failed to change mode of %s", target.Path)
		}
		modeChanged = true
	}

	ownershipChanged := false
	if expected.UID != target.UID || expected.GID != target.GID {
		err := os.Chown(target.Path, expected.UID, expected.GID)
		if err != nil {
			return false, errors.Wrapf(err, "failed to change mode of %s", target.Path)
		}
		ownershipChanged = true
	}
	return modeChanged || ownershipChanged, nil
}
