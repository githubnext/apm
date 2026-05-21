package postdepslocal

import (
	"reflect"
	"testing"
)

func TestHasLocalContentErrors_NoNewErrors(t *testing.T) {
	s := LocalContentState{LocalContentErrorsBefore: 5, CurrentErrorCount: 5}
	if HasLocalContentErrors(s) {
		t.Error("expected no new errors when counts are equal")
	}
}

func TestHasLocalContentErrors_FewerErrors(t *testing.T) {
	// CurrentErrorCount < before is not expected in practice, but should not panic.
	s := LocalContentState{LocalContentErrorsBefore: 3, CurrentErrorCount: 1}
	if HasLocalContentErrors(s) {
		t.Error("expected false when current <= before")
	}
}

func TestHasLocalContentErrors_NewErrors(t *testing.T) {
	s := LocalContentState{LocalContentErrorsBefore: 2, CurrentErrorCount: 4}
	if !HasLocalContentErrors(s) {
		t.Error("expected true when current > before")
	}
}

func TestDetectStaleLocalFiles_WithErrors(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles:       []string{"a.md"},
		OldLocalDeployed:         []string{"b.md"},
		LocalContentErrorsBefore: 0,
		CurrentErrorCount:        1,
	}
	got := DetectStaleLocalFiles(s)
	if len(got) != 0 {
		t.Error("expected nil/empty when errors occurred")
	}
}

func TestDetectStaleLocalFiles_EmptyOld(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles: []string{"a.md"},
		OldLocalDeployed:   nil,
	}
	got := DetectStaleLocalFiles(s)
	if len(got) != 0 {
		t.Error("expected nil/empty when OldLocalDeployed is empty")
	}
}

func TestDetectStaleLocalFiles_AllStale(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles: []string{},
		OldLocalDeployed:   []string{"old1.md", "old2.md"},
	}
	got := DetectStaleLocalFiles(s)
	if len(got) != 2 {
		t.Errorf("expected 2 stale files, got %d", len(got))
	}
}

func TestDetectStaleLocalFiles_PartialStale(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles: []string{"a.md", "b.md"},
		OldLocalDeployed:   []string{"a.md", "b.md", "c.md", "d.md"},
	}
	got := DetectStaleLocalFiles(s)
	if len(got) != 2 {
		t.Errorf("expected 2 stale files, got %d: %v", len(got), got)
	}
	for _, f := range got {
		if f != "c.md" && f != "d.md" {
			t.Errorf("unexpected stale file %q", f)
		}
	}
}

func TestSortedLocalDeployedFiles_Order(t *testing.T) {
	input := []string{"z.md", "a.md", "m.md"}
	got := SortedLocalDeployedFiles(input)
	expected := []string{"a.md", "m.md", "z.md"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("SortedLocalDeployedFiles: got %v, want %v", got, expected)
	}
	// Original must not be mutated.
	if input[0] != "z.md" {
		t.Error("original slice was mutated")
	}
}

func TestSortedLocalDeployedFiles_Empty(t *testing.T) {
	got := SortedLocalDeployedFiles(nil)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestShouldRun_NotProjectScope(t *testing.T) {
	if ShouldRun(false, true, true) {
		t.Error("should not run when not project scope")
	}
}

func TestShouldRun_NoContent(t *testing.T) {
	if ShouldRun(true, false, false) {
		t.Error("should not run when neither has local content")
	}
}

func TestShouldRun_HasLocalContent(t *testing.T) {
	if !ShouldRun(true, true, false) {
		t.Error("should run when hasLocalContent is true")
	}
}

func TestShouldRun_HasOldLocalContent(t *testing.T) {
	if !ShouldRun(true, false, true) {
		t.Error("should run when hasOldLocalContent is true")
	}
}
