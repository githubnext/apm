package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type realBehaviorCase struct {
	name   string
	args   []string
	env    map[string]string
	setup  func(t *testing.T, dir string)
	verify func(t *testing.T, dir, stdout, stderr string, code int) bool
}

func TestParityRealFunctionalAndStateDiffContracts(t *testing.T) {
	cases := []realBehaviorCase{
		{
			name: "init creates manifest",
			args: []string{"init", "--yes"},
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "apm.yml"), "dependencies:") && ok
				return ok
			},
		},
		{
			name:  "install local package materializes lock and modules",
			args:  []string{"install", "./packages/local-tools"},
			setup: realBehaviorSetupLocalPackage,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "apm.lock.yaml"), "local-tools") && ok
				ok = realBehaviorExpectDirHasEntries(t, filepath.Join(dir, "apm_modules")) && ok
				return ok
			},
		},
		{
			name:  "compile writes copilot target",
			args:  []string{"compile", "--target", "copilot"},
			setup: realBehaviorSetupProject,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, ".github", "copilot-instructions.md"), "real-behavior") && ok
				return ok
			},
		},
		{
			name:  "pack writes distributable output",
			args:  []string{"pack", "--output", "dist"},
			setup: realBehaviorSetupProjectWithLock,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectDirHasEntries(t, filepath.Join(dir, "dist")) && ok
				return ok
			},
		},
		{
			name:  "run executes project script",
			args:  []string{"run", "stamp"},
			setup: realBehaviorSetupRunnableProject,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "run-stamp.txt"), "real-run") && ok
				return ok
			},
		},
		{
			name:  "audit ci fails on planted hidden unicode",
			args:  []string{"audit", "--ci"},
			setup: realBehaviorSetupAuditFinding,
			verify: func(t *testing.T, _ string, stdout, stderr string, code int) bool {
				if code == 0 {
					t.Errorf("expected non-zero exit for hidden unicode finding\nstdout: %s\nstderr: %s", stdout, stderr)
					return false
				}
				return true
			},
		},
		{
			name:  "mcp install persists manifest dependency",
			args:  []string{"mcp", "install", "example-server"},
			setup: realBehaviorSetupProject,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "apm.yml"), "example-server") && ok
				return ok
			},
		},
		{
			name: "plugin init writes plugin manifest",
			args: []string{"plugin", "init"},
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "plugin.json"), "\"name\"") && ok
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "apm.yml"), "plugin") && ok
				return ok
			},
		},
		{
			name: "marketplace init writes marketplace block",
			args: []string{"marketplace", "init"},
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectFileContains(t, filepath.Join(dir, "apm.yml"), "marketplace:") && ok
				return ok
			},
		},
		{
			name:  "cache clean removes entries but preserves cache root",
			args:  []string{"cache", "clean"},
			env:   map[string]string{"APM_CACHE_DIR": "cache-root"},
			setup: realBehaviorSetupCacheRoot,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				cacheRoot := filepath.Join(dir, "cache-root")
				ok = realBehaviorExpectPathExists(t, cacheRoot) && ok
				ok = realBehaviorExpectPathMissing(t, filepath.Join(cacheRoot, "http_v1", "old", "body")) && ok
				return ok
			},
		},
		{
			name:  "prune removes unreferenced module",
			args:  []string{"prune"},
			setup: realBehaviorSetupStaleModule,
			verify: func(t *testing.T, dir, stdout, stderr string, code int) bool {
				ok := realBehaviorExpectExit(t, stdout, stderr, code, 0)
				ok = realBehaviorExpectPathMissing(t, filepath.Join(dir, "apm_modules", "stale-package")) && ok
				return ok
			},
		},
	}

	functionalPassing := 0
	stateDiffPassing := 0
	defer func() {
		emitCraneRatioGate("functional", functionalPassing, len(cases))
		emitCraneRatioGate("state_diff", stateDiffPassing, len(cases))
	}()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.setup != nil {
				tc.setup(t, dir)
			}
			stdout, stderr, code := realBehaviorRunGoInDir(t, dir, tc.env, tc.args...)
			if tc.verify(t, dir, stdout, stderr, code) {
				functionalPassing++
				stateDiffPassing++
			}
		})
	}
}

func realBehaviorRunGoInDir(t *testing.T, dir string, env map[string]string, args ...string) (string, string, int) {
	t.Helper()
	if goBinPath == "" {
		t.Skip("Go binary not built; skipping")
	}

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(goBinPath, args...)
	cmd.Dir = dir
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run apm %s: %v", strings.Join(args, " "), err)
		}
	}
	return outBuf.String(), errBuf.String(), code
}

func realBehaviorSetupProject(t *testing.T, dir string) {
	t.Helper()
	realBehaviorWriteFile(t, filepath.Join(dir, "apm.yml"), `name: real-behavior
version: 1.0.0
description: Real behavior fixture
author: Crane
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
scripts: {}
`)
	realBehaviorWriteFile(t, filepath.Join(dir, ".apm", "prompts", "real-behavior.md"), "real-behavior prompt\n")
}

func realBehaviorSetupProjectWithLock(t *testing.T, dir string) {
	t.Helper()
	realBehaviorSetupProject(t, dir)
	realBehaviorWriteFile(t, filepath.Join(dir, "apm.lock.yaml"), `lockfile_version: "1"
dependencies: []
local_deployed_files:
  - .apm/prompts/real-behavior.md
local_deployed_file_hashes: {}
`)
}

func realBehaviorSetupRunnableProject(t *testing.T, dir string) {
	t.Helper()
	realBehaviorWriteFile(t, filepath.Join(dir, "apm.yml"), `name: runnable
version: 1.0.0
description: Runnable fixture
author: Crane
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
scripts:
  stamp: "printf real-run > run-stamp.txt"
`)
}

func realBehaviorSetupLocalPackage(t *testing.T, dir string) {
	t.Helper()
	realBehaviorSetupProject(t, dir)
	pkgDir := filepath.Join(dir, "packages", "local-tools")
	realBehaviorWriteFile(t, filepath.Join(pkgDir, "apm.yml"), `name: local-tools
version: 1.0.0
description: Local tools package
author: Crane
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
scripts: {}
`)
	realBehaviorWriteFile(t, filepath.Join(pkgDir, ".apm", "prompts", "tool.md"), "local-tools prompt\n")
}

func realBehaviorSetupAuditFinding(t *testing.T, dir string) {
	t.Helper()
	realBehaviorSetupProjectWithLock(t, dir)
	realBehaviorWriteFile(t, filepath.Join(dir, "apm_modules", "unicode-package", "SKILL.md"), "safe text \u202eevil text\n")
	realBehaviorWriteFile(t, filepath.Join(dir, "apm.lock.yaml"), `lockfile_version: "1"
dependencies:
  - repo_url: local/unicode-package
    resolved_commit: fixture
    deployed_files:
      - apm_modules/unicode-package/SKILL.md
    deployed_file_hashes: {}
`)
}

func realBehaviorSetupCacheRoot(t *testing.T, dir string) {
	t.Helper()
	realBehaviorWriteFile(t, filepath.Join(dir, "cache-root", "http_v1", "old", "body"), "cached\n")
}

func realBehaviorSetupStaleModule(t *testing.T, dir string) {
	t.Helper()
	realBehaviorSetupProjectWithLock(t, dir)
	realBehaviorWriteFile(t, filepath.Join(dir, "apm_modules", "stale-package", "README.md"), "stale\n")
}

func realBehaviorWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create parent dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func realBehaviorExpectExit(t *testing.T, stdout, stderr string, got, want int) bool {
	t.Helper()
	if got != want {
		t.Errorf("exit code = %d, want %d\nstdout: %s\nstderr: %s", got, want, stdout, stderr)
		return false
	}
	return true
}

func realBehaviorExpectFileContains(t *testing.T, path, needle string) bool {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("expected file %s to exist: %v", path, err)
		return false
	}
	if !strings.Contains(string(content), needle) {
		t.Errorf("expected %s to contain %q, got:\n%s", path, needle, string(content))
		return false
	}
	return true
}

func realBehaviorExpectPathExists(t *testing.T, path string) bool {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected path %s to exist: %v", path, err)
		return false
	}
	return true
}

func realBehaviorExpectPathMissing(t *testing.T, path string) bool {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected path %s to be removed", path)
		return false
	} else if !os.IsNotExist(err) {
		t.Errorf("expected path %s to be absent, got: %v", path, err)
		return false
	}
	return true
}

func realBehaviorExpectDirHasEntries(t *testing.T, path string) bool {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Errorf("expected directory %s to exist: %v", path, err)
		return false
	}
	if len(entries) == 0 {
		t.Errorf("expected directory %s to contain at least one entry", path)
		return false
	}
	return true
}
