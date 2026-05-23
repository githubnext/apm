// Package paths_test provides parity tests for path utilities.
package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/paths"
)

// TestParityPathsPortableRelpath verifies portable_relpath behavior.
func TestParityPathsPortableRelpath(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name     string
		path     string
		base     string
		expected string
	}{
		{"nested path", filepath.Join(tmp, "a", "b", "c.txt"), tmp, "a/b/c.txt"},
		{"direct child", filepath.Join(tmp, "file.txt"), tmp, "file.txt"},
		{"same dir", tmp, tmp, "."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := paths.PortableRelpath(tc.path, tc.base)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}
