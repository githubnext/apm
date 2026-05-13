// Package outputwriter provides a single chokepoint for persisting compiled outputs.
package outputwriter

import (
"fmt"
"os"
"path/filepath"
"strings"

"github.com/githubnext/apm/internal/compilation/buildid"
"github.com/githubnext/apm/internal/compilation/compilationconst"
)

// CompiledOutputWriter persists compiled output with cross-cutting concerns applied.
type CompiledOutputWriter struct{}

// Write stabilizes the build ID, validates no placeholder remains, and writes atomically.
func (w *CompiledOutputWriter) Write(path, content string) error {
final := buildid.StabilizeBuildID(content)
if strings.Contains(final, compilationconst.BuildIDPlaceholder) {
return fmt.Errorf("build_id stabilization bypassed: placeholder still present after stabilization (target=%s)", path)
}
if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
return err
}
return atomicWrite(path, final)
}

func atomicWrite(path, content string) error {
dir := filepath.Dir(path)
tmp, err := os.CreateTemp(dir, ".apm-write-*")
if err != nil {
return err
}
tmpName := tmp.Name()
if _, err := tmp.WriteString(content); err != nil {
tmp.Close()
os.Remove(tmpName)
return err
}
if err := tmp.Close(); err != nil {
os.Remove(tmpName)
return err
}
return os.Rename(tmpName, path)
}
