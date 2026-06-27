package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ManifestFile is the on-disk name of the project manifest written at the root
// of a generated project. It is the source of truth for which modules and
// models exist, and is consumed by the generate and apply commands.
const ManifestFile = "brrr.yaml"

// Manifest is the persisted description of a generated project.
type Manifest struct {
	Name    string                    `yaml:"name"`
	Module  string                    `yaml:"module"`
	Modules map[string]ModuleManifest `yaml:"modules"`
}

// ModuleManifest records the models that belong to a single module.
type ModuleManifest struct {
	Models []string `yaml:"models,omitempty"`
}

// NewManifest creates a manifest seeded with the built-in user/auth module.
func NewManifest(c *Context) *Manifest {
	return &Manifest{
		Name:   c.AppName,
		Module: c.ModulePath,
		Modules: map[string]ModuleManifest{
			"user": {Models: []string{"User"}},
		},
	}
}

// LoadManifest reads and parses the manifest located in dir.
func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, ManifestFile))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", ManifestFile, err)
	}
	if m.Modules == nil {
		m.Modules = map[string]ModuleManifest{}
	}
	return &m, nil
}

// Save writes the manifest to dir.
func (m *Manifest) Save(dir string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ManifestFile), data, 0o644)
}

// AddModel records a model under a module, creating the module entry if needed
// and avoiding duplicates. It reports whether a change was made.
func (m *Manifest) AddModel(module, model string) bool {
	if m.Modules == nil {
		m.Modules = map[string]ModuleManifest{}
	}
	mm := m.Modules[module]
	for _, existing := range mm.Models {
		if existing == model {
			return false
		}
	}
	mm.Models = append(mm.Models, model)
	m.Modules[module] = mm
	return true
}
