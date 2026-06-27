package engine

import (
	"fmt"
	"os"
	"strings"
)

// Inject inserts snippet immediately before the first line in the file at path
// that contains marker, matching that line's leading indentation. It is
// idempotent: if the trimmed snippet already appears in the file, nothing is
// written. It returns whether the file was modified.
//
// Marker comments (e.g. "// brrr:routes") are left in place so further
// generation can keep injecting at the same point.
func Inject(path, marker, snippet string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(data)

	if !strings.Contains(content, marker) {
		return false, fmt.Errorf("marker %q not found in %s", marker, path)
	}
	if strings.Contains(content, strings.TrimSpace(snippet)) {
		return false, nil // already injected
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if !strings.Contains(line, marker) {
			continue
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		block := indentBlock(snippet, indent)
		out := make([]string, 0, len(lines)+1)
		out = append(out, lines[:i]...)
		out = append(out, block)
		out = append(out, lines[i:]...)
		return true, os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
	}
	return false, nil
}

// indentBlock prefixes every line of s with indent.
func indentBlock(s, indent string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		lines[i] = indent + l
	}
	return strings.Join(lines, "\n")
}
