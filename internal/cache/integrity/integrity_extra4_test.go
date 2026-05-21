package integrity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache/integrity"
)

func TestReadHeadSHA_NoGitDir_Empty_Extra4(t *testing.T) {
	dir := t.TempDir()
	got := integrity.ReadHeadSHA(dir)
	if got != "" {
		t.Errorf("expected empty for dir without .git, got %q", got)
	}
}

func TestReadHeadSHA_WithDetachedSHA_Extra4(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o644); err != nil {
		t.Fatalf("write HEAD: %v", err)
	}
	got := integrity.ReadHeadSHA(dir)
	if got != sha {
		t.Errorf("expected %q, got %q", sha, got)
	}
}

func TestVerifyCheckout_NonexistentDir_Extra4(t *testing.T) {
	ok := integrity.VerifyCheckout("/nonexistent/dir", "abc123")
	if ok {
		t.Error("expected false for nonexistent dir")
	}
}

func TestVerifyCheckout_EmptySHA_Extra4(t *testing.T) {
	dir := t.TempDir()
	ok := integrity.VerifyCheckout(dir, "")
	if ok {
		t.Error("expected false when no .git present and empty expected SHA")
	}
}

func TestVerifyCheckout_MatchingFullSHA_Extra4(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	sha := "1234567890abcdef1234567890abcdef12345678"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o644); err != nil {
		t.Fatalf("write HEAD: %v", err)
	}
	ok := integrity.VerifyCheckout(dir, sha)
	if !ok {
		t.Error("expected true for matching SHA")
	}
}

func TestVerifyCheckout_MismatchSHA_Extra4(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	sha := "1234567890abcdef1234567890abcdef12345678"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o644); err != nil {
		t.Fatalf("write HEAD: %v", err)
	}
	ok := integrity.VerifyCheckout(dir, "different000000000000000000000000000000")
	if ok {
		t.Error("expected false for mismatched SHA")
	}
}
