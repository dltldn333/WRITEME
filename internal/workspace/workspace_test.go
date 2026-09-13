package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dltldn333/WRITEME/internal/config"
)

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func defaultConfig(parts ...string) config.Config {
	return config.Config{Entry: "WRITEME.md", Parts: parts, Output: "README.md"}
}

func TestLoadFindsEntriesInEveryDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "WRITEME.md", "root")
	writeFile(t, root, "packages/core/WRITEME.md", "core")
	writeFile(t, root, "packages/cli/nested/WRITEME.md", "nested")
	writeFile(t, root, "packages/cli/notes.md", "not an entry")

	ws, err := Load(root, defaultConfig())
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := []Entry{
		{Dir: ".", Path: "WRITEME.md", Output: "README.md", Content: "root"},
		{
			Dir:     filepath.Join("packages", "cli", "nested"),
			Path:    filepath.Join("packages", "cli", "nested", "WRITEME.md"),
			Output:  filepath.Join("packages", "cli", "nested", "README.md"),
			Content: "nested",
		},
		{
			Dir:     filepath.Join("packages", "core"),
			Path:    filepath.Join("packages", "core", "WRITEME.md"),
			Output:  filepath.Join("packages", "core", "README.md"),
			Content: "core",
		},
	}

	if len(ws.Entries) != len(want) {
		t.Fatalf("found %d entries, want %d: %+v", len(ws.Entries), len(want), ws.Entries)
	}
	for i := range want {
		if ws.Entries[i] != want[i] {
			t.Errorf("Entries[%d]\n got: %+v\nwant: %+v", i, ws.Entries[i], want[i])
		}
	}
}

func TestLoadSkipsIgnoredDirectories(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "WRITEME.md", "root")
	writeFile(t, root, "node_modules/dep/WRITEME.md", "dep")
	writeFile(t, root, ".git/WRITEME.md", "git")
	writeFile(t, root, "vendor/mod/WRITEME.md", "vendored")

	ws, err := Load(root, defaultConfig())
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if len(ws.Entries) != 1 || ws.Entries[0].Path != "WRITEME.md" {
		t.Errorf("Entries = %+v, want only the root WRITEME.md", ws.Entries)
	}
}

func TestLoadNoEntries(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "hand written")

	if _, err := Load(root, defaultConfig()); err == nil {
		t.Fatal("Load() should fail when no WRITEME.md exists")
	}
}

func TestLoadPartsAtRootAndInFolders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "WRITEME.md", "::BASE")
	writeFile(t, root, "BASE.md", "base")
	writeFile(t, root, "parts/install.md", "install")
	writeFile(t, root, "docs/shared/license.md", "license")

	ws, err := Load(root, defaultConfig("BASE.md", "parts/install.md", "docs/shared/license.md"))
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := []Part{
		{Name: "BASE", Path: "BASE.md", Content: "base"},
		{Name: "install", Path: "parts/install.md", Content: "install"},
		{Name: "license", Path: "docs/shared/license.md", Content: "license"},
	}
	if len(ws.Parts) != len(want) {
		t.Fatalf("Parts = %+v, want %+v", ws.Parts, want)
	}
	for i := range want {
		if ws.Parts[i] != want[i] {
			t.Errorf("Parts[%d]\n got: %+v\nwant: %+v", i, ws.Parts[i], want[i])
		}
	}

	if p, ok := ws.Part("install"); !ok || p.Content != "install" {
		t.Errorf(`Part("install") = %+v, %v`, p, ok)
	}
	if _, ok := ws.Part("missing"); ok {
		t.Error(`Part("missing") should not be found`)
	}
}

func TestLoadPartErrors(t *testing.T) {
	tests := []struct {
		name    string
		files   map[string]string
		parts   []string
		errText string
	}{
		{
			name:    "missing file",
			parts:   []string{"parts/gone.md"},
			errText: "parts/gone.md",
		},
		{
			name:    "duplicate name",
			files:   map[string]string{"BASE.md": "a", "parts/BASE.md": "b"},
			parts:   []string{"BASE.md", "parts/BASE.md"},
			errText: `"BASE"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "WRITEME.md", "entry")
			for rel, content := range tt.files {
				writeFile(t, root, rel, content)
			}

			_, err := Load(root, defaultConfig(tt.parts...))
			if err == nil {
				t.Fatal("Load() expected an error, got none")
			}
			if !strings.Contains(err.Error(), tt.errText) {
				t.Errorf("error %q should mention %s", err, tt.errText)
			}
		})
	}
}

func TestLoadRejectsPathsForEntryAndOutput(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.Config
	}{
		{"entry with directory", config.Config{Entry: "docs/WRITEME.md", Output: "README.md"}},
		{"output with directory", config.Config{Entry: "WRITEME.md", Output: "out/README.md"}},
		{"entry empty", config.Config{Output: "README.md"}},
		{"output empty", config.Config{Entry: "WRITEME.md"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "WRITEME.md", "entry")

			if _, err := Load(root, tt.cfg); err == nil {
				t.Fatal("Load() expected an error, got none")
			}
		})
	}
}
