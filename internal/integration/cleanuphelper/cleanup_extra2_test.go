package cleanuphelper

import (
	"testing"
)

func TestValidateDeployPath_AbsoluteRejected(t *testing.T) {
	ok := ValidateDeployPath("/absolute/path", "/root", []string{".github/"})
	if ok {
		t.Error("expected absolute path to be rejected")
	}
}

func TestValidateDeployPath_DotDotComponent(t *testing.T) {
	ok := ValidateDeployPath("safe/../evil", "/root", []string{"safe/"})
	if ok {
		t.Error("expected dotdot path to be rejected")
	}
}

func TestValidateDeployPath_CoworkURI(t *testing.T) {
	ok := ValidateDeployPath("cowork://some/path", "/root", []string{".github/"})
	if ok {
		t.Error("expected cowork:// URI to be rejected")
	}
}

func TestValidateDeployPath_ValidPrefix(t *testing.T) {
	ok := ValidateDeployPath(".github/skills/foo.md", "/root", []string{".github/"})
	if !ok {
		t.Error("expected valid prefix to pass")
	}
}

func TestValidateDeployPath_NoMatchingPrefix(t *testing.T) {
	ok := ValidateDeployPath("docs/readme.md", "/root", []string{".github/"})
	if ok {
		t.Error("expected no matching prefix to fail")
	}
}

func TestDiagnosticCollector_Warn_Appends(t *testing.T) {
	d := &DiagnosticCollector{}
	d.Warn("pkg1", "warning message 1")
	d.Warn("pkg2", "warning message 2")
	if len(d.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(d.Warnings))
	}
	if d.Warnings[0].Package != "pkg1" {
		t.Errorf("unexpected Package %q", d.Warnings[0].Package)
	}
	if d.Warnings[1].Message != "warning message 2" {
		t.Errorf("unexpected Message %q", d.Warnings[1].Message)
	}
}

func TestDiagnosticCollector_InitiallyEmpty(t *testing.T) {
	d := &DiagnosticCollector{}
	if len(d.Warnings) != 0 {
		t.Errorf("expected empty warnings, got %d", len(d.Warnings))
	}
}

func TestDiagnostic_PackageField(t *testing.T) {
	diag := Diagnostic{Package: "mypkg", Message: "oops"}
	if diag.Package != "mypkg" {
		t.Errorf("unexpected Package %q", diag.Package)
	}
}

func TestDiagnostic_MessageField(t *testing.T) {
	diag := Diagnostic{Message: "something went wrong"}
	if diag.Message != "something went wrong" {
		t.Errorf("unexpected Message %q", diag.Message)
	}
}

func TestOptions_ZeroValue(t *testing.T) {
	var o Options
	if o.ProjectRoot != "" {
		t.Errorf("expected empty ProjectRoot")
	}
	if o.DepKey != "" {
		t.Errorf("expected empty DepKey")
	}
	if o.FailedPathRetained {
		t.Error("expected FailedPathRetained false")
	}
}

func TestOptions_FieldRoundtrip(t *testing.T) {
	o := Options{
		DepKey:      "org/repo",
		ProjectRoot: "/project",
		IntegrationPrefixes: []string{".github/"},
		FailedPathRetained:  true,
	}
	if o.DepKey != "org/repo" {
		t.Errorf("unexpected DepKey %q", o.DepKey)
	}
	if o.ProjectRoot != "/project" {
		t.Errorf("unexpected ProjectRoot %q", o.ProjectRoot)
	}
	if len(o.IntegrationPrefixes) != 1 {
		t.Errorf("expected 1 prefix, got %d", len(o.IntegrationPrefixes))
	}
	if !o.FailedPathRetained {
		t.Error("expected FailedPathRetained true")
	}
}
