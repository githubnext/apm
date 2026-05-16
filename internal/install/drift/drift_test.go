package drift_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/drift"
)

// mockDep implements DependencyRef.
type mockDep struct {
	ref      string
	key      string
	insecure bool
}

func (m *mockDep) Reference() string        { return m.ref }
func (m *mockDep) UniqueKey() string         { return m.key }
func (m *mockDep) IsInsecure() bool          { return m.insecure }
func (m *mockDep) Host() string              { return "" }
func (m *mockDep) ArtifactoryPrefix() string { return "" }

// mockLocked implements LockedDep.
type mockLocked struct {
	resolvedRef    string
	resolvedCommit string
	deployedFiles  []string
	insecure       bool
	allowInsecure  bool
	registryPrefix string
	host           string
}

func (m *mockLocked) ResolvedRef() string     { return m.resolvedRef }
func (m *mockLocked) ResolvedCommit() string  { return m.resolvedCommit }
func (m *mockLocked) DeployedFiles() []string { return m.deployedFiles }
func (m *mockLocked) IsInsecure() bool        { return m.insecure }
func (m *mockLocked) AllowInsecure() bool     { return m.allowInsecure }
func (m *mockLocked) RegistryPrefix() string  { return m.registryPrefix }
func (m *mockLocked) Host() string            { return m.host }

// mockLockFile implements LockFile.
type mockLockFile struct {
	deps map[string]drift.LockedDep
}

func (m *mockLockFile) Dependencies() map[string]drift.LockedDep { return m.deps }
func (m *mockLockFile) GetDependency(key string) drift.LockedDep {
	if d, ok := m.deps[key]; ok {
		return d
	}
	return nil
}

func TestDetectRefChange_UpdateRefs(t *testing.T) {
	dep := &mockDep{ref: "v1.0.0"}
	locked := &mockLocked{resolvedRef: "v0.9.0"}
	if drift.DetectRefChange(dep, locked, true) {
		t.Fatal("updateRefs=true should always return false")
	}
}

func TestDetectRefChange_NilLocked(t *testing.T) {
	dep := &mockDep{ref: "v1.0.0"}
	if drift.DetectRefChange(dep, nil, false) {
		t.Fatal("nil lockedDep should return false")
	}
}

func TestDetectRefChange_Changed(t *testing.T) {
	dep := &mockDep{ref: "v2.0.0"}
	locked := &mockLocked{resolvedRef: "v1.0.0"}
	if !drift.DetectRefChange(dep, locked, false) {
		t.Fatal("expected ref change detected")
	}
}

func TestDetectRefChange_Unchanged(t *testing.T) {
	dep := &mockDep{ref: "v1.0.0"}
	locked := &mockLocked{resolvedRef: "v1.0.0"}
	if drift.DetectRefChange(dep, locked, false) {
		t.Fatal("no change expected")
	}
}

func TestDetectRefChange_InsecureToggle(t *testing.T) {
	dep := &mockDep{ref: "main", insecure: true}
	locked := &mockLocked{resolvedRef: "main", insecure: false}
	if !drift.DetectRefChange(dep, locked, false) {
		t.Fatal("insecure toggle should be detected")
	}
}

func TestDetectOrphans_Empty(t *testing.T) {
	lf := &mockLockFile{deps: map[string]drift.LockedDep{}}
	orphans := drift.DetectOrphans(lf, map[string]struct{}{}, nil)
	if len(orphans) != 0 {
		t.Fatalf("expected no orphans, got %v", orphans)
	}
}

func TestDetectOrphans_WithOrphaned(t *testing.T) {
	lf := &mockLockFile{
		deps: map[string]drift.LockedDep{
			"old/dep": &mockLocked{deployedFiles: []string{"file1.md", "file2.md"}},
		},
	}
	intended := map[string]struct{}{} // old/dep NOT in intended
	orphans := drift.DetectOrphans(lf, intended, nil)
	if len(orphans) != 2 {
		t.Fatalf("expected 2 orphaned files, got %d", len(orphans))
	}
}

func TestDetectOrphans_OnlyPackages(t *testing.T) {
	// Partial installs skip orphan detection
	lf := &mockLockFile{
		deps: map[string]drift.LockedDep{
			"old/dep": &mockLocked{deployedFiles: []string{"f.md"}},
		},
	}
	orphans := drift.DetectOrphans(lf, map[string]struct{}{}, []string{"pkg-a"})
	if len(orphans) != 0 {
		t.Fatal("partial install should skip orphan detection")
	}
}

func TestDetectStaleFiles(t *testing.T) {
	old := []string{"a.md", "b.md", "c.md"}
	new_ := []string{"a.md", "c.md"}
	stale := drift.DetectStaleFiles(old, new_)
	if _, ok := stale["b.md"]; !ok {
		t.Fatal("b.md should be stale")
	}
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale file, got %d", len(stale))
	}
}

func TestDetectStaleFiles_NoStale(t *testing.T) {
	old := []string{"a.md"}
	new_ := []string{"a.md", "b.md"}
	stale := drift.DetectStaleFiles(old, new_)
	if len(stale) != 0 {
		t.Fatalf("expected no stale, got %v", stale)
	}
}

func TestDetectConfigDrift_Changed(t *testing.T) {
	current := map[string]interface{}{"srv": map[string]interface{}{"host": "new"}}
	stored := map[string]interface{}{"srv": map[string]interface{}{"host": "old"}}
	drifted := drift.DetectConfigDrift(current, stored)
	if _, ok := drifted["srv"]; !ok {
		t.Fatal("expected srv to be drifted")
	}
}

func TestDetectConfigDrift_NoDrift(t *testing.T) {
	cfg := map[string]interface{}{"srv": "val"}
	drifted := drift.DetectConfigDrift(cfg, cfg)
	if len(drifted) != 0 {
		t.Fatalf("expected no drift, got %v", drifted)
	}
}

func TestDetectConfigDrift_NewEntry(t *testing.T) {
	// New entries (not in stored) should be excluded
	current := map[string]interface{}{"new-srv": "val"}
	stored := map[string]interface{}{}
	drifted := drift.DetectConfigDrift(current, stored)
	if len(drifted) != 0 {
		t.Fatal("new-only entries should not count as drifted")
	}
}
