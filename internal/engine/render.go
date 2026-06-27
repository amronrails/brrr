package engine

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

// tmplSuffix is the extension every template file on disk must carry. It keeps
// template ".go" sources from being compiled into the brrr binary and makes the
// mapping from template path to output path unambiguous.
const tmplSuffix = ".tmpl"

// Render walks the template tree rooted at src (an fs.FS, typically an embedded
// sub-filesystem), renders every file's path and contents with data, and writes
// the results through w. Both the relative path and the file body are treated as
// templates so callers can drive directory and file names from data.
//
// A template whose rendered path is empty (or whose path renders to only
// whitespace) is skipped, which lets templates opt out of generation
// conditionally.
func Render(src fs.FS, data any, w *Writer) error {
	funcs := FuncMap()
	return fs.WalkDir(src, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		raw, err := fs.ReadFile(src, p)
		if err != nil {
			return fmt.Errorf("read template %s: %w", p, err)
		}

		outPath, err := renderString("path:"+p, strings.TrimSuffix(p, tmplSuffix), data, funcs)
		if err != nil {
			return err
		}
		outPath = strings.TrimSpace(outPath)
		if outPath == "" {
			return nil // template opted out
		}

		body, err := renderBytes("body:"+p, raw, data, funcs)
		if err != nil {
			return err
		}

		return w.Write(outPath, body, fileMode(outPath))
	})
}

// RenderBytes renders a single named template body with the engine's standard
// FuncMap and strict missing-key handling. It is used for one-off files whose
// output path is computed by the caller (e.g. regenerated aggregator files).
func RenderBytes(name string, content []byte, data any) ([]byte, error) {
	return renderBytes(name, content, data, FuncMap())
}

func renderString(name, text string, data any, funcs template.FuncMap) (string, error) {
	out, err := renderBytes(name, []byte(text), data, funcs)
	return string(out), err
}

func renderBytes(name string, text []byte, data any, funcs template.FuncMap) ([]byte, error) {
	t, err := template.New(name).Funcs(funcs).Option("missingkey=error").Parse(string(text))
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

// fileMode picks a sensible permission for an output path: scripts and a few
// well known executable files get the execute bit, everything else is 0644.
func fileMode(p string) os.FileMode {
	base := path.Base(p)
	if strings.HasSuffix(base, ".sh") || base == "gradlew" {
		return 0o755
	}
	return 0o644
}
