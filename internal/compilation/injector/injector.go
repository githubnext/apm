// Package injector implements the constitution injection workflow for compile command.
package injector

import (
"os"
"strings"

"github.com/githubnext/apm/internal/compilation/compilationconst"
)

// InjectionStatus represents the outcome of a constitution injection attempt.
type InjectionStatus string

const (
StatusCreated   InjectionStatus = "CREATED"
StatusUpdated   InjectionStatus = "UPDATED"
StatusUnchanged InjectionStatus = "UNCHANGED"
StatusSkipped   InjectionStatus = "SKIPPED"
StatusMissing   InjectionStatus = "MISSING"
)

// ConstitutionInjector encapsulates constitution detection and injection logic.
type ConstitutionInjector struct {
BaseDir string
}

// Inject returns final AGENTS.md content after optional constitution injection.
// Returns (finalContent, status, hashOrEmpty).
func (ci *ConstitutionInjector) Inject(compiledContent string, withConstitution bool, outputPath string) (string, InjectionStatus, string) {
existingContent := ""
if data, err := os.ReadFile(outputPath); err == nil {
existingContent = string(data)
}

if !withConstitution {
// Preserve any existing constitution block.
block := extractConstitutionBlock(existingContent)
if block == "" {
return compiledContent, StatusSkipped, ""
}
return injectBlock(compiledContent, block), StatusUnchanged, ""
}

// Read constitution file.
constitPath := ci.BaseDir + "/" + compilationconst.ConstitutionRelativePath
constitData, err := os.ReadFile(constitPath)
if err != nil {
return compiledContent, StatusMissing, ""
}
block := compilationconst.ConstitutionMarkerBegin + "\n" + string(constitData) + "\n" + compilationconst.ConstitutionMarkerEnd

existing := extractConstitutionBlock(existingContent)
status := StatusCreated
if existing != "" {
if existing == block {
status = StatusUnchanged
} else {
status = StatusUpdated
}
}
return injectBlock(compiledContent, block), status, ""
}

func extractConstitutionBlock(content string) string {
begin := strings.Index(content, compilationconst.ConstitutionMarkerBegin)
if begin < 0 {
return ""
}
end := strings.Index(content[begin:], compilationconst.ConstitutionMarkerEnd)
if end < 0 {
return ""
}
return content[begin : begin+end+len(compilationconst.ConstitutionMarkerEnd)]
}

func injectBlock(content, block string) string {
// Remove existing block if present
if idx := strings.Index(content, compilationconst.ConstitutionMarkerBegin); idx >= 0 {
endIdx := strings.Index(content[idx:], compilationconst.ConstitutionMarkerEnd)
if endIdx >= 0 {
after := content[idx+endIdx+len(compilationconst.ConstitutionMarkerEnd):]
content = content[:idx] + after
}
}
// Prepend block
return block + "\n\n" + content
}
