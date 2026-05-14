// Package constitutionblock provides rendering and parsing of the injected
// constitution block in AGENTS.md.
// Mirrors src/apm_cli/compilation/constitution_block.py.
package constitutionblock

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

// Constants used for the constitution block markers (imported from compilationconst).
const (
	MarkerBegin        = "<!-- BEGIN: APM CONSTITUTION -->"
	MarkerEnd          = "<!-- END: APM CONSTITUTION -->"
	ConstitutionRelPath = ".apm/constitution.md"
	HashPrefix         = "hash:"
)

// ComputeConstitutionHash returns a 12-character hex SHA-256 of the constitution content.
func ComputeConstitutionHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum)[:12]
}

// RenderBlock renders the full constitution block with markers and hash line.
func RenderBlock(constitutionContent string) string {
	h := ComputeConstitutionHash(constitutionContent)
	headerMeta := fmt.Sprintf("%s %s path: %s", HashPrefix, h, ConstitutionRelPath)
	body := strings.TrimRight(constitutionContent, "\n") + "\n"
	return fmt.Sprintf("%s\n%s\n%s%s\n\n", MarkerBegin, headerMeta, body, MarkerEnd)
}

// ExistingBlock represents a constitution block found in an AGENTS.md file.
type ExistingBlock struct {
	Raw        string
	Hash       string // may be empty if no hash line found
	StartIndex int
	EndIndex   int
}

var (
	blockRegex    = regexp.MustCompile(`(?s)(` + regexp.QuoteMeta(MarkerBegin) + `)(.*?)(` + regexp.QuoteMeta(MarkerEnd) + `)`)
	hashLineRegex = regexp.MustCompile(`hash:\s*([0-9a-fA-F]{6,64})`)
)

// FindExistingBlock locates an existing constitution block and extracts its hash.
// Returns nil if no block is found.
func FindExistingBlock(content string) *ExistingBlock {
	loc := blockRegex.FindStringIndex(content)
	if loc == nil {
		return nil
	}
	blockText := content[loc[0]:loc[1]]
	h := ""
	if hm := hashLineRegex.FindStringSubmatch(blockText); hm != nil {
		h = hm[1]
	}
	return &ExistingBlock{
		Raw:        blockText,
		Hash:       h,
		StartIndex: loc[0],
		EndIndex:   loc[1],
	}
}

// InjectionStatus represents the outcome of InjectOrUpdate.
type InjectionStatus string

const (
	StatusCreated   InjectionStatus = "CREATED"
	StatusUpdated   InjectionStatus = "UPDATED"
	StatusUnchanged InjectionStatus = "UNCHANGED"
)

// InjectOrUpdate inserts or updates the constitution block in existing AGENTS.md content.
// placeTop=true always prepends at the top (Phase 0 behaviour).
// Returns (updatedText, status).
func InjectOrUpdate(existingAgents, newBlock string, placeTop bool) (string, InjectionStatus) {
	existing := FindExistingBlock(existingAgents)
	if existing != nil {
		if existing.Raw == strings.TrimRight(newBlock, "\n") {
			return existingAgents, StatusUnchanged
		}
		updated := existingAgents[:existing.StartIndex] +
			strings.TrimRight(newBlock, "\n") +
			existingAgents[existing.EndIndex:]
		if placeTop && !strings.HasPrefix(updated, newBlock) {
			bodyWithoutBlock := strings.TrimLeft(strings.Replace(updated, strings.TrimRight(newBlock, "\n"), "", 1), "\n")
			updated = newBlock + bodyWithoutBlock
		}
		return updated, StatusUpdated
	}
	// No existing block.
	if placeTop {
		return newBlock + strings.TrimLeft(existingAgents, "\n"), StatusCreated
	}
	sep := ""
	if len(existingAgents) > 0 && !strings.HasSuffix(existingAgents, "\n") {
		sep = "\n"
	}
	return existingAgents + sep + newBlock, StatusCreated
}
