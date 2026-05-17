package integrity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache/integrity"
)

func TestReadHeadSHA_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	got := integrity.ReadHeadSHA(dir)
	if got != "" {
		t.Errorf("expected empty for dir without .git, got %q", got)
	}
}

func TestReadHeadSHA_DirectSHA(t *testing.T) {
	sha := "cafebabe00000000cafebabe00000000cafebabe"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA direct SHA: got %q want %q", got, sha)
	}
}

func TestReadHeadSHA_RefPointingToMissingFile(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/missing\n"), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != "" {
		t.Errorf("expected empty for dangling ref, got %q", got)
	}
}

func TestVerifyCheckout_EmptyExpected(t *testing.T) {
	sha := "aabbccddaabbccddaabbccddaabbccddaabbccdd"
	root := makeGitDir(t, sha)
	if integrity.VerifyCheckout(root, "") {
		t.Error("VerifyCheckout should be false when expectedSHA is empty")
	}
}

func TestVerifyCheckout_EmptyActual(t *testing.T) {
	root := t.TempDir()
	if integrity.VerifyCheckout(root, "anySHA") {
		t.Error("VerifyCheckout should be false when ReadHeadSHA returns empty")
	}
}

func TestReadHeadSHA_PackedRefsWithCommentLine(t *testing.T) {
	sha := "deadcafedeadcafedeadcafedeadcafedeadcafe"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	packedRefs := "# pack-refs with: peeled fully-peeled sorted\n^peeled-sha\n" + sha + " refs/heads/main\n"
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packedRefs), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA packed-refs with comment: got %q want %q", got, sha)
	}
}
