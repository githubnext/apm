package diagnostics_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/diagnostics"
)

func TestDiagnosticCollectorVerboseMode(t *testing.T) {
	d := diagnostics.New(true)
	if d == nil {
		t.Fatal("New(true) should not return nil")
	}
}

func TestDiagnosticCollectorSkip(t *testing.T) {
	d := diagnostics.New(false)
	d.Skip("/some/path", "pkg/skip")
	// Skip records a diagnostic but not an error
	if d.HasErrors() {
		t.Error("Skip should not count as an error")
	}
}

func TestDiagnosticCollectorOverwrite(t *testing.T) {
	d := diagnostics.New(false)
	d.Overwrite("/some/path", "pkg/overwrite", "forced overwrite")
	if d.HasErrors() {
		t.Error("Overwrite should not count as an error")
	}
}

func TestDiagnosticCollectorDrift(t *testing.T) {
	d := diagnostics.New(false)
	d.Drift("drift detected", "pkg/drift", "file changed")
	all := d.All()
	if len(all) == 0 {
		t.Error("Drift should add a diagnostic")
	}
}

func TestDiagnosticCollector_AllOrdering(t *testing.T) {
	d := diagnostics.New(false)
	// Add multiple types and verify All() returns them all
	d.Warn("w1", "p1", "")
	d.Error("e1", "p2", "")
	d.Info("i1", "p3", "")
	d.Security("s1", "p4", "", "high")
	d.Policy("pol1", "p5", "")
	d.Auth("a1", "p6", "")
	all := d.All()
	if len(all) != 6 {
		t.Errorf("expected 6 diagnostics, got %d", len(all))
	}
}

func TestDiagnosticCollector_CategoryValues(t *testing.T) {
	if diagnostics.CategoryWarning == "" {
		t.Error("CategoryWarning should not be empty")
	}
	if diagnostics.CategoryError == "" {
		t.Error("CategoryError should not be empty")
	}
	if diagnostics.CategorySecurity == "" {
		t.Error("CategorySecurity should not be empty")
	}
	if diagnostics.CategoryPolicy == "" {
		t.Error("CategoryPolicy should not be empty")
	}
	if diagnostics.CategoryAuth == "" {
		t.Error("CategoryAuth should not be empty")
	}
	if diagnostics.CategoryInfo == "" {
		t.Error("CategoryInfo should not be empty")
	}
}

func TestDiagnosticCollector_ErrorDoesNotAddMultiple(t *testing.T) {
	d := diagnostics.New(false)
	d.Error("err1", "pkg", "detail")
	d.Error("err2", "pkg2", "")
	all := d.All()
	if len(all) != 2 {
		t.Errorf("expected 2 diagnostics, got %d", len(all))
	}
}

func TestDiagnosticCollector_SecurityWithSeverity(t *testing.T) {
	for _, sev := range []string{"low", "medium", "high", "critical"} {
		d := diagnostics.New(false)
		d.Security("vuln", "pkg", "detail", sev)
		all := d.All()
		if len(all) != 1 {
			t.Errorf("sev=%s: expected 1, got %d", sev, len(all))
		}
		if all[0].Severity != sev {
			t.Errorf("sev=%s: got %q", sev, all[0].Severity)
		}
	}
}

func TestDiagnosticCollector_RenderSummaryNoOp(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("test warning", "pkg", "detail")
	// RenderSummary should not panic (it writes to stdout)
	d.RenderSummary()
}

func TestDiagnosticCollector_HasErrorsAfterMultiple(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("w", "p", "")
	if d.HasErrors() {
		t.Error("Warn should not set HasErrors")
	}
	d.Error("e", "p", "")
	if !d.HasErrors() {
		t.Error("Error should set HasErrors")
	}
}
