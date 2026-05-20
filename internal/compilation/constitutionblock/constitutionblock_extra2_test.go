package constitutionblock_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/constitutionblock"
)

func TestComputeConstitutionHash_Is12Chars(t *testing.T) {
	h := constitutionblock.ComputeConstitutionHash("test content")
	if len(h) != 12 {
		t.Errorf("expected 12-char hash, got %d chars: %q", len(h), h)
	}
}

func TestComputeConstitutionHash_IsHex(t *testing.T) {
	h := constitutionblock.ComputeConstitutionHash("test content")
	for _, c := range h {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("hash should be lowercase hex, got char %q in %q", c, h)
		}
	}
}

func TestComputeConstitutionHash_EmptyInput(t *testing.T) {
	h := constitutionblock.ComputeConstitutionHash("")
	if len(h) != 12 {
		t.Errorf("empty hash should be 12 chars, got %d", len(h))
	}
}

func TestRenderBlock_ContainsMarkers(t *testing.T) {
	block := constitutionblock.RenderBlock("my content")
	if !strings.Contains(block, constitutionblock.MarkerBegin) {
		t.Error("block should contain begin marker")
	}
	if !strings.Contains(block, constitutionblock.MarkerEnd) {
		t.Error("block should contain end marker")
	}
}

func TestRenderBlock_ContainsContent(t *testing.T) {
	block := constitutionblock.RenderBlock("unique-const-content")
	if !strings.Contains(block, "unique-const-content") {
		t.Error("block should contain the constitution content")
	}
}

func TestRenderBlock_ContainsHash(t *testing.T) {
	content := "my constitution text"
	hash := constitutionblock.ComputeConstitutionHash(content)
	block := constitutionblock.RenderBlock(content)
	if !strings.Contains(block, hash) {
		t.Errorf("block should contain hash %q", hash)
	}
}

func TestRenderBlock_ContainsConstitutionRelPath(t *testing.T) {
	block := constitutionblock.RenderBlock("content")
	if !strings.Contains(block, constitutionblock.ConstitutionRelPath) {
		t.Errorf("block should contain path %q", constitutionblock.ConstitutionRelPath)
	}
}

func TestFindExistingBlock_NotFound(t *testing.T) {
	result := constitutionblock.FindExistingBlock("no markers here")
	if result != nil {
		t.Error("expected nil for content without markers")
	}
}

func TestFindExistingBlock_Found(t *testing.T) {
	block := constitutionblock.RenderBlock("some content")
	result := constitutionblock.FindExistingBlock(block)
	if result == nil {
		t.Fatal("expected non-nil result for content with markers")
	}
}

func TestFindExistingBlock_HashExtracted2(t *testing.T) {
	block := constitutionblock.RenderBlock("content to hash")
	result := constitutionblock.FindExistingBlock(block)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	expected := constitutionblock.ComputeConstitutionHash("content to hash")
	if result.Hash != expected {
		t.Errorf("expected hash %q, got %q", expected, result.Hash)
	}
}

func TestInjectOrUpdate_CreateNew(t *testing.T) {
	block := constitutionblock.RenderBlock("init")
	updated, status := constitutionblock.InjectOrUpdate("existing agents\n", block, false)
	if status != constitutionblock.StatusCreated {
		t.Errorf("expected StatusCreated, got %q", status)
	}
	if !strings.Contains(updated, "existing agents") {
		t.Error("updated content should keep original content")
	}
}

func TestInjectOrUpdate_Updated(t *testing.T) {
	original := constitutionblock.RenderBlock("v1")
	doc := "preamble\n" + original
	newBlock := constitutionblock.RenderBlock("v2")
	_, status := constitutionblock.InjectOrUpdate(doc, newBlock, false)
	if status != constitutionblock.StatusUpdated {
		t.Errorf("expected StatusUpdated, got %q", status)
	}
}

func TestInjectOrUpdate_Unchanged(t *testing.T) {
	block := constitutionblock.RenderBlock("same content")
	doc := "header\n" + block
	_, status := constitutionblock.InjectOrUpdate(doc, block, false)
	if status != constitutionblock.StatusUnchanged {
		t.Errorf("expected StatusUnchanged, got %q", status)
	}
}

func TestInjectionStatus_Values(t *testing.T) {
	if constitutionblock.StatusCreated == "" {
		t.Error("StatusCreated should not be empty")
	}
	if constitutionblock.StatusUpdated == "" {
		t.Error("StatusUpdated should not be empty")
	}
	if constitutionblock.StatusUnchanged == "" {
		t.Error("StatusUnchanged should not be empty")
	}
	if constitutionblock.StatusCreated == constitutionblock.StatusUpdated {
		t.Error("StatusCreated and StatusUpdated should differ")
	}
}
