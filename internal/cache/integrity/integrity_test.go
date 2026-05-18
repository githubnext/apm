package integrity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache/integrity"
)

func makeGitDir(t *testing.T, sha string) string {
	t.Helper()
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	_ = os.MkdirAll(filepath.Join(gitDir, "refs", "heads"), 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	_ = os.WriteFile(filepath.Join(gitDir, "refs", "heads", "main"), []byte(sha+"\n"), 0o600)
	return dir
}

func TestReadHeadSHADetachedHead(t *testing.T) {
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)

	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA = %q, want %q", got, sha)
	}
}

func TestReadHeadSHASymbolicRef(t *testing.T) {
	sha := "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	root := makeGitDir(t, sha)

	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA = %q, want %q", got, sha)
	}
}

func TestReadHeadSHAPackedRefs(t *testing.T) {
	sha := "1111222233334444555566667777888899990000"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/feature\n"), 0o600)
	packedRefs := "# pack-refs with: peeled fully-peeled sorted\n" + sha + " refs/heads/feature\n"
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packedRefs), 0o600)

	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA from packed-refs = %q, want %q", got, sha)
	}
}

func TestReadHeadSHANoGitDir(t *testing.T) {
	root := t.TempDir()
	got := integrity.ReadHeadSHA(root)
	if got != "" {
		t.Errorf("ReadHeadSHA on non-git dir = %q, want empty", got)
	}
}

func TestVerifyCheckout(t *testing.T) {
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	root := makeGitDir(t, sha)

	if !integrity.VerifyCheckout(root, sha) {
		t.Error("VerifyCheckout should return true for matching SHA")
	}
	if integrity.VerifyCheckout(root, "wrongsha") {
		t.Error("VerifyCheckout should return false for mismatched SHA")
	}
}

func TestVerifyCheckoutNonGitDir(t *testing.T) {
	root := t.TempDir()
	if integrity.VerifyCheckout(root, "anySHA") {
		t.Error("VerifyCheckout should return false for non-git dir")
	}
}
