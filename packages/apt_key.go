package packages

import (
	"bufio"
	"io"

	"strings"

	"regexp"

	"os/exec"

	"bytes"

	"github.com/pkg/errors"
)

// NewAptKeyModule creates a new module for managing apt repository keys
func NewAptKeyModule(id string, state State) *AptKeyModule {
	return &AptKeyModule{
		ID:     id,
		Server: "hkp://keyserver.ubuntu.com:80",
		State:  state,
		system: &aptKey{},
	}
}

// AptKeyModule handles the state of apt repository keys
type AptKeyModule struct {
	ID     string
	Server string
	State  State
	system keySystem
}

func (module *AptKeyModule) Run() (bool, error) {
	present, err := module.system.IsPresent(module.ID)
	if err != nil {
		return false, err
	}

	if module.State == Present && !present {
		err := module.system.Add(module.Server, module.ID)
		if err != nil {
			return false, err
		}
		return true, nil
	} else if module.State == Absent && present {
		err := module.system.Remove(module.ID)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

type keySystem interface {
	Add(server, key string) error
	Remove(key string) error
	IsPresent(key string) (bool, error)
}

type aptKey struct {
}

func (sys *aptKey) Add(server string, id string) error {
	err := exec.Command("apt-key", "adv", "--recv-keys", "--keyserver", server, id).Run()
	if err != nil {
		return errors.Wrapf(err, "failed to add key %s from server %s", id, server)
	}
	return nil
}

func (sys *aptKey) Remove(id string) error {
	err := exec.Command("apt-key", "del", id).Run()
	if err != nil {
		return errors.Wrapf(err, "failed to remove key %s", id)
	}
	return nil
}

func (sys *aptKey) IsPresent(id string) (bool, error) {
	cmd := exec.Command("apt-key", "list")
	listing, err := cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to list keys")
	}

	contains, err := containsKey(bytes.NewReader(listing), id)
	if err != nil {
		return false, err
	}
	return contains, nil
}

var keyLine = regexp.MustCompile("pub\\s+[0-9A-Z]+/([0-9A-Z]+) [0-9]{4}-[0-9]{2}-[0-9]{2}")

func containsKey(listing io.Reader, key string) (bool, error) {
	reader := bufio.NewReader(listing)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, errors.Wrapf(err, "failed to check if listing contains key %s", key)
		}

		groups := keyLine.FindStringSubmatch(strings.TrimSpace(line))
		if len(groups) == 2 && groups[1] == key {
			return true, nil
		}

	}
	return false, nil
}
