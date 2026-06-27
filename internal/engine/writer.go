package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writer is responsible for materialising rendered files onto disk. It tracks
// what it has written and how conflicts with existing files should be handled.
type Writer struct {
	Root    string // destination root directory
	DryRun  bool   // when true, nothing is written; actions are only recorded
	Force   bool   // when true, existing files are overwritten
	Written []string
	Skipped []string
}

// NewWriter constructs a Writer rooted at dst.
func NewWriter(dst string) *Writer { return &Writer{Root: dst} }

// Write places content at relPath (relative to Root), creating parent
// directories as needed. Existing files are skipped unless Force is set.
func (w *Writer) Write(relPath string, content []byte, mode os.FileMode) error {
	dst := filepath.Join(w.Root, relPath)
	if _, err := os.Stat(dst); err == nil && !w.Force {
		w.Skipped = append(w.Skipped, relPath)
		return nil
	}
	w.Written = append(w.Written, relPath)
	if w.DryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}
	if mode == 0 {
		mode = 0o644
	}
	if err := os.WriteFile(dst, content, mode); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	return nil
}

// Summary returns a short human readable report of the writer's activity.
func (w *Writer) Summary() string {
	return fmt.Sprintf("%d files written, %d skipped", len(w.Written), len(w.Skipped))
}
