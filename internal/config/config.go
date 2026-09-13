package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	Filename      = "writeme.yaml"
	BaseFilename  = "BASE.md"
	EntryFilename = "WRITEME.md"
)

type Config struct {
	Entry  string   `yaml:"entry"`
	Parts  []string `yaml:"parts"`
	Output string   `yaml:"output"`
}

const defaultConfig = `# writeme.yaml
entry: WRITEME.md

# 등록하면 ::이름 으로 사용 가능
parts:
  - BASE.md

output: README.md
`

const defaultBase = `## License

MIT
`

const defaultEntry = `# Project Title

Describe your project here.

::BASE
`

// Init scaffolds a project in dir and reports which files it created.
// writeme.yaml must not already exist; BASE.md and WRITEME.md are only written
// when missing, so hand-written content is never clobbered.
func Init(dir string) ([]string, error) {
	configPath := filepath.Join(dir, Filename)

	exists, err := fileExists(configPath)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%s already exists", configPath)
	}

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0o644); err != nil {
		return nil, err
	}
	created := []string{Filename}

	scaffold := []struct {
		name    string
		content string
	}{
		{BaseFilename, defaultBase},
		{EntryFilename, defaultEntry},
	}

	for _, f := range scaffold {
		path := filepath.Join(dir, f.name)

		exists, err := fileExists(path)
		if err != nil {
			return created, err
		}
		if exists {
			continue
		}
		if err := os.WriteFile(path, []byte(f.content), 0o644); err != nil {
			return created, err
		}
		created = append(created, f.name)
	}

	return created, nil
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

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
