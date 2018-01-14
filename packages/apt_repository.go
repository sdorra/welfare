package packages

import (
	"io/ioutil"
	"path"

	"os"

	"bufio"

	"io"
	"strings"

	"path/filepath"

	"github.com/pkg/errors"
)

// NewAptRepositoryModule creates a new module for managing apt repositories
func NewAptRepositoryModule(name string, state State) *AptRepositoryModule {
	return &AptRepositoryModule{
		Name:      name,
		State:     state,
		Directory: "/etc/apt",
	}
}

// AptRepositoryModule is able to handle state of apt repositories
type AptRepositoryModule struct {
	Name       string
	Repository string
	State      State
	Directory  string
}

func (module *AptRepositoryModule) Run() (bool, error) {
	present, err := module.isPresent()
	if err != nil {
		return false, err
	}

	if module.State == Present && !present {
		err := module.registerRepository()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if module.State == Absent && present {
		return false, errors.New("not yet implemented")
	}

	return false, nil
}

func (module *AptRepositoryModule) registerRepository() error {
	sourcesDirectory := path.Join(module.Directory, "sources.list.d")
	if _, err := os.Stat(sourcesDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(sourcesDirectory, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory %s", sourcesDirectory)
		}
	}

	file := path.Join(module.Directory, "sources.list.d", module.Name+".list")

	content := "# " + module.Name + " repository created by welfare\n"
	content += module.Repository + "\n"
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write repository file %s", file)
	}

	return nil
}

func (module *AptRepositoryModule) isPresent() (bool, error) {
	sourcesList := path.Join(module.Directory, "sources.list")
	contains, err := containsRepository(sourcesList, module.Repository)
	if err != nil {
		return false, err
	}
	if contains {
		return true, nil
	}

	sourcesDirectory := path.Join(module.Directory, "sources.list.d")
	if _, err := os.Stat(sourcesDirectory); os.IsNotExist(err) {
		return false, nil
	}

	sourceFiles, err := filepath.Glob(path.Join(sourcesDirectory, "*.list"))
	if err != nil {
		return false, errors.Wrapf(err, "failed to list source files from directory %s", sourcesDirectory)
	}

	for _, sourceFile := range sourceFiles {
		contains, err := containsRepository(sourceFile, module.Repository)
		if err != nil {
			return false, err
		}
		if contains {
			return true, nil
		}
	}

	return false, nil
}

func containsRepository(path string, repository string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, errors.Wrapf(err, "failed to open file %s", path)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.TrimSpace(line) == repository {
			return true, nil
		}
	}

	if err != nil && err != io.EOF {
		return false, errors.Wrapf(err, "failed to read line from file %s", path)
	}

	return false, nil
}
