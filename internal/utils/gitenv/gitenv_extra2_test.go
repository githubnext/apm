package gitenv_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/gitenv"
)

func TestGetGitExecutable_ReturnsNonEmptyPath(t *testing.T) {
	path, err := gitenv.GetGitExecutable()
	if err != nil {
		t.Fatalf("GetGitExecutable: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty git path")
	}
}

func TestGetGitExecutable_IsCached(t *testing.T) {
	p1, err1 := gitenv.GetGitExecutable()
	p2, err2 := gitenv.GetGitExecutable()
	if err1 != nil || err2 != nil {
		t.Fatalf("errors: %v / %v", err1, err2)
	}
	if p1 != p2 {
		t.Errorf("expected same cached path, got %q vs %q", p1, p2)
	}
}

func TestGitSubprocessEnv_NotEmpty(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	if len(env) == 0 {
		t.Error("expected non-empty environment slice")
	}
}

func TestGitSubprocessEnv_StripsGITDIR(t *testing.T) {
	t.Setenv("GIT_DIR", "/some/dir")
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_DIR=") {
			t.Errorf("GIT_DIR should be stripped, but found: %q", kv)
		}
	}
}

func TestGitSubprocessEnv_StripsGITWORKTREE(t *testing.T) {
	t.Setenv("GIT_WORK_TREE", "/work")
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_WORK_TREE=") {
			t.Errorf("GIT_WORK_TREE should be stripped, but found: %q", kv)
		}
	}
}

func TestGitSubprocessEnv_StripsGITINDEXFILE(t *testing.T) {
	t.Setenv("GIT_INDEX_FILE", "/index")
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_INDEX_FILE=") {
			t.Errorf("GIT_INDEX_FILE should be stripped, but found: %q", kv)
		}
	}
}

func TestGitSubprocessEnv_PreservesCustomVar(t *testing.T) {
	t.Setenv("MY_CUSTOM_VAR", "hello")
	env := gitenv.GitSubprocessEnv()
	found := false
	for _, kv := range env {
		if kv == "MY_CUSTOM_VAR=hello" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected MY_CUSTOM_VAR to be preserved in environment")
	}
}

func TestGitSubprocessEnv_StripsGITNAMESPACE(t *testing.T) {
	t.Setenv("GIT_NAMESPACE", "myns")
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_NAMESPACE=") {
			t.Errorf("GIT_NAMESPACE should be stripped, but found: %q", kv)
		}
	}
}

func TestGitSubprocessEnv_StripsGITCEILINGDIRECTORIES(t *testing.T) {
	t.Setenv("GIT_CEILING_DIRECTORIES", "/ceiling")
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_CEILING_DIRECTORIES=") {
			t.Errorf("GIT_CEILING_DIRECTORIES should be stripped, but found: %q", kv)
		}
	}
}

func TestGitSubprocessEnv_EnvSliceIsKeyValuePairs(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if !strings.Contains(kv, "=") {
			t.Errorf("env entry %q should contain '='", kv)
		}
	}
}
