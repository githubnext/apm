package dryrun_test

import (
	"fmt"
	"testing"

	"github.com/githubnext/apm/internal/install/presentation/dryrun"
)

// mockLogger captures calls to the Logger interface.
type mockLogger struct {
	progress   []string
	dryRun     []string
	successMsg []string
}

func (m *mockLogger) Progress(msg string)      { m.progress = append(m.progress, msg) }
func (m *mockLogger) DryRunNotice(msg string)  { m.dryRun = append(m.dryRun, msg) }
func (m *mockLogger) Success(msg string)       { m.successMsg = append(m.successMsg, msg) }

// mockDep implements the Dep interface.
type mockDep struct {
	repoURL   string
	reference string
	key       string
}

func (d *mockDep) RepoURL() string      { return d.repoURL }
func (d *mockDep) Reference() string    { return d.reference }
func (d *mockDep) GetUniqueKey() string { return d.key }

// mockMCPDep implements fmt.Stringer.
type mockMCPDep struct{ name string }

func (m *mockMCPDep) String() string { return m.name }

func TestRenderAndExit_NoDeps(t *testing.T) {
	lg := &mockLogger{}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallAPM: false,
		ShouldInstallMCP: false,
	})
	found := false
	for _, msg := range lg.progress {
		if msg == "No dependencies found in apm.yml" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'No dependencies found' message; got: %v", lg.progress)
	}
}

func TestRenderAndExit_WithAPMDeps(t *testing.T) {
	lg := &mockLogger{}
	deps := []dryrun.Dep{
		&mockDep{repoURL: "https://github.com/a/b", reference: "v1.0"},
		&mockDep{repoURL: "https://github.com/c/d", reference: ""},
	}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallAPM: true,
		APMDeps:          deps,
	})
	// Should mention 2 APM deps
	foundHeader := false
	for _, msg := range lg.progress {
		if msg == fmt.Sprintf("APM dependencies (%d):", len(deps)) {
			foundHeader = true
		}
	}
	if !foundHeader {
		t.Errorf("expected APM dep count header; got: %v", lg.progress)
	}
	// Second dep should default ref to "main"
	foundMain := false
	for _, msg := range lg.progress {
		if msg == "  - https://github.com/c/d#main -> install" {
			foundMain = true
		}
	}
	if !foundMain {
		t.Errorf("expected default ref=main; got: %v", lg.progress)
	}
}

func TestRenderAndExit_UpdateAction(t *testing.T) {
	lg := &mockLogger{}
	deps := []dryrun.Dep{&mockDep{repoURL: "https://github.com/a/b", reference: "main"}}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallAPM: true,
		APMDeps:          deps,
		Update:           true,
	})
	found := false
	for _, msg := range lg.progress {
		if msg == "  - https://github.com/a/b#main -> update" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'update' action; got: %v", lg.progress)
	}
}

func TestRenderAndExit_WithMCPDeps(t *testing.T) {
	lg := &mockLogger{}
	mcpDeps := []fmt.Stringer{&mockMCPDep{"mcp-tool"}}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallMCP: true,
		MCPDeps:          mcpDeps,
	})
	found := false
	for _, msg := range lg.progress {
		if msg == "  - mcp-tool" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected mcp dep listed; got: %v", lg.progress)
	}
}

func TestRenderAndExit_WithOrphans(t *testing.T) {
	lg := &mockLogger{}
	orphans := []string{"file1.txt", "file2.txt"}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:          lg,
		LockfileOrphans: orphans,
	})
	found := false
	for _, msg := range lg.progress {
		if msg == fmt.Sprintf("Files that would be removed (packages no longer in apm.yml): %d", len(orphans)) {
			found = true
		}
	}
	if !found {
		t.Errorf("expected orphan count message; got: %v", lg.progress)
	}
}

func TestRenderAndExit_SuccessMessage(t *testing.T) {
	lg := &mockLogger{}
	dryrun.RenderAndExit(dryrun.Options{Logger: lg})
	if len(lg.successMsg) == 0 {
		t.Error("expected at least one success message")
	}
}
