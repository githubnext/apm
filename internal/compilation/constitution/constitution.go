// Package constitution reads Spec Kit constitution files.
package constitution

import (
"os"
"path/filepath"
"sync"

"github.com/githubnext/apm/internal/compilation/compilationconst"
)

var (
mu    sync.Mutex
cache = map[string]*string{}
)

// ClearCache clears the constitution read cache.
func ClearCache() {
mu.Lock()
defer mu.Unlock()
cache = map[string]*string{}
}

// FindConstitution returns the path to constitution.md relative to baseDir.
func FindConstitution(baseDir string) string {
return filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
}

// ReadConstitution reads the full constitution content if the file exists.
// Results are cached by resolved baseDir for the lifetime of the process.
func ReadConstitution(baseDir string) (string, bool) {
resolved, err := filepath.Abs(baseDir)
if err != nil {
resolved = baseDir
}
mu.Lock()
if v, ok := cache[resolved]; ok {
mu.Unlock()
if v == nil {
return "", false
}
return *v, true
}
mu.Unlock()

path := FindConstitution(resolved)
data, err := os.ReadFile(path)
mu.Lock()
defer mu.Unlock()
if err != nil {
cache[resolved] = nil
return "", false
}
s := string(data)
cache[resolved] = &s
return s, true
}
