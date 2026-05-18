package postdepslocal_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/postdepslocal"
)

func TestHasLocalContentErrors(t *testing.T) {
	tests := []struct {
		name     string
		state    postdepslocal.LocalContentState
		expected bool
	}{
		{
			name:     "no errors when counts equal",
			state:    postdepslocal.LocalContentState{LocalContentErrorsBefore: 2, CurrentErrorCount: 2},
			expected: false,
		},
		{
			name:     "errors when current exceeds before",
			state:    postdepslocal.LocalContentState{LocalContentErrorsBefore: 1, CurrentErrorCount: 3},
			expected: true,
		},
		{
			name:     "no errors at zero",
			state:    postdepslocal.LocalContentState{LocalContentErrorsBefore: 0, CurrentErrorCount: 0},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postdepslocal.HasLocalContentErrors(tt.state)
			if got != tt.expected {
				t.Errorf("HasLocalContentErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectStaleLocalFiles(t *testing.T) {
	tests := []struct {
		name     string
		state    postdepslocal.LocalContentState
		expected []string
	}{
		{
			name: "stale files detected",
			state: postdepslocal.LocalContentState{
				LocalDeployedFiles:       []string{"a.txt", "b.txt"},
				OldLocalDeployed:         []string{"a.txt", "b.txt", "c.txt"},
				LocalContentErrorsBefore: 0,
				CurrentErrorCount:        0,
			},
			expected: []string{"c.txt"},
		},
		{
			name: "no stale files when all still present",
			state: postdepslocal.LocalContentState{
				LocalDeployedFiles:       []string{"a.txt", "b.txt"},
				OldLocalDeployed:         []string{"a.txt"},
				LocalContentErrorsBefore: 0,
				CurrentErrorCount:        0,
			},
			expected: nil,
		},
		{
			name: "nil returned on errors",
			state: postdepslocal.LocalContentState{
				LocalDeployedFiles:       []string{"a.txt"},
				OldLocalDeployed:         []string{"b.txt"},
				LocalContentErrorsBefore: 0,
				CurrentErrorCount:        2,
			},
			expected: nil,
		},
		{
			name: "nil returned when old list is empty",
			state: postdepslocal.LocalContentState{
				LocalDeployedFiles:       []string{"a.txt"},
				OldLocalDeployed:         nil,
				LocalContentErrorsBefore: 0,
				CurrentErrorCount:        0,
			},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postdepslocal.DetectStaleLocalFiles(tt.state)
			if len(got) != len(tt.expected) {
				t.Errorf("DetectStaleLocalFiles() = %v, want %v", got, tt.expected)
				return
			}
			if tt.expected != nil {
				for i, v := range tt.expected {
					if got[i] != v {
						t.Errorf("DetectStaleLocalFiles()[%d] = %q, want %q", i, got[i], v)
					}
				}
			}
		})
	}
}

func TestSortedLocalDeployedFiles(t *testing.T) {
	input := []string{"c.txt", "a.txt", "b.txt"}
	got := postdepslocal.SortedLocalDeployedFiles(input)
	expected := []string{"a.txt", "b.txt", "c.txt"}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("SortedLocalDeployedFiles()[%d] = %q, want %q", i, got[i], v)
		}
	}
	// Ensure original is not mutated
	if input[0] != "c.txt" {
		t.Error("original slice was mutated")
	}
}

func TestShouldRun(t *testing.T) {
	tests := []struct {
		name             string
		isProjectScope   bool
		hasLocalContent  bool
		hasOldLocalContent bool
		expected         bool
	}{
		{"project scope with local content", true, true, false, true},
		{"project scope with old local content", true, false, true, true},
		{"project scope with both", true, true, true, true},
		{"project scope no content", true, false, false, false},
		{"non-project scope", false, true, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postdepslocal.ShouldRun(tt.isProjectScope, tt.hasLocalContent, tt.hasOldLocalContent)
			if got != tt.expected {
				t.Errorf("ShouldRun() = %v, want %v", got, tt.expected)
			}
		})
	}
}
