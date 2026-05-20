package cleanup

import (
	"sort"
	"testing"
)

func TestDetectStaleFiles_NoOldFiles(t *testing.T) {
	stale := DetectStaleFiles(nil, []string{"a.txt", "b.txt"})
	if len(stale) != 0 {
		t.Errorf("expected no stale files, got %v", stale)
	}
}

func TestDetectStaleFiles_NoNewFiles(t *testing.T) {
	stale := DetectStaleFiles([]string{"a.txt", "b.txt"}, nil)
	sort.Strings(stale)
	if len(stale) != 2 || stale[0] != "a.txt" || stale[1] != "b.txt" {
		t.Errorf("all old files should be stale, got %v", stale)
	}
}

func TestDetectStaleFiles_PartialOverlap(t *testing.T) {
	stale := DetectStaleFiles([]string{"a.txt", "b.txt", "c.txt"}, []string{"b.txt", "c.txt"})
	if len(stale) != 1 || stale[0] != "a.txt" {
		t.Errorf("expected [a.txt], got %v", stale)
	}
}

func TestCollectOrphanKeys_Empty(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{},
		IntendedDepKeys:  map[string]bool{},
		SelfKey:          "self",
	}
	keys := CollectOrphanKeys(cfg)
	if len(keys) != 0 {
		t.Errorf("expected no orphan keys, got %v", keys)
	}
}

func TestCollectOrphanKeys_AllIntended(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"pkg-a": {"f1.txt"},
			"pkg-b": {"f2.txt"},
		},
		IntendedDepKeys: map[string]bool{
			"pkg-a": true,
			"pkg-b": true,
		},
		SelfKey: "",
	}
	keys := CollectOrphanKeys(cfg)
	if len(keys) != 0 {
		t.Errorf("no orphans expected, got %v", keys)
	}
}

func TestCollectOrphanKeys_SelfKeySkipped(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"self":  {"selffile"},
			"pkg-x": {"x.txt"},
		},
		IntendedDepKeys: map[string]bool{},
		SelfKey:         "self",
	}
	keys := CollectOrphanKeys(cfg)
	if len(keys) != 1 || keys[0] != "pkg-x" {
		t.Errorf("expected [pkg-x], got %v", keys)
	}
}

func TestCollectStalePerPackage_NoOldFilesExtra2(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles:   map[string][]string{},
		NewDeployedFiles:   map[string][]string{"pkg-a": {"new.txt"}},
		PackageErrorCounts: map[string]int{},
	}
	result := CollectStalePerPackage(cfg)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestCollectStalePerPackage_EmptyOldForPackage(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles:   map[string][]string{"pkg-a": {}},
		NewDeployedFiles:   map[string][]string{"pkg-a": {"new.txt"}},
		PackageErrorCounts: map[string]int{},
	}
	result := CollectStalePerPackage(cfg)
	if len(result) != 0 {
		t.Errorf("expected empty result for empty old files, got %v", result)
	}
}

func TestCleanupResult_Append(t *testing.T) {
	r := CleanupResult{}
	r.Deleted = append(r.Deleted, "d1")
	r.Failed = append(r.Failed, "f1", "f2")
	if len(r.Deleted) != 1 || r.Deleted[0] != "d1" {
		t.Error("Deleted field append failed")
	}
	if len(r.Failed) != 2 {
		t.Error("Failed field append count mismatch")
	}
}
