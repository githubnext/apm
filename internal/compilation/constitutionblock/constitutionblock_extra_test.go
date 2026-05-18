package constitutionblock_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/constitutionblock"
)

func TestComputeConstitutionHash_Length(t *testing.T) {
	h := constitutionblock.ComputeConstitutionHash("some content")
	if len(h) != 12 {
		t.Errorf("hash length = %d, want 12", len(h))
	}
}

func TestComputeConstitutionHash_Deterministic(t *testing.T) {
	h1 := constitutionblock.ComputeConstitutionHash("abc")
	h2 := constitutionblock.ComputeConstitutionHash("abc")
	if h1 != h2 {
		t.Error("same input must produce same hash")
	}
}

func TestComputeConstitutionHash_Different(t *testing.T) {
	h1 := constitutionblock.ComputeConstitutionHash("version A")
	h2 := constitutionblock.ComputeConstitutionHash("version B")
	if h1 == h2 {
		t.Error("different content should produce different hashes")
	}
}

func TestComputeConstitutionHash_Empty(t *testing.T) {
	h := constitutionblock.ComputeConstitutionHash("")
	if len(h) != 12 {
		t.Errorf("empty input: hash length = %d, want 12", len(h))
	}
}

func TestRenderBlock_ContainsConstitutionPath(t *testing.T) {
	block := constitutionblock.RenderBlock("# Rules\n")
	if !strings.Contains(block, constitutionblock.ConstitutionRelPath) {
		t.Error("RenderBlock should contain ConstitutionRelPath")
	}
}

func TestRenderBlock_HashInBlock(t *testing.T) {
	content := "# Content\n"
	block := constitutionblock.RenderBlock(content)
	h := constitutionblock.ComputeConstitutionHash(content)
	if !strings.Contains(block, h) {
		t.Errorf("RenderBlock should contain hash %q", h)
	}
}

func TestRenderBlock_StartsWithMarkerBegin(t *testing.T) {
	block := constitutionblock.RenderBlock("content")
	if !strings.HasPrefix(block, constitutionblock.MarkerBegin) {
		t.Errorf("RenderBlock should start with MarkerBegin, got: %q", block[:50])
	}
}

func TestRenderBlock_EndsWithMarkerEnd(t *testing.T) {
	block := constitutionblock.RenderBlock("content")
	// trailing newlines may exist
	trimmed := strings.TrimRight(block, "\n")
	if !strings.HasSuffix(trimmed, constitutionblock.MarkerEnd) {
		t.Errorf("RenderBlock should end with MarkerEnd")
	}
}

func TestFindExistingBlock_HashExtracted(t *testing.T) {
	content := "# Rules\n"
	block := constitutionblock.RenderBlock(content)
	h := constitutionblock.ComputeConstitutionHash(content)
	existing := constitutionblock.FindExistingBlock(block)
	if existing == nil {
		t.Fatal("expected to find block")
	}
	if existing.Hash != h {
		t.Errorf("Hash = %q, want %q", existing.Hash, h)
	}
}

func TestFindExistingBlock_Indices(t *testing.T) {
	prefix := "some prefix\n"
	block := constitutionblock.RenderBlock("content")
	doc := prefix + block
	existing := constitutionblock.FindExistingBlock(doc)
	if existing == nil {
		t.Fatal("expected to find block")
	}
	if existing.StartIndex != len(prefix) {
		t.Errorf("StartIndex = %d, want %d", existing.StartIndex, len(prefix))
	}
}

func TestInjectOrUpdate_StatusCreated_Bottom(t *testing.T) {
	newBlock := constitutionblock.RenderBlock("rules")
	result, status := constitutionblock.InjectOrUpdate("existing content\n", newBlock, false)
	if status != constitutionblock.StatusCreated {
		t.Errorf("status = %q, want %q", status, constitutionblock.StatusCreated)
	}
	if !strings.Contains(result, "existing content") {
		t.Error("original content should be preserved")
	}
}

func TestInjectOrUpdate_StatusCreated_Top(t *testing.T) {
	newBlock := constitutionblock.RenderBlock("rules")
	result, status := constitutionblock.InjectOrUpdate("existing content\n", newBlock, true)
	if status != constitutionblock.StatusCreated {
		t.Errorf("status = %q, want %q", status, constitutionblock.StatusCreated)
	}
	if !strings.HasPrefix(result, constitutionblock.MarkerBegin) {
		t.Error("with placeTop=true, block should be at top")
	}
}

func TestInjectOrUpdate_StatusUpdated(t *testing.T) {
	oldBlock := constitutionblock.RenderBlock("old rules")
	newBlock := constitutionblock.RenderBlock("new rules")
	_, status := constitutionblock.InjectOrUpdate(oldBlock, newBlock, false)
	if status != constitutionblock.StatusUpdated {
		t.Errorf("status = %q, want %q", status, constitutionblock.StatusUpdated)
	}
}

func TestInjectOrUpdate_StatusUnchanged(t *testing.T) {
	block := constitutionblock.RenderBlock("same rules")
	result, status := constitutionblock.InjectOrUpdate(block, block, false)
	if status != constitutionblock.StatusUnchanged {
		t.Errorf("status = %q, want %q", status, constitutionblock.StatusUnchanged)
	}
	if result != block {
		t.Error("unchanged document should not be modified")
	}
}

func TestMarkerConstants_NonEmpty(t *testing.T) {
	if constitutionblock.MarkerBegin == "" {
		t.Error("MarkerBegin must not be empty")
	}
	if constitutionblock.MarkerEnd == "" {
		t.Error("MarkerEnd must not be empty")
	}
	if constitutionblock.HashPrefix == "" {
		t.Error("HashPrefix must not be empty")
	}
}
