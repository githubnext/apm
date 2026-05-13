package gitenv_test

import (
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
