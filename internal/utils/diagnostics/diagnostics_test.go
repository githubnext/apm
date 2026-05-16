package diagnostics_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/diagnostics"
)

func TestNewDiagnosticCollector(t *testing.T) {
	d := diagnostics.New(false)
	if d == nil {
		t.Fatal("New should not return nil")
	}
	if d.HasDiagnostics() {
		t.Error("new collector should have no diagnostics")
	}
	if d.HasErrors() {
		t.Error("new collector should have no errors")
	}
}

func TestDiagnosticCollectorWarn(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("something is off", "pkg/foo", "detail here")
	if !d.HasDiagnostics() {
		t.Error("collector should have diagnostics after Warn")
	}
	if d.HasErrors() {
		t.Error("Warn should not set HasErrors")
	}
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(all))
	}
	if all[0].Category != diagnostics.CategoryWarning {
		t.Errorf("category = %q, want %q", all[0].Category, diagnostics.CategoryWarning)
	}
}

func TestDiagnosticCollectorError(t *testing.T) {
	d := diagnostics.New(false)
	d.Error("fatal issue", "pkg/bar", "")
	if !d.HasErrors() {
		t.Error("Error should set HasErrors")
	}
}

func TestDiagnosticCollectorSecurity(t *testing.T) {
	d := diagnostics.New(false)
	d.Security("malicious content", "pkg/sec", "hash mismatch", "critical")
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(all))
	}
	if all[0].Category != diagnostics.CategorySecurity {
		t.Errorf("category = %q, want %q", all[0].Category, diagnostics.CategorySecurity)
	}
	if all[0].Severity != "critical" {
		t.Errorf("severity = %q, want critical", all[0].Severity)
	}
}

func TestDiagnosticCollectorMultiple(t *testing.T) {
	d := diagnostics.New(false)
	d.Info("info msg", "pkg/a", "")
	d.Warn("warn msg", "pkg/b", "")
	d.Error("error msg", "pkg/c", "")
	all := d.All()
	if len(all) != 3 {
		t.Errorf("expected 3 diagnostics, got %d", len(all))
	}
}

func TestDiagnosticCollectorPolicy(t *testing.T) {
	d := diagnostics.New(false)
	d.Policy("policy violation", "pkg/p", "rule xyz")
	all := d.All()
	if len(all) != 1 || all[0].Category != diagnostics.CategoryPolicy {
		t.Error("Policy diagnostic should have category 'policy'")
	}
}

func TestDiagnosticCollectorAuth(t *testing.T) {
	d := diagnostics.New(false)
	d.Auth("auth failed", "pkg/a", "token expired")
	all := d.All()
	if len(all) != 1 || all[0].Category != diagnostics.CategoryAuth {
		t.Error("Auth diagnostic should have category 'auth'")
	}
}
