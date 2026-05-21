package integrity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache/integrity"
)

func TestReadHeadSHA_DetachedHEAD(t *testing.T) {
	sha := "1234567890abcdef1234567890abcdef12345678"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("detached HEAD: got %q want %q", got, sha)
	}
}

func TestReadHeadSHA_PackedRefNoMatch(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/missing\n"), 0o600)
	packed := "aaaa1111aaaa1111aaaa1111aaaa1111aaaa1111 refs/heads/other\n"
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packed), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != "" {
		t.Errorf("unmatched ref should return empty, got %q", got)
	}
}

func TestReadHeadSHA_GitWorktreeFile(t *testing.T) {
	root := t.TempDir()
	realGitDir := filepath.Join(root, "real_git")
	sha := "cafebabe1234cafebabe1234cafebabe12345678"
	refsHeads := filepath.Join(realGitDir, "refs", "heads")
	_ = os.MkdirAll(refsHeads, 0o700)
	_ = os.WriteFile(filepath.Join(realGitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	_ = os.WriteFile(filepath.Join(refsHeads, "main"), []byte(sha+"\n"), 0o600)
	// Use relative path "real_git" from the root so Abs resolves correctly
	_ = os.WriteFile(filepath.Join(root, ".git"), []byte("gitdir: real_git\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("worktree gitdir file: got %q want %q", got, sha)
	}
}

func TestReadHeadSHA_InvalidGitdirFile(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, ".git"), []byte("not-a-gitdir-line\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != "" {
		t.Errorf("invalid gitdir file should return empty, got %q", got)
	}
}

func TestVerifyCheckout_EmptySHA(t *testing.T) {
	dir := t.TempDir()
	result := integrity.VerifyCheckout(dir, "")
	if result {
		t.Error("empty expected SHA should not verify")
	}
}

func TestVerifyCheckout_BothEmpty(t *testing.T) {
	dir := t.TempDir()
	result := integrity.VerifyCheckout(dir, "")
	if result {
		t.Error("empty expected SHA with no .git dir should not verify")
	}
}

func TestReadHeadSHA_MissingHEADFile(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	// No HEAD file
	got := integrity.ReadHeadSHA(root)
	if got != "" {
		t.Errorf("missing HEAD file should return empty, got %q", got)
	}
}

func TestReadHeadSHA_NonExistentDir(t *testing.T) {
	got := integrity.ReadHeadSHA("/nonexistent/xyz123/abc")
	if got != "" {
		t.Errorf("non-existent dir should return empty, got %q", got)
	}
}
