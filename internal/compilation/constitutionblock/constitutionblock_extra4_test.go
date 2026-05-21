package constitutionblock

import (
"strings"
"testing"
)

func TestMarkerBegin_ContainsBegin_Extra4(t *testing.T) {
if !strings.Contains(MarkerBegin, "BEGIN") {
t.Error("MarkerBegin should contain 'BEGIN'")
}
}

func TestMarkerEnd_ContainsEnd_Extra4(t *testing.T) {
if !strings.Contains(MarkerEnd, "END") {
t.Error("MarkerEnd should contain 'END'")
}
}

func TestHashPrefix_Value_Extra4(t *testing.T) {
if HashPrefix == "" {
t.Error("HashPrefix should not be empty")
}
}

func TestConstitutionRelPath_NonEmpty_Extra4(t *testing.T) {
if ConstitutionRelPath == "" {
t.Error("ConstitutionRelPath should not be empty")
}
}

func TestComputeConstitutionHash_Length_Extra4(t *testing.T) {
h := ComputeConstitutionHash("some content")
if len(h) != 12 {
t.Errorf("expected 12-char hash, got %d", len(h))
}
}

func TestComputeConstitutionHash_Hex_Extra4(t *testing.T) {
h := ComputeConstitutionHash("test")
for _, c := range h {
if !strings.ContainsRune("0123456789abcdef", c) {
t.Errorf("hash contains non-hex char: %c", c)
}
}
}

func TestRenderBlock_HasMarkers_Extra4(t *testing.T) {
block := RenderBlock("constitution text")
if !strings.Contains(block, MarkerBegin) {
t.Error("rendered block should contain MarkerBegin")
}
if !strings.Contains(block, MarkerEnd) {
t.Error("rendered block should contain MarkerEnd")
}
}

func TestRenderBlock_HasContent_Extra4(t *testing.T) {
block := RenderBlock("my constitution content")
if !strings.Contains(block, "my constitution content") {
t.Error("rendered block should contain the constitution content")
}
}

func TestInjectOrUpdate_CreatesNew_Extra4(t *testing.T) {
block := RenderBlock("content")
result, status := InjectOrUpdate("# AGENTS", block, false)
if status == "" {
t.Error("expected non-empty status")
}
if !strings.Contains(result, MarkerBegin) {
t.Error("result should contain MarkerBegin")
}
}

func TestInjectOrUpdate_Unchanged_Extra4(t *testing.T) {
block := RenderBlock("stable content")
// inject once
result1, _ := InjectOrUpdate("# AGENTS", block, false)
// inject again with same block
_, status2 := InjectOrUpdate(result1, block, false)
if status2 != StatusUnchanged {
t.Errorf("expected UNCHANGED, got %q", status2)
}
}
