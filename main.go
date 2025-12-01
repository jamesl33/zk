package main

import (
	"fmt"
	"os"

	"github.com/jamesl33/zk/cmd"
)

// main runs 'zk'.
func main() {
	err := cmd.Execute()
	if err == nil {
		return
	}

	defer os.Exit(1)

	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
}
