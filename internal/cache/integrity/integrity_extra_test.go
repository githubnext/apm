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

func TestVerifyCheckout_Match(t *testing.T) {
	sha := "aabbccddaabbccddaabbccddaabbccddaabbccdd"
	root := makeGitDir(t, sha)
	if !integrity.VerifyCheckout(root, sha) {
		t.Error("expected VerifyCheckout true for matching SHA")
	}
}

func TestVerifyCheckout_Mismatch(t *testing.T) {
	sha := "aabbccddaabbccddaabbccddaabbccddaabbccdd"
	root := makeGitDir(t, sha)
	if integrity.VerifyCheckout(root, "differentsha000000000000000000000000000000") {
		t.Error("expected VerifyCheckout false for mismatched SHA")
	}
}

func TestReadHeadSHA_RefPointingToFile(t *testing.T) {
	sha := "1234567890abcdef1234567890abcdef12345678"
	root := makeGitDir(t, sha)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA ref to file: got %q want %q", got, sha)
	}
}

func TestReadHeadSHA_WorktreeGitFile(t *testing.T) {
	// .git is a file pointing to a relative gitdir
	sha := "cafecafecafecafecafecafecafecafecafecafe"
	root := t.TempDir()
	realGitDir := filepath.Join(root, "dotgit")
	_ = os.MkdirAll(filepath.Join(realGitDir, "refs", "heads"), 0o700)
	_ = os.WriteFile(filepath.Join(realGitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o600)
	_ = os.WriteFile(filepath.Join(realGitDir, "refs", "heads", "main"), []byte(sha+"\n"), 0o600)
	worktree := filepath.Join(root, "worktree")
	_ = os.MkdirAll(worktree, 0o755)
	_ = os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: ../dotgit\n"), 0o600)
	got := integrity.ReadHeadSHA(worktree)
	if got != sha {
		t.Errorf("ReadHeadSHA worktree gitfile: got %q want %q", got, sha)
	}
}

func TestReadHeadSHA_PackedRefsSimple(t *testing.T) {
	sha := "0000000000000000000000000000000000000001"
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	_ = os.MkdirAll(gitDir, 0o700)
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/feature\n"), 0o600)
	packed := sha + " refs/heads/feature\n"
	_ = os.WriteFile(filepath.Join(gitDir, "packed-refs"), []byte(packed), 0o600)
	got := integrity.ReadHeadSHA(root)
	if got != sha {
		t.Errorf("ReadHeadSHA packed-refs simple: got %q want %q", got, sha)
	}
}
