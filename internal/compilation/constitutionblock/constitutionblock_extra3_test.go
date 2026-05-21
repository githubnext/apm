package constitutionblock

import (
	"strings"
	"testing"
)

func TestComputeConstitutionHash_NotEmpty(t *testing.T) {
	h := ComputeConstitutionHash("some content")
	if h == "" {
		t.Fatal("hash should not be empty")
	}
}

func TestRenderBlock_ContainsContentInside(t *testing.T) {
	block := RenderBlock("my constitution")
	if !strings.Contains(block, "my constitution") {
		t.Fatalf("rendered block should contain the content, got %q", block)
	}
}

func TestRenderBlock_MarkersPresent(t *testing.T) {
	block := RenderBlock("text")
	if !strings.Contains(block, MarkerBegin) {
		t.Fatalf("expected MarkerBegin in block")
	}
	if !strings.Contains(block, MarkerEnd) {
		t.Fatalf("expected MarkerEnd in block")
	}
}

func TestRenderBlock_HashPrefixPresent(t *testing.T) {
	block := RenderBlock("content")
	if !strings.Contains(block, HashPrefix) {
		t.Fatalf("expected HashPrefix in block")
	}
}

func TestInjectOrUpdate_PreservesLeadingContent(t *testing.T) {
	existing := "# Title\n\nsome text\n"
	result, _ := InjectOrUpdate(existing, RenderBlock("constitution"), false)
	if !strings.Contains(result, "# Title") {
		t.Fatalf("expected title preserved, got %q", result)
	}
}

func TestInjectOrUpdate_InjectsBlock(t *testing.T) {
	existing := "some content"
	result, _ := InjectOrUpdate(existing, RenderBlock("constitution body"), false)
	if !strings.Contains(result, MarkerBegin) {
		t.Fatalf("expected constitution block in output")
	}
}

func TestInjectOrUpdate_UpdatesExisting(t *testing.T) {
	existing := "before\n" + RenderBlock("old body") + "\nafter\n"
	result, status := InjectOrUpdate(existing, RenderBlock("new body"), false)
	if status != StatusUpdated && status != StatusUnchanged {
		t.Fatalf("expected updated or unchanged, got %v", status)
	}
	if !strings.Contains(result, "new body") {
		t.Fatalf("expected new body in result")
	}
}

func TestInjectOrUpdate_StatusCreatedOnNew(t *testing.T) {
	_, status := InjectOrUpdate("no block here", RenderBlock("my constitution"), false)
	if status != StatusCreated {
		t.Fatalf("expected StatusCreated, got %v", status)
	}
}

func TestFindExistingBlock_FoundWhenPresent(t *testing.T) {
	block := RenderBlock("content")
	doc := "before\n" + block + "\nafter\n"
	eb := FindExistingBlock(doc)
	if eb == nil {
		t.Fatal("expected block to be found")
	}
	if eb.StartIndex < 0 || eb.EndIndex <= eb.StartIndex {
		t.Fatalf("invalid start/end: %d, %d", eb.StartIndex, eb.EndIndex)
	}
}

func TestFindExistingBlock_NotFoundEmpty(t *testing.T) {
	eb := FindExistingBlock("no markers here")
	if eb != nil {
		t.Fatal("expected block not found")
	}
}

func TestInjectionStatusValues_Distinct(t *testing.T) {
	if StatusCreated == StatusUpdated || StatusUpdated == StatusUnchanged || StatusCreated == StatusUnchanged {
		t.Fatal("injection status values should be distinct")
	}
}

func TestConstitutionRelPath_NonEmpty(t *testing.T) {
	if ConstitutionRelPath == "" {
		t.Fatal("ConstitutionRelPath should not be empty")
	}
}
