package main

import (
	"fmt"
	"os"

	"github.com/ez0000001000000/Goop/cmd/goop/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
