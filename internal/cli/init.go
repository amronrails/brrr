package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amronrails/brrr/internal/engine"
	"github.com/amronrails/brrr/internal/gen"
	"github.com/amronrails/brrr/internal/project"
	"github.com/amronrails/brrr/internal/spec"
	"github.com/amronrails/brrr/internal/templates"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var (
		modulePath string
		dir        string
		specPath   string
		force      bool
		dryRun     bool
	)
	cmd := &cobra.Command{
		Use:   "init <app-name>",
		Short: "Scaffold a new Go + React project",
		Long: "Initialise a new project: a Go modular monolith with a secure\n" +
			"user/auth module plus a React admin frontend under web/.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if specPath != "" && dryRun {
				return fmt.Errorf("--spec cannot be combined with --dry-run")
			}
			ctx, err := project.NewContext(appName, modulePath)
			if err != nil {
				return err
			}

			// Parse the spec up front so a bad spec fails before we write files.
			var fileSpec *spec.FileSpec
			if specPath != "" {
				fileSpec, err = spec.ParseFileSpec(specPath)
				if err != nil {
					return err
				}
			}

			dest := dir
			if dest == "" {
				dest = appName
			}
			dest, err = filepath.Abs(dest)
			if err != nil {
				return err
			}

			tmpl, err := templates.Init()
			if err != nil {
				return err
			}

			w := engine.NewWriter(dest)
			w.Force = force
			w.DryRun = dryRun
			if err := engine.Render(tmpl, ctx, w); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			rel, _ := filepath.Rel(mustCwd(), dest)
			fmt.Fprintf(out, "✓ initialised %s (%s)\n", appName, w.Summary())

			if !dryRun {
				if err := project.NewManifest(ctx).Save(dest); err != nil {
					return err
				}
			}

			if fileSpec != nil {
				applied, err := gen.ApplySpec(dest, fileSpec)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "\nApplied spec — generated %d model(s):\n", len(applied.Models))
				for _, r := range applied.Models {
					fmt.Fprintf(out, "  • %s/%s\n", r.Module, r.Model)
				}
				if len(applied.Packages) > 0 {
					fmt.Fprintf(out, "\nNote: standard packages %v are recognised but not generated yet.\n", applied.Packages)
				}
			}

			fmt.Fprintf(out, "\nNext steps:\n  cd %s\n  make setup     # create .env, generate sqlc, install deps\n  make db-up     # start Postgres (docker)\n  make migrate   # apply migrations\n  make dev       # run backend + frontend\n\nEdit .env to set JWT_SECRET before deploying.\n", rel)
			return nil
		},
	}
	cmd.Flags().StringVarP(&modulePath, "module", "m", "", "Go module path (default: app name)")
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "target directory (default: ./<app-name>)")
	cmd.Flags().StringVarP(&specPath, "spec", "s", "", "YAML spec of modules/models to generate after scaffolding")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be generated without writing")
	return cmd
}

func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
