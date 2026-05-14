// Package compilationconst defines shared constants for compilation extensions.
package compilationconst

// ConstitutionMarkerBegin marks the start of a constitution injection block.
const ConstitutionMarkerBegin = "<!-- SPEC-KIT CONSTITUTION: BEGIN -->"

// ConstitutionMarkerEnd marks the end of a constitution injection block.
const ConstitutionMarkerEnd = "<!-- SPEC-KIT CONSTITUTION: END -->"

// ConstitutionRelativePath is the repo-root-relative path to constitution.md.
const ConstitutionRelativePath = ".specify/memory/constitution.md"

// BuildIDPlaceholder is the sentinel line inserted by formatters before stabilization.
const BuildIDPlaceholder = "<!-- Build ID: __BUILD_ID__ -->"
