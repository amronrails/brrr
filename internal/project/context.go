package project

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// Context carries the values that parameterise a generated project. It is the
// data object handed to every init template.
type Context struct {
	AppName    string // application name, e.g. "acme"
	ModulePath string // Go module path, e.g. "github.com/acme/acme"
	GoVersion  string // Go version for go.mod, e.g. "1.26"
	HTTPPort   int    // default HTTP port

	DBName     string // Postgres database name
	DBUser     string // Postgres user
	DBPassword string // Postgres password (local/dev only)
	DBPort     int    // Postgres port exposed by docker-compose
}

var appNameRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// NewContext builds a Context from an app name and module path, filling in
// sensible local-development defaults for anything not supplied.
func NewContext(appName, modulePath string) (*Context, error) {
	appName = strings.TrimSpace(appName)
	if !appNameRe.MatchString(appName) {
		return nil, fmt.Errorf("invalid app name %q: must start with a letter and contain only lowercase letters, digits, '-' or '_'", appName)
	}
	if modulePath == "" {
		modulePath = appName
	}
	if !validModulePath(modulePath) {
		return nil, fmt.Errorf("invalid module path %q", modulePath)
	}
	db := strings.ReplaceAll(appName, "-", "_")
	return &Context{
		AppName:    appName,
		ModulePath: modulePath,
		GoVersion:  "1.26",
		HTTPPort:   8080,
		DBName:     db,
		DBUser:     db,
		DBPassword: db + "_dev",
		DBPort:     5432,
	}, nil
}

// validModulePath performs a light sanity check on a Go module path. It is not
// a full validator; it rejects obviously malformed input.
func validModulePath(p string) bool {
	if p == "" || strings.ContainsAny(p, " \t\n") {
		return false
	}
	for _, seg := range strings.Split(path.Clean(p), "/") {
		if seg == "" || seg == "." || seg == ".." {
			return false
		}
	}
	return true
}
