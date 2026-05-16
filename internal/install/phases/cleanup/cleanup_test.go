package cleanup

import (
	"sort"
	"testing"
)

func TestDetectStaleFiles_NoStale(t *testing.T) {
	old := []string{"a.txt", "b.txt"}
	new_ := []string{"a.txt", "b.txt"}
	stale := DetectStaleFiles(old, new_)
	if len(stale) != 0 {
		t.Errorf("expected no stale files, got %v", stale)
	}
}

func TestDetectStaleFiles_AllStale(t *testing.T) {
	old := []string{"a.txt", "b.txt"}
	new_ := []string{}
	stale := DetectStaleFiles(old, new_)
	if len(stale) != 2 {
		t.Errorf("expected 2 stale files, got %v", stale)
	}
}

func TestDetectStaleFiles_PartialStale(t *testing.T) {
	old := []string{"a.txt", "b.txt", "c.txt"}
	new_ := []string{"a.txt", "c.txt"}
	stale := DetectStaleFiles(old, new_)
	if len(stale) != 1 || stale[0] != "b.txt" {
		t.Errorf("expected [b.txt], got %v", stale)
	}
}

func TestDetectStaleFiles_NewFilesAdded(t *testing.T) {
	old := []string{"a.txt"}
	new_ := []string{"a.txt", "b.txt"}
	stale := DetectStaleFiles(old, new_)
	if len(stale) != 0 {
		t.Errorf("expected no stale, got %v", stale)
	}
}

func TestCollectOrphanKeys_NoOrphans(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"pkg/a": {"file.txt"},
		},
		IntendedDepKeys: map[string]bool{"pkg/a": true},
	}
	orphans := CollectOrphanKeys(cfg)
	if len(orphans) != 0 {
		t.Errorf("expected no orphans, got %v", orphans)
	}
}

func TestCollectOrphanKeys_WithOrphans(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"pkg/a": {"a.txt"},
			"pkg/b": {"b.txt"},
		},
		IntendedDepKeys: map[string]bool{"pkg/a": true},
	}
	orphans := CollectOrphanKeys(cfg)
	if len(orphans) != 1 || orphans[0] != "pkg/b" {
		t.Errorf("expected [pkg/b], got %v", orphans)
	}
}

func TestCollectOrphanKeys_SkipSelfKey(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"self":  {"self.txt"},
			"pkg/b": {"b.txt"},
		},
		IntendedDepKeys: map[string]bool{},
		SelfKey:         "self",
	}
	orphans := CollectOrphanKeys(cfg)
	sort.Strings(orphans)
	for _, o := range orphans {
		if o == "self" {
			t.Error("self key should be skipped")
		}
	}
}

func TestCollectOrphanKeys_SkipEmptyFiles(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"pkg/a": {},
		},
		IntendedDepKeys: map[string]bool{},
	}
	orphans := CollectOrphanKeys(cfg)
	if len(orphans) != 0 {
		t.Errorf("expected no orphans for empty file list, got %v", orphans)
	}
}

func TestCollectStalePerPackage_NoStale(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{"pkg/a": {"a.txt"}},
		NewDeployedFiles: map[string][]string{"pkg/a": {"a.txt"}},
		PackageErrorCounts: map[string]int{},
	}
	result := CollectStalePerPackage(cfg)
	if len(result) != 0 {
		t.Errorf("expected no stale per-package, got %v", result)
	}
}

func TestCollectStalePerPackage_WithStale(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{"pkg/a": {"a.txt", "b.txt"}},
		NewDeployedFiles: map[string][]string{"pkg/a": {"a.txt"}},
		PackageErrorCounts: map[string]int{},
	}
	result := CollectStalePerPackage(cfg)
	if len(result["pkg/a"]) != 1 || result["pkg/a"][0] != "b.txt" {
		t.Errorf("expected pkg/a stale=[b.txt], got %v", result)
	}
}

func TestCollectStalePerPackage_SkipsErrorPackages(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{"pkg/a": {"a.txt", "b.txt"}},
		NewDeployedFiles: map[string][]string{"pkg/a": {"a.txt"}},
		PackageErrorCounts: map[string]int{"pkg/a": 1},
	}
	result := CollectStalePerPackage(cfg)
	if _, ok := result["pkg/a"]; ok {
		t.Error("package with errors should be skipped")
	}
}

func TestCollectStalePerPackage_NoOldFiles(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{},
		NewDeployedFiles: map[string][]string{"pkg/a": {"a.txt"}},
		PackageErrorCounts: map[string]int{},
	}
	result := CollectStalePerPackage(cfg)
	if len(result) != 0 {
		t.Errorf("expected no results when no old files, got %v", result)
	}
}
