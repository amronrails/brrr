package cli

import (
	"fmt"

	"github.com/amronrails/brrr/internal/gen"
	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate <module> <Model> [field:type ...]",
		Aliases: []string{"g"},
		Short:   "Generate CRUD for a model (backend + frontend)",
		Long: "Generate all backend and frontend files for a model and wire them in.\n\n" +
			"Fields are name:type pairs; belongs_to relations are name:belongs_to:Target.\n" +
			"Add :required and/or :unique modifiers to a scalar field.\n\n" +
			"Example:\n  brrr g blog Post title:string:required body:text published:bool author:belongs_to:User",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := gen.FindProjectRoot(mustCwd())
			if err != nil {
				return err
			}

			module, model, fields := args[0], args[1], args[2:]
			res, err := gen.Generate(root, module, model, fields)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ generated %s in module %s (%d files)\n", res.Model, res.Module, len(res.Files))
			for _, f := range res.Files {
				fmt.Fprintf(out, "  + %s\n", f)
			}
			fmt.Fprintf(out, "\nNext steps:\n  make sqlc      # regenerate db code for the new queries\n  make migrate   # apply the new migration\n  make dev       # run it\n")
			return nil
		},
	}
	return cmd
}
