package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const Filename = "writeme.yaml"

type Config struct {
	Entry	string		`yaml:"entry"`
	Parts	[]string	`yaml:"parts"`
	Output	string		`yaml:"output"`
}
const defaultConfig = `# writeme.yaml
entry: WRITEME.md
parts: []
output: README.md
`

func Init(dir string) error {
	path := filepath.Join(dir, Filename)
	_, err := os.Stat(path)

	if err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
	return err
	}
	return os.WriteFile(path, []byte(defaultConfig), 0o644)
}

func Load(dir string) (Config, error) {
	path := filepath.Join(ddir, Filename)

	data, err := 
	
}