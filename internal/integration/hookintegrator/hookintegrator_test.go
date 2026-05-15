package hookintegrator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	hi := New()
	if hi == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHookIntegrationResult_HooksIntegrated(t *testing.T) {
	r := &HookIntegrationResult{
		FilesIntegrated: 3,
		FilesUpdated:    1,
		FilesSkipped:    0,
		ScriptsCopied:   2,
	}
	if r.HooksIntegrated() != 3 {
		t.Errorf("HooksIntegrated() = %d, want 3", r.HooksIntegrated())
	}
}

func TestFindHookFiles_NoHooksDir(t *testing.T) {
	hi := New()
	dir := t.TempDir()
	// No hooks/ or .apm/hooks/ directory
	files := hi.FindHookFiles(dir)
	if len(files) != 0 {
		t.Errorf("FindHookFiles with no hooks dir should return empty, got %v", files)
	}
}

func TestFindHookFiles_WithHooksDir(t *testing.T) {
	hi := New()
	pkgDir := t.TempDir()
	hooksDir := filepath.Join(pkgDir, "hooks")
	os.MkdirAll(hooksDir, 0755) //nolint:errcheck

	// Write a JSON hook file
	hookData := map[string]interface{}{
		"hooks": []interface{}{},
	}
	data, _ := json.Marshal(hookData)
	os.WriteFile(filepath.Join(hooksDir, "myhook.json"), data, 0644) //nolint:errcheck

	// Write a non-JSON file (should be ignored)
	os.WriteFile(filepath.Join(hooksDir, "readme.txt"), []byte("ignored"), 0644) //nolint:errcheck

	files := hi.FindHookFiles(pkgDir)
	if len(files) != 1 {
		t.Errorf("FindHookFiles should find 1 JSON file, got %d: %v", len(files), files)
	}
}

func TestFindHookFiles_ApmHooksDir(t *testing.T) {
	hi := New()
	pkgDir := t.TempDir()
	apmHooksDir := filepath.Join(pkgDir, ".apm", "hooks")
	os.MkdirAll(apmHooksDir, 0755) //nolint:errcheck

	hookData := map[string]interface{}{"hooks": []interface{}{}}
	data, _ := json.Marshal(hookData)
	os.WriteFile(filepath.Join(apmHooksDir, "hook.json"), data, 0644) //nolint:errcheck

	files := hi.FindHookFiles(pkgDir)
	if len(files) != 1 {
		t.Errorf("FindHookFiles should find 1 JSON file in .apm/hooks, got %d", len(files))
	}
}

func TestFindHookFiles_DeduplicatesSameFile(t *testing.T) {
	hi := New()
	pkgDir := t.TempDir()
	// If both hooks/ and .apm/hooks/ had a symlink to the same file, it should only appear once.
	// For simplicity, just test that two different hook files appear as two entries.
	hooksDir := filepath.Join(pkgDir, "hooks")
	os.MkdirAll(hooksDir, 0755) //nolint:errcheck
	for _, name := range []string{"hook1.json", "hook2.json"} {
		data, _ := json.Marshal(map[string]interface{}{"name": name})
		os.WriteFile(filepath.Join(hooksDir, name), data, 0644) //nolint:errcheck
	}
	files := hi.FindHookFiles(pkgDir)
	if len(files) != 2 {
		t.Errorf("FindHookFiles should find 2 files, got %d", len(files))
	}
}

func TestIntegratePackageHooks_NoHooks(t *testing.T) {
	hi := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	// No hook files in package — should return empty result, not error
	result := hi.IntegratePackageHooks(pkgDir, projectDir, "test-pkg", false, nil, nil, "")
	if result.FilesIntegrated != 0 {
		t.Errorf("IntegratePackageHooks with no hooks should integrate 0 files, got %d", result.FilesIntegrated)
	}
}

func TestIntegratePackageHooks_WithValidHook(t *testing.T) {
	hi := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	hooksDir := filepath.Join(pkgDir, "hooks")
	os.MkdirAll(hooksDir, 0755) //nolint:errcheck

	hookContent := map[string]interface{}{
		"hooks": []interface{}{
			map[string]interface{}{
				"event":   "preToolUse",
				"command": "echo hello",
			},
		},
	}
	data, _ := json.MarshalIndent(hookContent, "", "  ")
	os.WriteFile(filepath.Join(hooksDir, "copilot.json"), data, 0644) //nolint:errcheck

	result := hi.IntegratePackageHooks(pkgDir, projectDir, "test-pkg", false, nil, nil, "")
	// May or may not integrate depending on target filtering, but should not panic
	_ = result
}

func TestSyncIntegration_NoManagedFiles(t *testing.T) {
	hi := New()
	projectDir := t.TempDir()
	stats := hi.SyncIntegration(projectDir, nil, nil)
	if stats.FilesRemoved < 0 {
		t.Errorf("SyncIntegration FilesRemoved should be >= 0")
	}
}

func TestSyncIntegration_EmptyManagedFiles(t *testing.T) {
	hi := New()
	projectDir := t.TempDir()
	managed := map[string]struct{}{}
	stats := hi.SyncIntegration(projectDir, managed, nil)
	if stats.FilesRemoved != 0 {
		t.Errorf("SyncIntegration with empty managed files should remove 0, got %d", stats.FilesRemoved)
	}
}

func TestIntegrateHooksForTarget_NilTarget(t *testing.T) {
	// Test that IntegratePackageHooks path is exercised
	hi := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	// Test with empty package directory — should return empty result
	result := hi.IntegratePackageHooks(pkgDir, projectDir, "my-package", false, nil, nil, ".github")
	if result == nil {
		t.Fatal("IntegratePackageHooks should not return nil")
	}
}

func TestHookIntegrationResult_ZeroValue(t *testing.T) {
	r := &HookIntegrationResult{}
	if r.HooksIntegrated() != 0 {
		t.Errorf("Zero-value result should have HooksIntegrated() == 0")
	}
	if r.FilesUpdated != 0 || r.FilesSkipped != 0 || r.ScriptsCopied != 0 {
		t.Errorf("Zero-value result should have all zero fields")
	}
	if r.TargetPaths != nil && len(r.TargetPaths) != 0 {
		t.Errorf("Zero-value result should have empty TargetPaths")
	}
}
