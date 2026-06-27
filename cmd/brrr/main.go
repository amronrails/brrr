// Command brrr is a CRUD code generator for Go (modular monolith) + React.
package main

import (
	"fmt"
	"os"

	"github.com/amronrails/brrr/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
