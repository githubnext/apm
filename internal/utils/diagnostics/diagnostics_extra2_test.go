package diagnostics_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/diagnostics"
)

func TestDiagnosticCollector_InfoDiagnostic(t *testing.T) {
	d := diagnostics.New(false)
	d.Info("info msg", "pkg", "detail")
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(all))
	}
	if all[0].Category != "info" {
		t.Errorf("expected category=info, got %q", all[0].Category)
	}
}

func TestDiagnosticCollector_SkipAdded(t *testing.T) {
	d := diagnostics.New(false)
	d.Skip("/path/to/file", "some/pkg")
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 diagnostic after Skip, got %d", len(all))
	}
}

func TestDiagnosticCollector_OverwriteCategory(t *testing.T) {
	d := diagnostics.New(false)
	d.Overwrite("/path", "pkg", "detail")
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1, got %d", len(all))
	}
	if all[0].Category != "overwrite" {
		t.Errorf("expected category=overwrite, got %q", all[0].Category)
	}
}

func TestDiagnosticCollector_PolicyCategory(t *testing.T) {
	d := diagnostics.New(false)
	d.Policy("policy violation", "pkg", "detail")
	all := d.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(all))
	}
	if all[0].Category != "policy" {
		t.Errorf("expected category=policy, got %q", all[0].Category)
	}
}

func TestDiagnosticCollector_AuthCategory(t *testing.T) {
	d := diagnostics.New(false)
	d.Auth("auth failure", "pkg", "detail")
	all := d.All()
	if len(all) != 1 || all[0].Category != "auth" {
		t.Errorf("expected auth category, got %v", all)
	}
}

func TestDiagnosticCollector_DriftCategory(t *testing.T) {
	d := diagnostics.New(false)
	d.Drift("drift detected", "pkg", "detail")
	all := d.All()
	if len(all) != 1 || all[0].Category != "drift" {
		t.Errorf("expected drift category, got %v", all)
	}
}

func TestDiagnosticCollector_HasDiagnostics_EmptyFalse(t *testing.T) {
	d := diagnostics.New(false)
	if d.HasDiagnostics() {
		t.Error("expected HasDiagnostics=false for empty collector")
	}
}

func TestDiagnosticCollector_HasDiagnostics_TrueAfterAdd(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("w", "p", "d")
	if !d.HasDiagnostics() {
		t.Error("expected HasDiagnostics=true after adding warning")
	}
}

func TestDiagnosticCollector_HasErrors_False_ForWarn(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("w", "p", "d")
	if d.HasErrors() {
		t.Error("expected HasErrors=false for warning-only collector")
	}
}

func TestDiagnosticCollector_HasErrors_True(t *testing.T) {
	d := diagnostics.New(false)
	d.Error("e", "p", "d")
	if !d.HasErrors() {
		t.Error("expected HasErrors=true after adding error")
	}
}

func TestDiagnosticCollector_MessageField(t *testing.T) {
	d := diagnostics.New(false)
	d.Warn("my message", "my/pkg", "my detail")
	all := d.All()
	if all[0].Message != "my message" {
		t.Errorf("expected message='my message', got %q", all[0].Message)
	}
	if all[0].Package != "my/pkg" {
		t.Errorf("expected pkg='my/pkg', got %q", all[0].Package)
	}
	if all[0].Detail != "my detail" {
		t.Errorf("expected detail='my detail', got %q", all[0].Detail)
	}
}

func TestDiagnosticCollector_SecuritySeverityField(t *testing.T) {
	d := diagnostics.New(false)
	d.Security("hidden char", "pkg", "detail", "critical")
	all := d.All()
	if all[0].Severity != "critical" {
		t.Errorf("expected severity=critical, got %q", all[0].Severity)
	}
}

func TestDiagnosticCollector_MultipleTypes(t *testing.T) {
	d := diagnostics.New(false)
	d.Info("i", "p", "")
	d.Warn("w", "p", "")
	d.Error("e", "p", "")
	d.Policy("pol", "p", "")
	if len(d.All()) != 4 {
		t.Errorf("expected 4 diagnostics, got %d", len(d.All()))
	}
}

func TestDiagnosticCollector_RenderSummaryWithDiagnostics(t *testing.T) {
	d := diagnostics.New(false)
	d.Error("error msg", "pkg", "detail")
	d.Warn("warn msg", "pkg", "detail")
	d.RenderSummary()
}

func TestDiagnosticCollector_VerboseWithInfo(t *testing.T) {
	d := diagnostics.New(true)
	d.Info("verbose info", "pkg", "")
	if len(d.All()) != 1 {
		t.Errorf("expected 1 diagnostic, got %d", len(d.All()))
	}
}
