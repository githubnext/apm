package postdepslocal

import (
	"reflect"
	"testing"
)

func TestHasLocalContentErrors_ExactlyEqual(t *testing.T) {
	s := LocalContentState{
		LocalContentErrorsBefore: 3,
		CurrentErrorCount:        3,
	}
	if HasLocalContentErrors(s) {
		t.Error("equal counts should not be an error")
	}
}

func TestDetectStaleLocalFiles_NoStale(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles:       []string{"a.txt", "b.txt"},
		OldLocalDeployed:         []string{"a.txt", "b.txt"},
		LocalContentErrorsBefore: 0,
		CurrentErrorCount:        0,
	}
	result := DetectStaleLocalFiles(s)
	if len(result) != 0 {
		t.Errorf("expected no stale files, got %v", result)
	}
}

func TestDetectStaleLocalFiles_NewFileAdded(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles:       []string{"a.txt", "b.txt", "c.txt"},
		OldLocalDeployed:         []string{"a.txt"},
		LocalContentErrorsBefore: 0,
		CurrentErrorCount:        0,
	}
	result := DetectStaleLocalFiles(s)
	if len(result) != 0 {
		t.Errorf("expected no stale (old is subset of new), got %v", result)
	}
}

func TestSortedLocalDeployedFiles_DoesNotMutateInput(t *testing.T) {
	original := []string{"c.txt", "a.txt", "b.txt"}
	cp := make([]string, len(original))
	copy(cp, original)
	SortedLocalDeployedFiles(original)
	if !reflect.DeepEqual(original, cp) {
		t.Error("input slice should not be mutated")
	}
}

func TestSortedLocalDeployedFiles_SingleElement(t *testing.T) {
	result := SortedLocalDeployedFiles([]string{"only.txt"})
	if len(result) != 1 || result[0] != "only.txt" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestShouldRun_ProjectScopeNoContent(t *testing.T) {
	if ShouldRun(true, false, false) {
		t.Error("should not run with no content")
	}
}

func TestShouldRun_AllTrue(t *testing.T) {
	if !ShouldRun(true, true, true) {
		t.Error("should run when project scope with content")
	}
}

func TestDetectStaleLocalFiles_EmptyBothSides(t *testing.T) {
	s := LocalContentState{
		LocalDeployedFiles: nil,
		OldLocalDeployed:   nil,
	}
	result := DetectStaleLocalFiles(s)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
