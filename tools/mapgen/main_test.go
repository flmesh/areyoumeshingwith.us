package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// Characterization: run() regenerates the committed SVG byte-identically.
// Restores the original file via t.Cleanup so a failed assertion cannot leave
// the working tree dirty.
func TestRun_RegeneratesCommittedSVG(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	outPath := filepath.Join(repoRoot, outputFile)

	original, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read committed SVG: %v", err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(outPath, original, 0644); err != nil {
			t.Errorf("restore committed SVG: %v", err)
		}
	})

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read regenerated SVG: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("regenerated SVG differs from committed (%d bytes vs %d bytes)", len(got), len(original))
	}
}
