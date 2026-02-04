package main

import (
	"flag"
	"os"
	"fmt"
)

func main() {
	target := flag.String("file", "README.fmd", "Target .fmd file to assemble")
	flag.Parse()

	fmt.Println("WRITEME: Assembling...")
	fmt.Printf("Target Source: %s\n", *target)

	// TODO: parse logic here
	
	if _, err := os.Stat(*target); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", *target)
		return
	}
	fmt.Println("Done! (Just kidding, implementation coming soon)")
}