package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dltldn333/WRITEME/internal/assemble"
	"github.com/dltldn333/WRITEME/internal/config"
	"github.com/dltldn333/WRITEME/internal/workspace"
)

func usage() {
	fmt.Fprintln(os.Stderr, `WRITEME — assemble README.md from WRITEME.md

Usage:
  writeme init    create writeme.yaml
  writeme build   compile entry into output`)
}

func runBuild() error {
	cfg, err := config.Load(".")
	if err != nil {
		return err
	}

	ws, err := workspace.Load(".", cfg)
	if err != nil {
		return err
	}

	a := assemble.New(ws)
	outputs := make([]string, len(ws.Entries))
	for i, e := range ws.Entries {
		if outputs[i], err = a.Entry(e); err != nil {
			return err
		}
	}

	// Write only after every entry assembled, so one bad file leaves nothing half-built.
	for i, e := range ws.Entries {
		if err := os.WriteFile(filepath.Join(ws.Root, e.Output), []byte(outputs[i]), 0o644); err != nil {
			return err
		}
		fmt.Println("Wrote", e.Output)
	}

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
