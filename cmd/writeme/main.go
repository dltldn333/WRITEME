package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dltldn333/WRITEME/internal/config"
)

func usage() {
	fmt.Fprintln(os.Stderr, `WRITEME — assemble README.md from WRITEME.md

Usage:
  writeme init    create writeme.yaml
  writeme build   compile entry into output`)
}

// partName turns "docs/BASE.md" into the component name "BASE".
func partName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func runBuild() error {
	cfg, err := config.Load(".")
	if err != nil {
		return err
	}

	entry, err := os.ReadFile(cfg.Entry)
	if err != nil {
		return fmt.Errorf("entry %s: %w", cfg.Entry, err)
	}

	parts := make(map[string]string)
	for _, p := range cfg.Parts {
		data, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("part %s: %w", p, err)
		}

		name := partName(p)
		if _, dup := parts[name]; dup {
			return fmt.Errorf("duplicate part name %q", name)
		}
		parts[name] = string(data)
	}

	fmt.Printf("entry:  %s (%d bytes)\n", cfg.Entry, len(entry))
	for _, p := range cfg.Parts {
		name := partName(p)
		fmt.Printf("part:   %s (%d bytes)\n", name, len(parts[name]))
	}
	fmt.Printf("output: %s\n", cfg.Output)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		created, err := config.Init(".")
		if err != nil {
			fmt.Fprintln(os.Stderr, "writeme:", err)
			os.Exit(1)
		}
		for _, name := range created {
			fmt.Println("Created", name)
		}
	case "build":
		if err := runBuild(); err != nil {
			fmt.Fprintln(os.Stderr, "writeme:", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}
