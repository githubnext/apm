package drift_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/drift"
)

func TestDetectStaleFiles_NoChange(t *testing.T) {
	files := []string{"a.txt", "b.txt"}
	stale := drift.DetectStaleFiles(files, files)
	if len(stale) != 0 {
		t.Errorf("expected no stale files when sets are identical, got %v", stale)
	}
}

func TestDetectStaleFiles_AllRemoved(t *testing.T) {
	old := []string{"a.txt", "b.txt", "c.txt"}
	stale := drift.DetectStaleFiles(old, []string{})
	if len(stale) != 3 {
		t.Errorf("expected 3 stale files, got %d: %v", len(stale), stale)
	}
}

func TestDetectStaleFiles_PartialRemoval(t *testing.T) {
	old := []string{"a.txt", "b.txt"}
	newFiles := []string{"b.txt"}
	stale := drift.DetectStaleFiles(old, newFiles)
	if _, ok := stale["a.txt"]; !ok {
		t.Error("a.txt should be stale")
	}
	if _, ok := stale["b.txt"]; ok {
		t.Error("b.txt should not be stale")
	}
}

func TestDetectStaleFiles_NewFilesAdded(t *testing.T) {
	old := []string{"a.txt"}
	newFiles := []string{"a.txt", "b.txt"}
	stale := drift.DetectStaleFiles(old, newFiles)
	if len(stale) != 0 {
		t.Errorf("expected no stale files when only adding, got %v", stale)
	}
}

func TestDetectStaleFiles_EmptyBoth(t *testing.T) {
	stale := drift.DetectStaleFiles([]string{}, []string{})
	if len(stale) != 0 {
		t.Errorf("expected empty stale set, got %v", stale)
	}
}

func TestDetectConfigDrift_IdenticalConfigs(t *testing.T) {
	cfg := map[string]interface{}{
		"key": "value",
	}
	drifted := drift.DetectConfigDrift(cfg, cfg)
	if len(drifted) != 0 {
		t.Errorf("expected no drifted keys for identical configs, got %v", drifted)
	}
}

func TestDetectConfigDrift_ValueChanged(t *testing.T) {
	current := map[string]interface{}{"key": "new"}
	stored := map[string]interface{}{"key": "old"}
	drifted := drift.DetectConfigDrift(current, stored)
	if _, ok := drifted["key"]; !ok {
		t.Error("expected 'key' in drifted set")
	}
}

func TestDetectConfigDrift_KeyMissing(t *testing.T) {
	// DetectConfigDrift only iterates currentConfigs keys;
	// a key present only in stored (but not in current) is not flagged as drift.
	current := map[string]interface{}{}
	stored := map[string]interface{}{"key": "val"}
	drifted := drift.DetectConfigDrift(current, stored)
	// No keys in current to compare; result should be empty.
	if len(drifted) != 0 {
		t.Errorf("expected no drifted keys, got %v", drifted)
	}
}

func TestDetectConfigDrift_NewKeyInCurrent(t *testing.T) {
	current := map[string]interface{}{"key": "val", "newkey": "x"}
	stored := map[string]interface{}{"key": "val"}
	drifted := drift.DetectConfigDrift(current, stored)
	// newkey only in current is not drift (it was added, not changed)
	if _, ok := drifted["key"]; ok {
		t.Error("key with same value should not be drifted")
	}
}

func TestSimpleDepRef_Fields(t *testing.T) {
	dep := &drift.SimpleDepRef{
		Ref:            "main",
		Key:            "mykey",
		Insecure:       true,
		HostVal:        "github.com",
		ArtifactoryPfx: "prefix",
	}
	if dep.Reference() != "main" {
		t.Errorf("Reference() = %q, want main", dep.Reference())
	}
	if dep.UniqueKey() != "mykey" {
		t.Errorf("UniqueKey() = %q, want mykey", dep.UniqueKey())
	}
	if !dep.IsInsecure() {
		t.Error("IsInsecure() should be true")
	}
	if dep.Host() != "github.com" {
		t.Errorf("Host() = %q, want github.com", dep.Host())
	}
	if dep.ArtifactoryPrefix() != "prefix" {
		t.Errorf("ArtifactoryPrefix() = %q, want prefix", dep.ArtifactoryPrefix())
	}
}
