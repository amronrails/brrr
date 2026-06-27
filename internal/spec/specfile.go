package spec

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileSpec is a whole-project specification: an ordered set of modules, each
// with an ordered set of models, plus optional standard packages. Order is
// preserved (via yaml.Node) so generated columns and migrations are
// deterministic.
type FileSpec struct {
	Modules  []ModuleEntry
	Packages []string
}

// ModuleEntry is a module and its models, in declaration order.
type ModuleEntry struct {
	Name   string
	Models []ModelEntry
}

// ModelEntry is a single model definition from the spec.
type ModelEntry struct {
	Module        string
	Name          string
	Fields        []FieldEntry
	Relationships []RelEntry
}

// FieldEntry is a scalar field: a name, a type, and optional modifiers
// (required, unique).
type FieldEntry struct {
	Name      string
	Type      string
	Modifiers []string
}

// RelEntry is a relationship: a name, a kind (belongs_to, ...) and a target
// model name.
type RelEntry struct {
	Name   string
	Kind   string
	Target string
}

// ParseFileSpec reads and parses a YAML project spec, preserving field and
// model order.
func ParseFileSpec(path string) (*FileSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil, fmt.Errorf("spec file %s is empty", path)
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("spec root must be a mapping")
	}

	modulesNode := mapValue(root, "modules")
	if modulesNode == nil || modulesNode.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("spec must contain a 'modules' mapping")
	}

	fs := &FileSpec{}
	for i := 0; i+1 < len(modulesNode.Content); i += 2 {
		modName := modulesNode.Content[i].Value
		me := ModuleEntry{Name: modName}

		modelsNode := mapValue(modulesNode.Content[i+1], "models")
		if modelsNode == nil || modelsNode.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("module %q must contain a 'models' mapping", modName)
		}
		for j := 0; j+1 < len(modelsNode.Content); j += 2 {
			modelName := modelsNode.Content[j].Value
			modelBody := modelsNode.Content[j+1]
			model := ModelEntry{Module: modName, Name: modelName}

			if fnode := mapValue(modelBody, "fields"); fnode != nil {
				for k := 0; k+1 < len(fnode.Content); k += 2 {
					name := fnode.Content[k].Value
					toks := strings.Fields(fnode.Content[k+1].Value)
					if len(toks) == 0 {
						return nil, fmt.Errorf("%s.%s: field %q has no type", modName, modelName, name)
					}
					model.Fields = append(model.Fields, FieldEntry{Name: name, Type: toks[0], Modifiers: toks[1:]})
				}
			}
			if rnode := mapValue(modelBody, "relationships"); rnode != nil {
				for k := 0; k+1 < len(rnode.Content); k += 2 {
					name := rnode.Content[k].Value
					toks := strings.Fields(rnode.Content[k+1].Value)
					if len(toks) < 2 {
						return nil, fmt.Errorf("%s.%s: relationship %q must be '<kind> <Target>'", modName, modelName, name)
					}
					model.Relationships = append(model.Relationships, RelEntry{Name: name, Kind: toks[0], Target: toks[1]})
				}
			}
			me.Models = append(me.Models, model)
		}
		fs.Modules = append(fs.Modules, me)
	}

	if pkgs := mapValue(root, "packages"); pkgs != nil && pkgs.Kind == yaml.SequenceNode {
		for _, n := range pkgs.Content {
			fs.Packages = append(fs.Packages, n.Value)
		}
	}
	return fs, nil
}

// Args converts a ModelEntry into the generate-style "name:type[:mod...]" and
// "name:kind:Target" argument list consumed by ParseModel.
func (m ModelEntry) Args() []string {
	args := make([]string, 0, len(m.Fields)+len(m.Relationships))
	for _, f := range m.Fields {
		parts := append([]string{f.Name, f.Type}, f.Modifiers...)
		args = append(args, strings.Join(parts, ":"))
	}
	for _, r := range m.Relationships {
		args = append(args, strings.Join([]string{r.Name, r.Kind, r.Target}, ":"))
	}
	return args
}

func mapValue(m *yaml.Node, key string) *yaml.Node {
	if m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}
