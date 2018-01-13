package files

import (
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
	// Link represents a symbolic link
	Link
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
	}

	file.FileMode = stat.Mode()
	sysStat, cast := stat.Sys().(*syscall.Stat_t)
	if !cast {
		return file, errors.New("stat not of type syscall.Stat_t")
	}

	file.UID = int(sysStat.Uid)
	file.GID = int(sysStat.Gid)

	hash, err := checksum(path)
	if err != nil {
		return file, err
	}

	file.Checksum = hash

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
