package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesFiles(t *testing.T) {
	dir := t.TempDir()

	created, err := Init(dir)
	if err != nil {
		t.Fatalf("Init() unexpected error: %v", err)
	}
	want := []string{Filename, BaseFilename, EntryFilename}
	if len(created) != len(want) {
		t.Errorf("created = %v, want %v", created, want)
	}
	for _, name := range want {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("%s was not created: %v", name, err)
		}
	}
}

func TestInitKeepsExistingFiles(t *testing.T) {
	for _, name := range []string{BaseFilename, EntryFilename} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, name)

			if err := os.WriteFile(path, []byte("mine\n"), 0o644); err != nil {
				t.Fatal(err)
			}

			created, err := Init(dir)
			if err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
			for _, c := range created {
				if c == name {
					t.Errorf("%s should not have been rewritten", name)
				}
			}

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "mine\n" {
				t.Errorf("%s was overwritten: %q", name, got)
			}
		})
	}
}

// Init should leave a directory that "writeme build" can run in as-is.
func TestInitProducesBuildableProject(t *testing.T) {
	dir := t.TempDir()

	if _, err := Init(dir); err != nil {
		t.Fatalf("Init() unexpected error: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, cfg.Entry)); err != nil {
		t.Errorf("entry %s from config does not exist: %v", cfg.Entry, err)
	}
	for _, part := range cfg.Parts {
		if _, err := os.Stat(filepath.Join(dir, part)); err != nil {
			t.Errorf("part %s from config does not exist: %v", part, err)
		}
	}
}

func TestInitRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()

	if _, err := Init(dir); err != nil {
		t.Fatalf("first Init() unexpected error: %v", err)
	}
	if _, err := Init(dir); err == nil {
		t.Fatal("second Init() should fail when the file already exists")
	}
}

func TestLoadReadsDefaults(t *testing.T) {
	dir := t.TempDir()

	if _, err := Init(dir); err != nil {
		t.Fatalf("Init() unexpected error: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Entry != "WRITEME.md" {
		t.Errorf("Entry = %q, want %q", cfg.Entry, "WRITEME.md")
	}
	if cfg.Output != "README.md" {
		t.Errorf("Output = %q, want %q", cfg.Output, "README.md")
	}
	if len(cfg.Parts) != 1 || cfg.Parts[0] != BaseFilename {
		t.Errorf("Parts = %v, want [%s]", cfg.Parts, BaseFilename)
	}
}

func TestLoadParts(t *testing.T) {
	dir := t.TempDir()
	body := "entry: WRITEME.md\nparts:\n  - BASE.md\n  - install.md\noutput: README.md\n"

	if err := os.WriteFile(filepath.Join(dir, Filename), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	want := []string{"BASE.md", "install.md"}
	if len(cfg.Parts) != len(want) {
		t.Fatalf("Parts = %v, want %v", cfg.Parts, want)
	}
	for i := range want {
		if cfg.Parts[i] != want[i] {
			t.Errorf("Parts[%d] = %q, want %q", i, cfg.Parts[i], want[i])
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(t.TempDir()); err == nil {
		t.Fatal("Load() should fail when writeme.yaml does not exist")
	}
}
