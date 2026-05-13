// Package gitenv provides cached git binary lookup and subprocess
// environment sanitization. Mirrors src/apm_cli/utils/git_env.py.
package gitenv

import (
	"errors"
	"os"
	"os/exec"
	"sync"
)

// stripGitVars lists ambient git state variables that are stripped from
// subprocess environments to avoid biasing APM's git operations.
var stripGitVars = map[string]struct{}{
	"GIT_DIR":                           {},
	"GIT_WORK_TREE":                     {},
	"GIT_INDEX_FILE":                    {},
	"GIT_OBJECT_DIRECTORY":              {},
	"GIT_ALTERNATE_OBJECT_DIRECTORIES":  {},
	"GIT_COMMON_DIR":                    {},
	"GIT_NAMESPACE":                     {},
	"GIT_INDEX_VERSION":                 {},
	"GIT_CEILING_DIRECTORIES":           {},
	"GIT_DISCOVERY_ACROSS_FILESYSTEM":   {},
	"GIT_REPLACE_REF_BASE":              {},
	"GIT_GRAFTS_FILE":                   {},
	"GIT_SHALLOW_FILE":                  {},
}

var (
	once          sync.Once
	gitExecutable string
	gitErr        error
)

// GetGitExecutable returns the path to the git executable (cached after first lookup).
func GetGitExecutable() (string, error) {
	once.Do(func() {
		gitExecutable, gitErr = exec.LookPath("git")
		if gitErr != nil {
			gitErr = errors.New("git executable not found on PATH. Please install git: https://git-scm.com/downloads")
		}
	})
	return gitExecutable, gitErr
}

// GitSubprocessEnv returns a sanitized environment slice for git subprocesses.
// Strips ambient git state variables while preserving user-controlled configuration.
func GitSubprocessEnv() []string {
	env := os.Environ()
	result := make([]string, 0, len(env))
	for _, kv := range env {
		key := kv
		for i, c := range kv {
			if c == '=' {
				key = kv[:i]
				break
			}
		}
		if _, strip := stripGitVars[key]; !strip {
			result = append(result, kv)
		}
	}
	return result
}

// ResetGitCache resets the cached git executable (for testing purposes only).
func ResetGitCache() {
	once = sync.Once{}
	gitExecutable = ""
	gitErr = nil
}
