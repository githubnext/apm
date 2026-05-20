package dryrun_test

import (
	"fmt"
	"testing"

	"github.com/githubnext/apm/internal/install/presentation/dryrun"
)

type mockDryLogger2 struct {
	progress []string
	dryRun   []string
	success  []string
}

func (m *mockDryLogger2) Progress(msg string)     { m.progress = append(m.progress, msg) }
func (m *mockDryLogger2) DryRunNotice(msg string) { m.dryRun = append(m.dryRun, msg) }
func (m *mockDryLogger2) Success(msg string)      { m.success = append(m.success, msg) }

type mockDryDep2 struct {
	url string
	ref string
	key string
}

func (d *mockDryDep2) RepoURL() string      { return d.url }
func (d *mockDryDep2) Reference() string    { return d.ref }
func (d *mockDryDep2) GetUniqueKey() string { return d.key }

type mockStringer2 struct{ val string }

func (s *mockStringer2) String() string { return s.val }

func TestRenderAndExit_NoDepsExtra2(t *testing.T) {
	log := &mockDryLogger2{}
	opts := dryrun.Options{Logger: log}
	dryrun.RenderAndExit(opts)
	if len(log.success) == 0 {
		t.Error("expected at least one success message")
	}
	found := false
	for _, msg := range log.progress {
		if msg == "No dependencies found in apm.yml" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'No dependencies found' message when no deps set")
	}
}

func TestRenderAndExit_APMDepsUpdateMode(t *testing.T) {
	log := &mockDryLogger2{}
	opts := dryrun.Options{
		Logger:           log,
		ShouldInstallAPM: true,
		APMDeps:          []dryrun.Dep{&mockDryDep2{url: "https://github.com/a/b", ref: "main", key: "k"}},
		Update:           true,
	}
	dryrun.RenderAndExit(opts)
	found := false
	for _, msg := range log.progress {
		if len(msg) > 0 && containsStr2(msg, "update") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'update' in progress messages when Update=true")
	}
}

func containsStr2(s, sub string) bool {
	return len(s) >= len(sub) && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}()
}

func TestRenderAndExit_APMDepNoRef(t *testing.T) {
	log := &mockDryLogger2{}
	opts := dryrun.Options{
		Logger:           log,
		ShouldInstallAPM: true,
		APMDeps:          []dryrun.Dep{&mockDryDep2{url: "https://github.com/a/b", ref: "", key: "k"}},
	}
	dryrun.RenderAndExit(opts)
	found := false
	for _, msg := range log.progress {
		if containsStr2(msg, "main") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'main' as default ref when Reference() returns empty")
	}
}

func TestRenderAndExit_OrphansExactly10(t *testing.T) {
	log := &mockDryLogger2{}
	orphans := make([]string, 10)
	for i := range orphans {
		orphans[i] = fmt.Sprintf("file%d.txt", i)
	}
	opts := dryrun.Options{
		Logger:          log,
		LockfileOrphans: orphans,
	}
	dryrun.RenderAndExit(opts)
	foundOrphanHeader := false
	for _, msg := range log.progress {
		if containsStr2(msg, "Files that would be removed") {
			foundOrphanHeader = true
		}
	}
	if !foundOrphanHeader {
		t.Error("expected orphan header message")
	}
}

func TestOptions_ZeroValue(t *testing.T) {
	var opts dryrun.Options
	if opts.Logger != nil || opts.ShouldInstallAPM || opts.Update {
		t.Error("Options zero value should have nil/false fields")
	}
}
