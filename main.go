package main

import (
	"fmt"
	"os"

	"github.com/jamesl33/zk/cmd"
)

// main - TODO
func main() {
	err := cmd.Execute()
	if err == nil {
		return
	}

	// TODO
	defer os.Exit(1)

	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
}
