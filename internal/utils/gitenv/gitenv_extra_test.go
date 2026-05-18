package gitenv_test

import (
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/gitenv"
)

func TestGitSubprocessEnv_NoGitDir(t *testing.T) {
	orig := os.Getenv("GIT_DIR")
	os.Setenv("GIT_DIR", "/some/git/dir")
	defer os.Setenv("GIT_DIR", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_DIR=") {
			t.Errorf("GIT_DIR should be stripped, but found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitWorkTree(t *testing.T) {
	orig := os.Getenv("GIT_WORK_TREE")
	os.Setenv("GIT_WORK_TREE", "/some/work/tree")
	defer os.Setenv("GIT_WORK_TREE", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_WORK_TREE=") {
			t.Errorf("GIT_WORK_TREE should be stripped, but found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitIndexFile(t *testing.T) {
	orig := os.Getenv("GIT_INDEX_FILE")
	os.Setenv("GIT_INDEX_FILE", "/tmp/index")
	defer os.Setenv("GIT_INDEX_FILE", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_INDEX_FILE=") {
			t.Errorf("GIT_INDEX_FILE should be stripped, but found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitObjectDirectory(t *testing.T) {
	orig := os.Getenv("GIT_OBJECT_DIRECTORY")
	os.Setenv("GIT_OBJECT_DIRECTORY", "/objects")
	defer os.Setenv("GIT_OBJECT_DIRECTORY", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_OBJECT_DIRECTORY=") {
			t.Errorf("GIT_OBJECT_DIRECTORY should be stripped, but found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitNamespace(t *testing.T) {
	orig := os.Getenv("GIT_NAMESPACE")
	os.Setenv("GIT_NAMESPACE", "testns")
	defer os.Setenv("GIT_NAMESPACE", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_NAMESPACE=") {
			t.Errorf("GIT_NAMESPACE should be stripped, but found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_PreservesNonGitVars(t *testing.T) {
	orig := os.Getenv("MY_CUSTOM_VAR")
	os.Setenv("MY_CUSTOM_VAR", "myvalue")
	defer os.Setenv("MY_CUSTOM_VAR", orig)

	env := gitenv.GitSubprocessEnv()
	found := false
	for _, kv := range env {
		if kv == "MY_CUSTOM_VAR=myvalue" {
			found = true
			break
		}
	}
	if !found {
		t.Error("MY_CUSTOM_VAR should be preserved in subprocess env")
	}
}

func TestGitSubprocessEnv_IsSlice(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	if env == nil {
		t.Error("expected non-nil slice from GitSubprocessEnv")
	}
}

func TestGitSubprocessEnv_AllKeyValuePairs(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if !strings.Contains(kv, "=") {
			t.Errorf("expected KEY=VALUE format, got %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitCeilingDirectories(t *testing.T) {
	orig := os.Getenv("GIT_CEILING_DIRECTORIES")
	os.Setenv("GIT_CEILING_DIRECTORIES", "/home")
	defer os.Setenv("GIT_CEILING_DIRECTORIES", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_CEILING_DIRECTORIES=") {
			t.Errorf("GIT_CEILING_DIRECTORIES should be stripped, found %q", kv)
		}
	}
}

func TestGitSubprocessEnv_NoGitDiscoveryAcrossFilesystem(t *testing.T) {
	orig := os.Getenv("GIT_DISCOVERY_ACROSS_FILESYSTEM")
	os.Setenv("GIT_DISCOVERY_ACROSS_FILESYSTEM", "1")
	defer os.Setenv("GIT_DISCOVERY_ACROSS_FILESYSTEM", orig)

	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_DISCOVERY_ACROSS_FILESYSTEM=") {
			t.Errorf("GIT_DISCOVERY_ACROSS_FILESYSTEM should be stripped, found %q", kv)
		}
	}
}

func TestResetGitCache_AllowsReinitialization(t *testing.T) {
	gitenv.ResetGitCache()
	exe1, err1 := gitenv.GetGitExecutable()
	gitenv.ResetGitCache()
	exe2, err2 := gitenv.GetGitExecutable()
	if err1 != err2 {
		t.Errorf("errors differ after reset: %v vs %v", err1, err2)
	}
	if exe1 != exe2 {
		t.Errorf("executables differ after reset: %q vs %q", exe1, exe2)
	}
}

func TestGetGitExecutable_ReturnsPath(t *testing.T) {
	gitenv.ResetGitCache()
	defer gitenv.ResetGitCache()
	exe, err := gitenv.GetGitExecutable()
	if err != nil {
		t.Skipf("git not available: %v", err)
	}
	if !strings.Contains(exe, "git") {
		t.Errorf("expected path containing 'git', got %q", exe)
	}
}
