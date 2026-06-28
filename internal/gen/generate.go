package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amronrails/brrr/internal/engine"
	"github.com/amronrails/brrr/internal/project"
	"github.com/amronrails/brrr/internal/spec"
	"github.com/amronrails/brrr/internal/templates"
)

// Result summarises what a generate run did.
type Result struct {
	Files     []string // files written (relative to project root)
	NewModule bool
	Module    string
	Model     string
}

// Generate scaffolds CRUD for a model inside the project rooted at root. It
// writes the per-model backend and frontend files, (re)generates the module's
// wiring, injects route/registry/nav entries, and updates the manifest.
func Generate(root, module string, model string, fieldArgs []string) (*Result, error) {
	manifest, err := project.LoadManifest(root)
	if err != nil {
		return nil, fmt.Errorf("load %s (run inside a brrr project): %w", project.ManifestFile, err)
	}

	modelName := engine.Pascal(model)
	if existing, ok := manifest.Modules[module]; ok {
		for _, m := range existing.Models {
			if engine.Pascal(m) == modelName {
				return nil, fmt.Errorf("model %q already exists in module %q", modelName, module)
			}
		}
	}
	_, isExisting := manifest.Modules[module]
	newModule := !isExisting

	parsed, err := spec.ParseModel(module, modelName, fieldArgs)
	if err != nil {
		return nil, err
	}

	seq, err := nextMigrationSeq(root)
	if err != nil {
		return nil, err
	}

	view, err := BuildModelView(manifest.Module, parsed, seq)
	if err != nil {
		return nil, err
	}
	if len(view.Columns) == 0 {
		return nil, fmt.Errorf("model %q needs at least one field, e.g. %s name:string", modelName, strings.ToLower(modelName))
	}

	res := &Result{NewModule: newModule, Module: module, Model: modelName}

	// 1. Per-model backend + frontend leaf files (never overwrite existing).
	for _, sub := range []string{"backend/leaf", "frontend/leaf"} {
		tree, err := templates.Generate(sub)
		if err != nil {
			return nil, err
		}
		w := engine.NewWriter(root)
		if err := engine.Render(tree, view, w); err != nil {
			return nil, fmt.Errorf("render %s: %w", sub, err)
		}
		res.Files = append(res.Files, w.Written...)
	}

	// 2. Update the manifest in memory so module.go reflects the new model.
	manifest.AddModel(module, modelName)

	// 3. (Re)generate the module's wiring + public API files from the full model
	//    list, so other modules can depend on this one through its interface.
	if err := regenerateModule(root, manifest, module); err != nil {
		return nil, err
	}
	res.Files = append(res.Files,
		filepath.Join("internal/modules", view.ModulePkg, "module.go"),
		filepath.Join("internal/modules", view.ModulePkg, "api.go"),
	)

	// 4. Wire the module into the registry (only needed once per module).
	if newModule {
		if err := wireRegistry(root, view); err != nil {
			return nil, err
		}
	}

	// 5. Wire the frontend routes and navigation.
	if err := wireFrontend(root, view); err != nil {
		return nil, err
	}

	// 6. Persist the manifest.
	if err := manifest.Save(root); err != nil {
		return nil, err
	}

	return res, nil
}

func regenerateModule(root string, manifest *project.Manifest, module string) error {
	view := ModuleView{
		ModulePath: manifest.Module,
		Module:     module,
		ModulePkg:  pkgName(module),
		HTTPPkg:    pkgName(module) + "http",
	}
	for _, name := range manifest.Modules[module].Models {
		view.Models = append(view.Models, ModelRef{
			Pascal: engine.Pascal(name),
			Camel:  engine.Camel(name),
		})
	}
	// Both files are aggregates over the module's full model list: module.go
	// wires and mounts the slice; api.go declares the public contract other
	// modules depend on. Per-model detail lives in the leaf files.
	for tmpl, out := range map[string]string{
		"backend/module.go.tmpl": "module.go",
		"backend/api.go.tmpl":    "api.go",
	} {
		content, err := templates.GenerateFile(tmpl)
		if err != nil {
			return err
		}
		rendered, err := engine.RenderBytes(out, content, view)
		if err != nil {
			return err
		}
		dst := filepath.Join(root, "internal", "modules", view.ModulePkg, out)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, rendered, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func wireRegistry(root string, v *ModelView) error {
	registry := filepath.Join(root, "internal", "modules", "registry.go")
	if _, err := engine.Inject(registry, "// brrr:module-imports",
		fmt.Sprintf("%q", v.ModulePath+"/internal/modules/"+v.ModulePkg)); err != nil {
		return err
	}
	// Phase 1: construct the module into a named provider variable, so other
	// modules can be handed its public API at wiring time.
	construct := fmt.Sprintf("%sMod := %s.New(%s.Deps{Queries: d.Queries, Tokens: d.Tokens, Validator: d.Validator})", v.ModulePkg, v.ModulePkg, v.ModulePkg)
	if _, err := engine.Inject(registry, "// brrr:module-construct", construct); err != nil {
		return err
	}
	// Phase 2: add the provider to the mounted list.
	_, err := engine.Inject(registry, "// brrr:modules", v.ModulePkg+"Mod,")
	return err
}

func wireFrontend(root string, v *ModelView) error {
	router := filepath.Join(root, "web", "src", "router.tsx")
	imports := fmt.Sprintf(
		"import { %sListPage, %sFormPage } from \"@/features/%s/%s\";",
		v.Pascal, v.Pascal, v.ModulePkg, v.Snake,
	)
	if _, err := engine.Inject(router, "// brrr:imports-fe", imports); err != nil {
		return err
	}
	routes := fmt.Sprintf(
		"{ path: %q, element: <%sListPage /> },\n{ path: %q, element: <%sFormPage /> },\n{ path: %q, element: <%sFormPage /> },",
		v.RoutePrefix, v.Pascal,
		v.RoutePrefix+"/new", v.Pascal,
		v.RoutePrefix+"/:id/edit", v.Pascal,
	)
	if _, err := engine.Inject(router, "// brrr:routes-fe", routes); err != nil {
		return err
	}

	layout := filepath.Join(root, "web", "src", "components", "layout", "DashboardLayout.tsx")
	nav := fmt.Sprintf("<NavLink to=\"/%s\" className={navItemClass}>%s</NavLink>", v.RoutePrefix, engine.Title(v.Plural))
	_, err := engine.Inject(layout, "{/* brrr:nav */}", nav)
	return err
}

// nextMigrationSeq returns one past the highest numeric prefix found in
// db/migrations.
func nextMigrationSeq(root string) (int, error) {
	dir := filepath.Join(root, "db", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read migrations dir: %w", err)
	}
	max := 0
	for _, e := range entries {
		name := e.Name()
		idx := strings.IndexByte(name, '_')
		if idx <= 0 {
			continue
		}
		if n, err := strconv.Atoi(name[:idx]); err == nil && n > max {
			max = n
		}
	}
	return max + 1, nil
}

// FindProjectRoot walks up from start looking for a brrr.yaml manifest.
func FindProjectRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, project.ManifestFile)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside a brrr project (no %s found)", project.ManifestFile)
		}
		dir = parent
	}
}
