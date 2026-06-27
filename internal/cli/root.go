package cli

import (
	"github.com/spf13/cobra"
)

// version is overridden at build time via -ldflags.
var version = "dev"

// NewRootCmd builds the root `brrr` command and wires up its subcommands.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "brrr",
		Short:         "brrr — a CRUD code generator for Go + React",
		Long:          "brrr scaffolds and grows Go (modular monolith) + React applications.\n\nIt can initialise a new project, generate CRUD for a model, and apply a\nwhole-project YAML spec.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.CompletionOptions.HiddenDefaultCmd = true

	root.AddCommand(newInitCmd())
	root.AddCommand(newGenerateCmd())
	return root
}
