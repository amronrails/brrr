package gen

import (
	"fmt"

	"github.com/amronrails/brrr/internal/engine"
	"github.com/amronrails/brrr/internal/spec"
)

// ApplyResult summarises an ApplySpec run.
type ApplyResult struct {
	Models   []*Result
	Packages []string // requested standard packages (not yet generated)
}

// ApplySpec generates every model described by fs into the project at root, in
// dependency order so that belongs-to targets are created before their
// dependents. It drives the same pipeline as `brrr generate`.
func ApplySpec(root string, fs *spec.FileSpec) (*ApplyResult, error) {
	var models []spec.ModelEntry
	for _, mod := range fs.Modules {
		models = append(models, mod.Models...)
	}

	ordered, err := topoSortModels(models)
	if err != nil {
		return nil, err
	}

	res := &ApplyResult{Packages: fs.Packages}
	for _, m := range ordered {
		r, err := Generate(root, m.Module, m.Name, m.Args())
		if err != nil {
			return nil, fmt.Errorf("module %s, model %s: %w", m.Module, m.Name, err)
		}
		res.Models = append(res.Models, r)
	}
	return res, nil
}

// topoSortModels orders models so that any model is generated after the models
// it belongs_to. Targets that are not defined in the spec are assumed to
// already exist (the built-in User, or models from a prior run) and impose no
// ordering. Input order is preserved among independent models.
func topoSortModels(models []spec.ModelEntry) ([]spec.ModelEntry, error) {
	index := make(map[string]spec.ModelEntry, len(models))
	order := make([]string, 0, len(models))
	known := map[string]bool{"User": true}

	for _, m := range models {
		p := engine.Pascal(m.Name)
		if _, dup := index[p]; dup {
			return nil, fmt.Errorf("duplicate model %q in spec", m.Name)
		}
		index[p] = m
		order = append(order, p)
		known[p] = true
	}

	deps := make(map[string][]string)
	for _, m := range models {
		p := engine.Pascal(m.Name)
		for _, r := range m.Relationships {
			t := engine.Pascal(r.Target)
			if !known[t] {
				return nil, fmt.Errorf("model %q references unknown target %q (not defined in the spec and not the built-in User)", m.Name, r.Target)
			}
			if _, ok := index[t]; ok && t != p {
				deps[p] = append(deps[p], t)
			}
		}
	}

	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)
	state := make(map[string]int, len(order))
	out := make([]spec.ModelEntry, 0, len(order))

	var visit func(p string) error
	visit = func(p string) error {
		switch state[p] {
		case done:
			return nil
		case visiting:
			return fmt.Errorf("relationship cycle involving %q", index[p].Name)
		}
		state[p] = visiting
		for _, d := range deps[p] {
			if err := visit(d); err != nil {
				return err
			}
		}
		state[p] = done
		out = append(out, index[p])
		return nil
	}

	for _, p := range order {
		if err := visit(p); err != nil {
			return nil, err
		}
	}
	return out, nil
}
