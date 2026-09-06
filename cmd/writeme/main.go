package main

import (
	"fmt"
	"os"
)

func usage(){
		fmt.Fprintln(os.Stderr, `WRITEME — assemble README.md from WRITEME.md

Usage:
  writeme init    create writeme.yaml
  writeme build   compile entry into output`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		//config.Init(".")
	case "build":
		// config.Load(".")
	default:
		usage()
		os.Exit(1)
	}
}