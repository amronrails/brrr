// Package templates embeds the code-generation template trees that ship with
// brrr. Every template file carries a ".tmpl" suffix so that template Go and
// TypeScript sources are never compiled or type-checked as part of brrr itself.
package templates

import (
	"embed"
	"io/fs"
)

//go:embed all:init
var initFS embed.FS

//go:embed all:generate
var generateFS embed.FS

// Init returns the template tree used by `brrr init`, rooted so that paths are
// relative to the project root (e.g. "go.mod.tmpl", "web/package.json.tmpl").
func Init() (fs.FS, error) {
	return fs.Sub(initFS, "init")
}

// Generate returns a subtree of the `brrr generate` templates rooted at the
// given relative path (e.g. "backend/leaf", "frontend/leaf").
func Generate(sub string) (fs.FS, error) {
	return fs.Sub(generateFS, "generate/"+sub)
}

// GenerateFile reads a single template file from the generate tree (path
// relative to "generate/", e.g. "backend/module.go.tmpl").
func GenerateFile(rel string) ([]byte, error) {
	return generateFS.ReadFile("generate/" + rel)
}
