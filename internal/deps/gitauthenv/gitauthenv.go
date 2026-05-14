// Package gitauthenv builds the various git environment dicts the downloader needs.
// Migrated from src/apm_cli/deps/git_auth_env.py
//
// Three env flavours:
//  1. SetupEnvironment   -- auth-bearing env for primary git ops
//  2. NoninteractiveEnv  -- non-auth env for unauthenticated fallback
//  3. SubprocessEnvDict  -- sanitized env for cache-layer subprocess calls
package gitauthenv

import (
	"os"
	"runtime"
	"strings"
)

// GitAuthEnvBuilder builds the various git env dicts the downloader needs.
type GitAuthEnvBuilder struct {
	baseEnv map[string]string
}

// New returns a new GitAuthEnvBuilder.
// baseEnv is the auth-bearing environment provided by the token manager
// (analogous to token_manager.setup_environment() in Python).
func New(baseEnv map[string]string) *GitAuthEnvBuilder {
	return &GitAuthEnvBuilder{baseEnv: baseEnv}
}

// SetupEnvironment builds the auth-bearing primary git env.
// Sets GIT_TERMINAL_PROMPT, GIT_ASKPASS, GIT_CONFIG_NOSYSTEM,
// GIT_SSH_COMMAND (with ConnectTimeout=30), and GIT_CONFIG_GLOBAL.
func (b *GitAuthEnvBuilder) SetupEnvironment() map[string]string {
	env := copyEnv(b.baseEnv)

	env["GIT_TERMINAL_PROMPT"] = "0"
	env["GIT_ASKPASS"] = "echo"
	env["GIT_CONFIG_NOSYSTEM"] = "1"

	// Ensure SSH connections fail fast (30 s timeout).
	const sshTimeout = "-o ConnectTimeout=30"
	existingSSH := strings.TrimSpace(os.Getenv("GIT_SSH_COMMAND"))
	if existingSSH != "" {
		if !strings.Contains(strings.ToLower(existingSSH), "connecttimeout") {
			env["GIT_SSH_COMMAND"] = existingSSH + " " + sshTimeout
		} else {
			env["GIT_SSH_COMMAND"] = existingSSH
		}
	} else {
		env["GIT_SSH_COMMAND"] = "ssh " + sshTimeout
	}

	if runtime.GOOS == "windows" {
		// On Windows, point GIT_CONFIG_GLOBAL at an empty file.
		tmpDir := os.TempDir()
		emptyCfg := tmpDir + "\\.apm_empty_gitconfig"
		// Create the empty file (ignore errors -- best-effort).
		f, err := os.OpenFile(emptyCfg, os.O_CREATE|os.O_WRONLY, 0o644)
		if err == nil {
			f.Close()
		}
		env["GIT_CONFIG_GLOBAL"] = emptyCfg
	} else {
		env["GIT_CONFIG_GLOBAL"] = "/dev/null"
	}

	return env
}

// NoninteractiveEnvOptions controls the credential-helper suppression fence.
type NoninteractiveEnvOptions struct {
	// PreserveConfigIsolation keeps GIT_CONFIG_NOSYSTEM and GIT_CONFIG_GLOBAL.
	PreserveConfigIsolation bool
	// SuppressCredentialHelpers applies the full credential-helper fence
	// (use for HTTP transport to avoid leaking tokens in plaintext).
	SuppressCredentialHelpers bool
}

// NoninteractiveEnv builds a non-interactive git env for unauthenticated operations.
//
// Credential-helper policy (two-stage):
//  1. Always clear GIT_ASKPASS so system credential helpers resolve naturally.
//  2. Re-set the full suppression fence only when SuppressCredentialHelpers is true.
func NoninteractiveEnv(baseGitEnv map[string]string, opts NoninteractiveEnvOptions) map[string]string {
	env := copyEnv(baseGitEnv)

	env["GIT_TERMINAL_PROMPT"] = "0"
	delete(env, "GIT_ASKPASS")

	if opts.PreserveConfigIsolation || opts.SuppressCredentialHelpers {
		env["GIT_CONFIG_NOSYSTEM"] = "1"
		if v, ok := baseGitEnv["GIT_CONFIG_GLOBAL"]; ok {
			env["GIT_CONFIG_GLOBAL"] = v
		}
	} else {
		delete(env, "GIT_CONFIG_GLOBAL")
		delete(env, "GIT_CONFIG_NOSYSTEM")
	}

	if opts.SuppressCredentialHelpers {
		env["GIT_ASKPASS"] = "echo"
		env["GIT_CONFIG_COUNT"] = "1"
		env["GIT_CONFIG_KEY_0"] = "credential.helper"
		env["GIT_CONFIG_VALUE_0"] = ""
	} else {
		delete(env, "GIT_CONFIG_COUNT")
		delete(env, "GIT_CONFIG_KEY_0")
		delete(env, "GIT_CONFIG_VALUE_0")
	}

	return env
}

// SubprocessEnvDict returns a sanitized git env dict for cache-layer subprocess calls.
// Merges the auth-aware baseGitEnv over a sanitized ambient env so the subprocess
// never inherits a stray GIT_DIR or GIT_CEILING_DIRECTORIES.
func SubprocessEnvDict(baseGitEnv map[string]string) map[string]string {
	env := gitSubprocessEnv()
	for k, v := range baseGitEnv {
		env[k] = v
	}
	return env
}

// gitSubprocessEnv returns the current process environment with git-state variables
// stripped so cache-layer subprocess calls start with a clean slate.
func gitSubprocessEnv() map[string]string {
	stripKeys := map[string]bool{
		"GIT_DIR":                  true,
		"GIT_CEILING_DIRECTORIES":  true,
		"GIT_WORK_TREE":            true,
		"GIT_INDEX_FILE":           true,
		"GIT_OBJECT_DIRECTORY":     true,
		"GIT_ALTERNATE_OBJECT_DIRECTORIES": true,
	}
	env := make(map[string]string)
	for _, kv := range os.Environ() {
		idx := strings.IndexByte(kv, '=')
		if idx < 0 {
			continue
		}
		k, v := kv[:idx], kv[idx+1:]
		if !stripKeys[k] {
			env[k] = v
		}
	}
	return env
}

func copyEnv(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
