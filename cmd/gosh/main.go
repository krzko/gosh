// cmd/gosh/main.go
package main

import (
	"fmt"
	"os"

	"github.com/krzko/gosh/internal/shell"
)

func main() {
	sh, err := shell.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing shell: %v\n", err)
		os.Exit(1)
	}

	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
