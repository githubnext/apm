package mcp

import "testing"

func TestSearchOptions_LimitZeroDefault(t *testing.T) {
	// Limit of 0 should be handled gracefully by callers (default to 20).
	opts := SearchOptions{Query: "test", Limit: 0}
	if opts.Limit != 0 {
		t.Errorf("Limit should remain 0 in struct, got %d", opts.Limit)
	}
}

func TestTruncate_ShortString(t *testing.T) {
	// Strings shorter than n should be returned as-is.
	got := truncate("hello", 10)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestTruncate_EmptyString(t *testing.T) {
	got := truncate("", 5)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestTruncate_NGreaterThanLen(t *testing.T) {
	got := truncate("abc", 100)
	if got != "abc" {
		t.Errorf("expected 'abc', got %q", got)
	}
}

func TestTruncate_ExactBoundary(t *testing.T) {
	// A string of exactly n chars should not be truncated.
	s := "abcde"
	got := truncate(s, 5)
	if got != "abcde" {
		t.Errorf("expected 'abcde', got %q", got)
	}
}

func TestTruncate_TruncatesWithEllipsis(t *testing.T) {
	// Strings longer than n should be truncated.
	s := "hello world this is a long string"
	got := truncate(s, 10)
	if len([]rune(got)) > 10 {
		t.Errorf("truncated string too long: %q (len %d)", got, len([]rune(got)))
	}
}

func TestInstallOptions_UserScope(t *testing.T) {
	opts := InstallOptions{ServerRef: "pkg/name", UserScope: true}
	if !opts.UserScope {
		t.Error("UserScope should be true")
	}
}

func TestInstallOptions_ForceAndRuntime(t *testing.T) {
	opts := InstallOptions{Force: true, Runtime: "node"}
	if !opts.Force {
		t.Error("Force should be true")
	}
	if opts.Runtime != "node" {
		t.Errorf("Runtime = %q", opts.Runtime)
	}
}

func TestInfoOptions_EmptyFormat(t *testing.T) {
	opts := InfoOptions{ServerRef: "pkg/name"}
	if opts.Format != "" {
		t.Errorf("Format should be empty by default, got %q", opts.Format)
	}
}
