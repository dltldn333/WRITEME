package main

import (
	"fmt"
	"os"

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

	fmt.Println("entries:")
	for _, e := range ws.Entries {
		fmt.Printf("  %s -> %s\n", e.Path, e.Output)
	}
	fmt.Println("parts:")
	for _, p := range ws.Parts {
		fmt.Printf("  ::%s  (%s)\n", p.Name, p.Path)
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
