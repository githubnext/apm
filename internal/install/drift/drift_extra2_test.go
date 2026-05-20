package drift

import (
	"testing"
)

// ---------------------------------------------------------------------------
// RefChangeResult
// ---------------------------------------------------------------------------

func TestRefChangeResult_ZeroValue(t *testing.T) {
	var r RefChangeResult
	if r.Changed {
		t.Error("expected Changed=false in zero value")
	}
}

func TestRefChangeResult_Set(t *testing.T) {
	r := RefChangeResult{Changed: true}
	if !r.Changed {
		t.Error("expected Changed=true")
	}
}

// ---------------------------------------------------------------------------
// DownloadRefOptions
// ---------------------------------------------------------------------------

func TestDownloadRefOptions_ZeroValue(t *testing.T) {
	var o DownloadRefOptions
	if o.UpdateRefs || o.RefChanged {
		t.Error("expected zero value DownloadRefOptions")
	}
}

func TestDownloadRefOptions_Fields(t *testing.T) {
	o := DownloadRefOptions{UpdateRefs: true, RefChanged: true}
	if !o.UpdateRefs {
		t.Error("expected UpdateRefs=true")
	}
	if !o.RefChanged {
		t.Error("expected RefChanged=true")
	}
}

// ---------------------------------------------------------------------------
// SimpleDepRef methods
// ---------------------------------------------------------------------------

func TestSimpleDepRef_ZeroValue(t *testing.T) {
	var s SimpleDepRef
	if s.Reference() != "" || s.UniqueKey() != "" || s.Host() != "" {
		t.Error("expected empty strings in zero value SimpleDepRef")
	}
	if s.IsInsecure() || s.ArtifactoryPrefix() != "" {
		t.Error("expected false/empty in zero value")
	}
}

func TestSimpleDepRef_ArtifactoryPrefix(t *testing.T) {
	s := SimpleDepRef{ArtifactoryPfx: "myregistry"}
	if s.ArtifactoryPrefix() != "myregistry" {
		t.Errorf("expected myregistry, got %q", s.ArtifactoryPrefix())
	}
}

func TestSimpleDepRef_Insecure(t *testing.T) {
	s := SimpleDepRef{Insecure: true}
	if !s.IsInsecure() {
		t.Error("expected IsInsecure=true")
	}
}

// ---------------------------------------------------------------------------
// Internal mockLockFile for DetectOrphans tests
// ---------------------------------------------------------------------------

type mockLockFileInternal struct {
	deps map[string]LockedDep
}

func (m *mockLockFileInternal) Dependencies() map[string]LockedDep { return m.deps }
func (m *mockLockFileInternal) GetDependency(key string) LockedDep { return m.deps[key] }

type mockLockedInternal struct{}

func (m *mockLockedInternal) ResolvedRef() string     { return "v1.0.0" }
func (m *mockLockedInternal) ResolvedCommit() string  { return "" }
func (m *mockLockedInternal) DeployedFiles() []string { return nil }
func (m *mockLockedInternal) IsInsecure() bool        { return false }
func (m *mockLockedInternal) AllowInsecure() bool     { return false }
func (m *mockLockedInternal) RegistryPrefix() string  { return "" }
func (m *mockLockedInternal) Host() string            { return "" }

func TestDetectOrphans_AllPresent(t *testing.T) {
	deps := map[string]LockedDep{"a": &mockLockedInternal{}, "b": &mockLockedInternal{}}
	lf := &mockLockFileInternal{deps: deps}
	intended := map[string]struct{}{"a": {}, "b": {}}
	orphans := DetectOrphans(lf, intended, nil)
	if len(orphans) != 0 {
		t.Errorf("expected no orphans, got %v", orphans)
	}
}

func TestDetectOrphans_AllOrphaned(t *testing.T) {
	// DetectOrphans returns orphaned *files* (from dep.DeployedFiles())
	// When DeployedFiles is empty, no file orphans are returned even if the key is orphaned
	deps := map[string]LockedDep{"a": &mockLockedInternal{}, "b": &mockLockedInternal{}}
	lf := &mockLockFileInternal{deps: deps}
	intended := map[string]struct{}{}
	orphans := DetectOrphans(lf, intended, nil)
	// mockLockedInternal.DeployedFiles() returns nil, so orphaned file set is empty
	if len(orphans) != 0 {
		t.Errorf("expected 0 file orphans (no deployed files), got %d: %v", len(orphans), orphans)
	}
}

type mockLockedWithFiles struct {
	files []string
}

func (m *mockLockedWithFiles) ResolvedRef() string     { return "v1.0.0" }
func (m *mockLockedWithFiles) ResolvedCommit() string  { return "" }
func (m *mockLockedWithFiles) DeployedFiles() []string { return m.files }
func (m *mockLockedWithFiles) IsInsecure() bool        { return false }
func (m *mockLockedWithFiles) AllowInsecure() bool     { return false }
func (m *mockLockedWithFiles) RegistryPrefix() string  { return "" }
func (m *mockLockedWithFiles) Host() string            { return "" }

func TestDetectOrphans_PartialOrphaned(t *testing.T) {
	deps := map[string]LockedDep{
		"a": &mockLockedWithFiles{files: []string{"a/file1.md"}},
		"b": &mockLockedWithFiles{files: []string{"b/file1.md"}},
		"c": &mockLockedWithFiles{files: []string{"c/file1.md"}},
	}
	lf := &mockLockFileInternal{deps: deps}
	intended := map[string]struct{}{"a": {}}
	orphans := DetectOrphans(lf, intended, nil)
	if _, ok := orphans["b/file1.md"]; !ok {
		t.Error("expected b/file1.md to be orphaned")
	}
	if _, ok := orphans["c/file1.md"]; !ok {
		t.Error("expected c/file1.md to be orphaned")
	}
}

func TestDetectOrphans_EmptyLockFile(t *testing.T) {
	lf := &mockLockFileInternal{deps: map[string]LockedDep{}}
	intended := map[string]struct{}{"a": {}}
	orphans := DetectOrphans(lf, intended, nil)
	if len(orphans) != 0 {
		t.Errorf("expected no orphans from empty lockfile, got %v", orphans)
	}
}

// ---------------------------------------------------------------------------
// DetectConfigDrift with nested values
// ---------------------------------------------------------------------------

func TestDetectConfigDrift_NestedMapUnchanged(t *testing.T) {
	current := map[string]interface{}{"key": map[string]interface{}{"inner": "val"}}
	stored := map[string]interface{}{"key": map[string]interface{}{"inner": "val"}}
	drifted := DetectConfigDrift(current, stored)
	if len(drifted) != 0 {
		t.Errorf("expected no drift for identical nested maps, got %v", drifted)
	}
}

func TestDetectConfigDrift_NestedMapChanged(t *testing.T) {
	current := map[string]interface{}{"key": map[string]interface{}{"inner": "new"}}
	stored := map[string]interface{}{"key": map[string]interface{}{"inner": "old"}}
	drifted := DetectConfigDrift(current, stored)
	if _, ok := drifted["key"]; !ok {
		t.Error("expected 'key' to be drifted")
	}
}

func TestDetectConfigDrift_EmptyBoth(t *testing.T) {
	drifted := DetectConfigDrift(map[string]interface{}{}, map[string]interface{}{})
	if len(drifted) != 0 {
		t.Errorf("expected no drift for empty maps, got %v", drifted)
	}
}


