package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const Filename = "writeme.yaml"

type Config struct {
	Entry  string   `yaml:"entry"`
	Parts  []string `yaml:"parts"`
	Output string   `yaml:"output"`
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
	path := filepath.Join(dir, Filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	return cfg, nil
}
