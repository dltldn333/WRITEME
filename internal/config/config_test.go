package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesFile(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("Init() unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, Filename)); err != nil {
		t.Fatalf("%s was not created: %v", Filename, err)
	}
}

func TestInitRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("first Init() unexpected error: %v", err)
	}
	if err := Init(dir); err == nil {
		t.Fatal("second Init() should fail when the file already exists")
	}
}

func TestLoadReadsDefaults(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
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
	if len(cfg.Parts) != 0 {
		t.Errorf("Parts = %v, want empty", cfg.Parts)
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
