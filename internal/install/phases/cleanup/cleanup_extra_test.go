package cleanup

import (
	"sort"
	"testing"
)

func TestCleanupResult_ZeroValue(t *testing.T) {
	var r CleanupResult
	if len(r.Deleted) != 0 {
		t.Error("Deleted should be nil/empty")
	}
	if len(r.DeletedTargets) != 0 {
		t.Error("DeletedTargets should be nil/empty")
	}
	if len(r.Failed) != 0 {
		t.Error("Failed should be nil/empty")
	}
	if len(r.SkippedUserEdit) != 0 {
		t.Error("SkippedUserEdit should be nil/empty")
	}
}

func TestCleanupResult_Fields(t *testing.T) {
	r := CleanupResult{
		Deleted:         []string{"a.txt", "b.txt"},
		DeletedTargets:  []string{"t.txt"},
		Failed:          []string{"f.txt"},
		SkippedUserEdit: []string{"s.txt"},
	}
	if len(r.Deleted) != 2 {
		t.Errorf("expected 2 deleted, got %d", len(r.Deleted))
	}
	if len(r.DeletedTargets) != 1 {
		t.Errorf("expected 1 deleted target, got %d", len(r.DeletedTargets))
	}
}

func TestDetectStaleFiles_DuplicatesInOld(t *testing.T) {
	// DetectStaleFiles uses a set for newFiles; duplicates in old are counted once per value
	old := []string{"a.txt", "a.txt", "b.txt"}
	new_ := []string{"a.txt"}
	stale := DetectStaleFiles(old, new_)
	// a.txt appears twice in old; the set check means each occurrence is evaluated
	// Since newSet["a.txt"]=true, both "a.txt" entries survive; only b.txt is stale
	if len(stale) != 1 || stale[0] != "b.txt" {
		t.Errorf("expected 1 stale entry [b.txt], got %v", stale)
	}
}

func TestDetectStaleFiles_EmptyBoth(t *testing.T) {
	stale := DetectStaleFiles(nil, nil)
	if len(stale) != 0 {
		t.Errorf("expected empty stale for empty inputs, got %v", stale)
	}
}

func TestCollectOrphanKeys_SkipSelfKey_WithFiles(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"self": {"file.txt"},
			"pkg":  {"other.txt"},
		},
		IntendedDepKeys: map[string]bool{},
		SelfKey:         "self",
	}
	orphans := CollectOrphanKeys(cfg)
	if len(orphans) != 1 || orphans[0] != "pkg" {
		t.Errorf("expected [pkg], got %v", orphans)
	}
}

func TestCollectOrphanKeys_MultipleOrphans(t *testing.T) {
	cfg := OrphanCleanupConfig{
		ExistingLockDeps: map[string][]string{
			"a": {"a.txt"},
			"b": {"b.txt"},
			"c": {"c.txt"},
		},
		IntendedDepKeys: map[string]bool{"a": true},
	}
	orphans := CollectOrphanKeys(cfg)
	sort.Strings(orphans)
	if len(orphans) != 2 {
		t.Errorf("expected 2 orphans, got %v", orphans)
	}
	if orphans[0] != "b" || orphans[1] != "c" {
		t.Errorf("unexpected orphans: %v", orphans)
	}
}

func TestCollectStalePerPackage_MultiplePackages(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{
			"pkg-a": {"a1.txt", "a2.txt"},
			"pkg-b": {"b1.txt"},
		},
		NewDeployedFiles: map[string][]string{
			"pkg-a": {"a1.txt"},
			"pkg-b": {"b1.txt", "b2.txt"},
		},
		PackageErrorCounts: map[string]int{},
	}
	stale := CollectStalePerPackage(cfg)
	if len(stale["pkg-a"]) != 1 || stale["pkg-a"][0] != "a2.txt" {
		t.Errorf("pkg-a stale: expected [a2.txt], got %v", stale["pkg-a"])
	}
	if len(stale["pkg-b"]) != 0 {
		t.Errorf("pkg-b stale: expected none, got %v", stale["pkg-b"])
	}
}

func TestCollectStalePerPackage_AllFilesRemoved(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{
			"pkg": {"x.txt", "y.txt"},
		},
		NewDeployedFiles: map[string][]string{
			"pkg": {},
		},
		PackageErrorCounts: map[string]int{},
	}
	stale := CollectStalePerPackage(cfg)
	sort.Strings(stale["pkg"])
	if len(stale["pkg"]) != 2 {
		t.Errorf("expected 2 stale files, got %v", stale["pkg"])
	}
}

func TestCollectStalePerPackage_ErrorPackageSkipped(t *testing.T) {
	cfg := StaleCleanupConfig{
		OldDeployedFiles: map[string][]string{
			"bad-pkg": {"bad.txt"},
		},
		NewDeployedFiles: map[string][]string{
			"bad-pkg": {},
		},
		PackageErrorCounts: map[string]int{
			"bad-pkg": 1,
		},
	}
	stale := CollectStalePerPackage(cfg)
	if _, ok := stale["bad-pkg"]; ok {
		t.Error("error package should be skipped")
	}
}

func TestOrphanCleanupConfig_ZeroValue(t *testing.T) {
	var cfg OrphanCleanupConfig
	orphans := CollectOrphanKeys(cfg)
	if len(orphans) != 0 {
		t.Errorf("zero-value config should yield no orphans, got %v", orphans)
	}
}

func TestStaleCleanupConfig_ZeroValue(t *testing.T) {
	var cfg StaleCleanupConfig
	stale := CollectStalePerPackage(cfg)
	if len(stale) != 0 {
		t.Errorf("zero-value config should yield no stale, got %v", stale)
	}
}
