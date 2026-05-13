// Package discovery finds workflow definition files.
package discovery

import (
"os"
"path/filepath"
"strings"

"github.com/githubnext/apm/internal/workflow/wfparser"
)

// DiscoverWorkflows finds all .prompt.md files under baseDir.
func DiscoverWorkflows(baseDir string) ([]*wfparser.WorkflowDefinition, []error) {
if baseDir == "" {
var err error
baseDir, err = os.Getwd()
if err != nil {
return nil, []error{err}
}
}

var files []string
_ = filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
if err != nil {
return nil
}
if !d.IsDir() && strings.HasSuffix(path, ".prompt.md") {
files = append(files, path)
}
return nil
})

// Deduplicate
seen := map[string]bool{}
var workflows []*wfparser.WorkflowDefinition
var errs []error
for _, f := range files {
if seen[f] {
continue
}
seen[f] = true
w, err := wfparser.ParseWorkflowFile(f)
if err != nil {
errs = append(errs, err)
continue
}
workflows = append(workflows, w)
}
return workflows, errs
}
