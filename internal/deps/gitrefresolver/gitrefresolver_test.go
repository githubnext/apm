package gitrefresolver

import (
	"testing"
)

func TestIsFullSHA(t *testing.T) {
	tests := []struct {
		ref  string
		want bool
	}{
		{"abcdef1234567890abcdef1234567890abcdef12", true},
		{"0000000000000000000000000000000000000000", true},
		{"abc123", false},
		{"", false},
		{"abcdef1234567890abcdef1234567890abcdef1", false},  // 39 chars
		{"abcdef1234567890abcdef1234567890abcdef123", false}, // 41 chars
		{"ABCDEF1234567890abcdef1234567890abcdef12", false},  // uppercase
		{"main", false},
		{"v1.2.3", false},
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			got := IsFullSHA(tc.ref)
			if got != tc.want {
				t.Errorf("IsFullSHA(%q) = %v, want %v", tc.ref, got, tc.want)
			}
		})
	}
}

func TestIsShortSHA(t *testing.T) {
	tests := []struct {
		ref  string
		want bool
	}{
		{"abcdef1", true},
		{"abcdef1234567", true},
		{"abcdef1234567890abcdef1234567890abcdef12", true}, // full sha also matches
		{"abc", false},                                     // too short (< 7)
		{"", false},
		{"main", false},
		{"ABCDEF1", false}, // uppercase not hex
		{"1234567", true},
		{"123456g", false}, // non-hex char
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			got := IsShortSHA(tc.ref)
			if got != tc.want {
				t.Errorf("IsShortSHA(%q) = %v, want %v", tc.ref, got, tc.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	r := New("github.com", "mytoken")
	if r.Host != "github.com" {
		t.Errorf("expected Host=github.com, got %q", r.Host)
	}
	if r.AuthToken != "mytoken" {
		t.Errorf("expected AuthToken=mytoken, got %q", r.AuthToken)
	}
	if r.Timeout == 0 {
		t.Error("expected non-zero timeout")
	}
}
