package integrity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache/integrity"
)

func TestReadHeadSHA_PackedRefMultipleRefs(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	packed := "# pack-refs with: peeled fully-peeled\nabc123def456abc123def456abc123def456abc1 refs/heads/other\ndeadbeef11111111deadbeef11111111deadbeef refs/heads/main\n"
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packed), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != "deadbeef11111111deadbeef11111111deadbeef" {
		t.Errorf("ReadHeadSHA packed-refs multiple: got %q", got)
	}
}

func TestReadHeadSHA_RefFileOverridesPackedRefs(t *testing.T) {
	sha := "1111111111111111111111111111111111111111"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	refsHeads := filepath.Join(gitDir, "refs", "heads")
	_ = os.MkdirAll(refsHeads, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	_ = os.WriteFile(filepath.Join(refsHeads, "main"), []byte(sha+"\n"), 0o600)
	packed := "2222222222222222222222222222222222222222 refs/heads/main\n"
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packed), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("loose ref should override packed-refs: got %q want %q", got, sha)
	}
}

func TestVerifyCheckout_SHAPrefixNotMatch(t *testing.T) {
	sha := "aaaa1111aaaa1111aaaa1111aaaa1111aaaa1111"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)
	if integrity.VerifyCheckout(root, "bbbb") {
		t.Error("prefix mismatch should return false")
	}
}

func TestVerifyCheckout_ExactMatch(t *testing.T) {
	sha := "cafecafecafecafecafecafecafecafecafecafe"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)
	if !integrity.VerifyCheckout(root, sha) {
		t.Error("exact SHA should verify as true")
	}
}

func TestVerifyCheckout_NonExistentDir(t *testing.T) {
	result := integrity.VerifyCheckout("/nonexistent/path/xyz/abc", "any-sha")
	if result {
		t.Error("non-existent dir should return false")
	}
}

func TestReadHeadSHA_HeadRefWithWindowsNewline(t *testing.T) {
	sha := "eeee1111eeee1111eeee1111eeee1111eeee1111"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\r\n"), 0o600)
	// Depending on implementation, may or may not trim \r -- just should not panic
	got := integrity.ReadHeadSHA(root)
	_ = got
}

func TestReadHeadSHA_ShortSHA(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("abc123\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != "abc123" {
		t.Errorf("short SHA: got %q want 'abc123'", got)
	}
}
