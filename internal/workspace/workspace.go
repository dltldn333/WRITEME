package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dltldn333/WRITEME/internal/config"
)

// skipDirs are never searched for entry files.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
}

// Entry is one WRITEME.md found under the root. All paths are relative to the root.
type Entry struct {
	Dir     string // directory holding the entry, "." for the root itself
	Path    string // path to the entry file
	Output  string // where its README.md will be written
	Content string
}

// Part is a shared fragment registered in writeme.yaml, called as ::Name.
type Part struct {
	Name    string
	Path    string // path relative to the root, as written in writeme.yaml
	Content string
}

type Workspace struct {
	Root    string
	Entries []Entry
	Parts   []Part
}

// Load reads every part listed in cfg and every entry file found under root.
func Load(root string, cfg config.Config) (*Workspace, error) {
	if err := checkFileName("entry", cfg.Entry); err != nil {
		return nil, err
	}
	if err := checkFileName("output", cfg.Output); err != nil {
		return nil, err
	}

	parts, err := loadParts(root, cfg.Parts)
	if err != nil {
		return nil, err
	}

	entries, err := findEntries(root, cfg.Entry, cfg.Output)
	if err != nil {
		return nil, err
	}

	return &Workspace{Root: root, Entries: entries, Parts: parts}, nil
}

// Part looks up a registered part by its component name.
func (w *Workspace) Part(name string) (Part, bool) {
	for _, p := range w.Parts {
		if p.Name == name {
			return p, true
		}
	}
	return Part{}, false
}

// entry and output are matched in every directory, so they must be bare file names.
func checkFileName(field, name string) error {
	if name == "" {
		return fmt.Errorf("%s is not set in %s", field, config.Filename)
	}
	if filepath.Base(name) != name {
		return fmt.Errorf("%s must be a file name, not a path: %q", field, name)
	}
	return nil
}

func loadParts(root string, paths []string) ([]Part, error) {
	parts := make([]Part, 0, len(paths))
	seen := make(map[string]string)

	for _, p := range paths {
		data, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return nil, fmt.Errorf("part %s: %w", p, err)
		}

		name := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		if prev, dup := seen[name]; dup {
			return nil, fmt.Errorf("part name %q is used by both %s and %s", name, prev, p)
		}
		seen[name] = p

		parts = append(parts, Part{Name: name, Path: p, Content: string(data)})
	}
	return parts, nil
}

func findEntries(root, entryName, outputName string) ([]Entry, error) {
	var entries []Entry

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != entryName {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("entry %s: %w", rel, err)
		}

		dir := filepath.Dir(rel)
		entries = append(entries, Entry{
			Dir:     dir,
			Path:    rel,
			Output:  filepath.Join(dir, outputName),
			Content: string(data),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no %s found under %s", entryName, root)
	}
	return entries, nil
}
