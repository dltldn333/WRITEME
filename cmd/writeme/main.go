package main

import (
	"fmt"
	"os"

	"github.com/dltldn333/WRITEME/internal/config"
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
			if err := config.Init("." ); err != nil{
			fmt.Fprintln(os.Stderr, "writeme:", err)
			os.Exit(1)
		}
		fmt.Println("Created", config.Filename)
	case "build":
		// config.Load(".")
	default:
		usage()
		os.Exit(1)
	}
}