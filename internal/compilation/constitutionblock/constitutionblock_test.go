package constitutionblock_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/constitutionblock"
)

func TestComputeConstitutionHash(t *testing.T) {
	h1 := constitutionblock.ComputeConstitutionHash("hello world")
	h2 := constitutionblock.ComputeConstitutionHash("hello world")
	if h1 != h2 {
		t.Error("same input should produce same hash")
	}
	if len(h1) != 12 {
		t.Errorf("hash length = %d, want 12", len(h1))
	}
	h3 := constitutionblock.ComputeConstitutionHash("different content")
	if h1 == h3 {
		t.Error("different inputs should produce different hashes")
	}
}

func TestRenderBlock(t *testing.T) {
	content := "# My Constitution\n\nSome rules here.\n"
	block := constitutionblock.RenderBlock(content)

	if !strings.Contains(block, constitutionblock.MarkerBegin) {
		t.Error("RenderBlock should contain MarkerBegin")
	}
	if !strings.Contains(block, constitutionblock.MarkerEnd) {
		t.Error("RenderBlock should contain MarkerEnd")
	}
	if !strings.Contains(block, constitutionblock.HashPrefix) {
		t.Error("RenderBlock should contain hash prefix")
	}
	if !strings.Contains(block, "Some rules here.") {
		t.Error("RenderBlock should contain constitution content")
	}
}

func TestFindExistingBlockNone(t *testing.T) {
	result := constitutionblock.FindExistingBlock("# Just some markdown\n\nNo constitution here.\n")
	if result != nil {
		t.Error("FindExistingBlock should return nil when no block exists")
	}
}

func TestFindExistingBlockFound(t *testing.T) {
	content := "# My Constitution\n\nSome rules here.\n"
	block := constitutionblock.RenderBlock(content)
	agents := "# AGENTS.md\n\n" + block + "\n## Other section\n"

	result := constitutionblock.FindExistingBlock(agents)
	if result == nil {
		t.Fatal("FindExistingBlock should find the block")
	}
	if result.Hash == "" {
		t.Error("FindExistingBlock should extract hash")
	}
	if result.StartIndex < 0 {
		t.Error("StartIndex should be non-negative")
	}
	if result.EndIndex <= result.StartIndex {
		t.Error("EndIndex should be greater than StartIndex")
	}
}

func TestInjectOrUpdateCreatesNew(t *testing.T) {
	content := "# My Constitution\n"
	block := constitutionblock.RenderBlock(content)
	existing := "# AGENTS.md\n\nSome content.\n"

	result, status := constitutionblock.InjectOrUpdate(existing, block, false)
	if status != constitutionblock.StatusCreated {
		t.Errorf("expected StatusCreated, got %q", status)
	}
	if !strings.Contains(result, constitutionblock.MarkerBegin) {
		t.Error("result should contain constitution block")
	}
}

func TestInjectOrUpdateCreatesNewAtTop(t *testing.T) {
	content := "# My Constitution\n"
	block := constitutionblock.RenderBlock(content)
	existing := "# AGENTS.md\n\nSome content.\n"

	result, status := constitutionblock.InjectOrUpdate(existing, block, true)
	if status != constitutionblock.StatusCreated {
		t.Errorf("expected StatusCreated, got %q", status)
	}
	if !strings.HasPrefix(result, constitutionblock.MarkerBegin) {
		t.Error("result should start with constitution block when placeTop=true")
	}
}

func TestInjectOrUpdateUnchanged(t *testing.T) {
	content := "# My Constitution\n"
	block := constitutionblock.RenderBlock(content)
	existing := block + "\n# AGENTS.md\n"

	_, status := constitutionblock.InjectOrUpdate(existing, block, true)
	if status != constitutionblock.StatusUnchanged {
		t.Errorf("expected StatusUnchanged for identical block, got %q", status)
	}
}

func TestInjectOrUpdateUpdates(t *testing.T) {
	oldContent := "# Old Constitution\n"
	oldBlock := constitutionblock.RenderBlock(oldContent)
	existing := "# AGENTS.md\n\n" + oldBlock

	newContent := "# New Constitution\n\nDifferent rules.\n"
	newBlock := constitutionblock.RenderBlock(newContent)

	result, status := constitutionblock.InjectOrUpdate(existing, newBlock, false)
	if status != constitutionblock.StatusUpdated {
		t.Errorf("expected StatusUpdated, got %q", status)
	}
	if strings.Contains(result, "Old Constitution") {
		t.Error("result should not contain old constitution content")
	}
	if !strings.Contains(result, "New Constitution") {
		t.Error("result should contain new constitution content")
	}
}
