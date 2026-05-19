package dryrun_test

import (
	"fmt"
	"testing"

	"github.com/githubnext/apm/internal/install/presentation/dryrun"
)

type extraMockLogger struct {
	progress   []string
	dryRun     []string
	successMsg []string
}

func (m *extraMockLogger) Progress(msg string)     { m.progress = append(m.progress, msg) }
func (m *extraMockLogger) DryRunNotice(msg string) { m.dryRun = append(m.dryRun, msg) }
func (m *extraMockLogger) Success(msg string)      { m.successMsg = append(m.successMsg, msg) }

type extraMockDep struct {
	url string
	ref string
	key string
}

func (d *extraMockDep) RepoURL() string      { return d.url }
func (d *extraMockDep) Reference() string    { return d.ref }
func (d *extraMockDep) GetUniqueKey() string { return d.key }

func TestRenderAndExit_DevAPMDeps(t *testing.T) {
	lg := &extraMockLogger{}
	deps := []dryrun.Dep{&extraMockDep{url: "https://github.com/dev/pkg", ref: "main"}}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallAPM: true,
		DevAPMDeps:       deps,
	})
	// Should emit dry-run notice about integration not previewed
	if len(lg.dryRun) == 0 {
		t.Error("expected DryRunNotice for DevAPMDeps")
	}
}

func TestRenderAndExit_ManyOrphans(t *testing.T) {
	lg := &extraMockLogger{}
	orphans := make([]string, 15)
	for i := range orphans {
		orphans[i] = fmt.Sprintf("file%d.txt", i)
	}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:          lg,
		LockfileOrphans: orphans,
	})
	// Check "... and N more" message
	foundMore := false
	for _, msg := range lg.progress {
		if msg == "  ... and 5 more" {
			foundMore = true
		}
	}
	if !foundMore {
		t.Errorf("expected '... and 5 more' for 15 orphans (limit 10); got: %v", lg.progress)
	}
}

func TestRenderAndExit_NilOrphans(t *testing.T) {
	lg := &extraMockLogger{}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:          lg,
		LockfileOrphans: nil,
	})
	// nil orphans should not produce orphan output
	for _, msg := range lg.progress {
		if len(msg) > 0 && msg[:2] == "  " {
			t.Errorf("unexpected indented message for nil orphans: %q", msg)
		}
	}
}

func TestRenderAndExit_BothAPMAndMCP(t *testing.T) {
	lg := &extraMockLogger{}
	apmDeps := []dryrun.Dep{&extraMockDep{url: "https://github.com/a/b", ref: "v1.0"}}
	mcpDeps := []fmt.Stringer{stringerDep{"mcp-tool"}}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallAPM: true,
		APMDeps:          apmDeps,
		ShouldInstallMCP: true,
		MCPDeps:          mcpDeps,
	})
	// success message should still be emitted
	if len(lg.successMsg) == 0 {
		t.Error("expected success message with both APM and MCP deps")
	}
}

func TestRenderAndExit_EmptyOrphansSlice(t *testing.T) {
	lg := &extraMockLogger{}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:          lg,
		LockfileOrphans: []string{},
	})
	// empty (non-nil) slice should not produce orphan output
	for _, msg := range lg.progress {
		if msg[:1] == " " {
			t.Errorf("unexpected indented message for empty orphans: %q", msg)
		}
	}
}

func TestRenderAndExit_OnlyMCPDeps(t *testing.T) {
	lg := &extraMockLogger{}
	mcpDeps := []fmt.Stringer{
		stringerDep{"tool-a"},
		stringerDep{"tool-b"},
	}
	dryrun.RenderAndExit(dryrun.Options{
		Logger:           lg,
		ShouldInstallMCP: true,
		MCPDeps:          mcpDeps,
	})
	count := 0
	for _, msg := range lg.progress {
		if len(msg) > 4 && msg[:4] == "  - " {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 MCP dep lines, got %d; progress: %v", count, lg.progress)
	}
}

type stringerDep struct{ s string }

func (sd stringerDep) String() string { return sd.s }
