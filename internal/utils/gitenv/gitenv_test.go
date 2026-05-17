package gitenv_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/gitenv"
)

func TestGetGitExecutable(t *testing.T) {
	gitenv.ResetGitCache()
	path, err := gitenv.GetGitExecutable()
	if err != nil {
		t.Fatalf("GetGitExecutable: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty git path")
	}
}

func TestGitSubprocessEnv(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	if len(env) == 0 {
		t.Error("expected non-empty env")
	}
	for _, kv := range env {
		for _, stripped := range []string{
			"GIT_DIR=", "GIT_WORK_TREE=", "GIT_INDEX_FILE=",
		} {
			if len(kv) >= len(stripped) && kv[:len(stripped)] == stripped {
				t.Errorf("env contains stripped var: %s", kv)
			}
		}
	}
}

func TestGitSubprocessEnvStripsAllKnownVars(t *testing.T) {
	strippedVars := []string{
		"GIT_DIR",
		"GIT_WORK_TREE",
		"GIT_INDEX_FILE",
		"GIT_OBJECT_DIRECTORY",
		"GIT_ALTERNATE_OBJECT_DIRECTORIES",
		"GIT_COMMON_DIR",
		"GIT_NAMESPACE",
		"GIT_INDEX_VERSION",
		"GIT_CEILING_DIRECTORIES",
		"GIT_DISCOVERY_ACROSS_FILESYSTEM",
		"GIT_REPLACE_REF_BASE",
		"GIT_GRAFTS_FILE",
		"GIT_SHALLOW_FILE",
	}
	env := gitenv.GitSubprocessEnv()
	envMap := make(map[string]bool)
	for _, kv := range env {
		idx := strings.IndexByte(kv, '=')
		if idx > 0 {
			envMap[kv[:idx]] = true
		}
	}
	for _, v := range strippedVars {
		if envMap[v] {
			t.Errorf("env should not contain stripped variable %q", v)
		}
	}
}

func TestGetGitExecutableCached(t *testing.T) {
	gitenv.ResetGitCache()
	path1, err1 := gitenv.GetGitExecutable()
	if err1 != nil {
		t.Skipf("git not found: %v", err1)
	}
	// Second call should return the same cached result.
	path2, err2 := gitenv.GetGitExecutable()
	if err2 != nil {
		t.Fatalf("second GetGitExecutable: %v", err2)
	}
	if path1 != path2 {
		t.Errorf("cached path mismatch: %q vs %q", path1, path2)
	}
}

func TestGitSubprocessEnvKeyValueFormat(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	for _, kv := range env {
		if !strings.Contains(kv, "=") {
			t.Errorf("env entry %q does not contain '='", kv)
		}
	}
}

func TestGitSubprocessEnvPreservesPath(t *testing.T) {
	env := gitenv.GitSubprocessEnv()
	found := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			found = true
			break
		}
	}
	if !found {
		t.Error("GitSubprocessEnv should preserve PATH")
	}
}

func TestResetGitCacheAllowsReinit(t *testing.T) {
	gitenv.ResetGitCache()
	p1, err := gitenv.GetGitExecutable()
	if err != nil {
		t.Skipf("git not found: %v", err)
	}
	gitenv.ResetGitCache()
	p2, err := gitenv.GetGitExecutable()
	if err != nil {
		t.Fatalf("re-init GetGitExecutable: %v", err)
	}
	if p1 != p2 {
		t.Errorf("path changed after reset: %q vs %q", p1, p2)
	}
}
