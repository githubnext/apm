package builder

import (
	"strings"
	"testing"
)

func TestBuildDiagnostic_Fields(t *testing.T) {
	d := BuildDiagnostic{Level: "warning", Message: "package skipped"}
	if d.Level != "warning" {
		t.Errorf("Level = %q", d.Level)
	}
	if d.Message != "package skipped" {
		t.Errorf("Message = %q", d.Message)
	}
}

func TestResolvedPackage_Fields(t *testing.T) {
	p := ResolvedPackage{
		Name:       "mypkg",
		SourceRepo: "owner/repo",
		Ref:        "v1.2.3",
		SHA:        "abc1234",
	}
	if p.Name != "mypkg" {
		t.Errorf("Name = %q", p.Name)
	}
	if p.SourceRepo != "owner/repo" {
		t.Errorf("SourceRepo = %q", p.SourceRepo)
	}
}

func TestResolvedPackage_IsPrerelease_False(t *testing.T) {
	p := ResolvedPackage{Ref: "v1.0.0"}
	if p.IsPrerelease {
		t.Error("stable version should not be prerelease")
	}
}

func TestBuildReport_ZeroValue(t *testing.T) {
	r := BuildReport{}
	if r.AddedCount != 0 || r.RemovedCount != 0 || r.UpdatedCount != 0 {
		t.Errorf("zero BuildReport has non-zero counts: %+v", r)
	}
}

func TestBuildReport_WithResolved(t *testing.T) {
	r := BuildReport{
		Resolved: []ResolvedPackage{
			{Name: "pkg1", SourceRepo: "o/r1"},
			{Name: "pkg2", SourceRepo: "o/r2"},
		},
		AddedCount: 2,
	}
	if len(r.Resolved) != 2 {
		t.Errorf("Resolved len = %d", len(r.Resolved))
	}
	if r.AddedCount != 2 {
		t.Errorf("AddedCount = %d", r.AddedCount)
	}
}

func TestBuildOptions_DefaultConcurrency(t *testing.T) {
	o := DefaultBuildOptions()
	if o.Concurrency <= 0 {
		t.Errorf("Concurrency = %d, want > 0", o.Concurrency)
	}
}

func TestBuildError_Message(t *testing.T) {
	e := &BuildError{Msg: "build failed: missing ref"}
	if e.Error() != "build failed: missing ref" {
		t.Errorf("Error() = %q", e.Error())
	}
}

func TestHeadNotAllowedError_Message(t *testing.T) {
	e := &HeadNotAllowedError{Package: "mypkg", Ref: "HEAD"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") {
		t.Errorf("error message should contain pkg name: %q", msg)
	}
}

func TestRefNotFoundError_Message(t *testing.T) {
	e := &RefNotFoundError{Package: "mypkg", Ref: "v99.0.0", OwnerRepo: "owner/repo"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") {
		t.Errorf("error message should contain pkg name: %q", msg)
	}
}

func TestNoMatchingVersionError_Message(t *testing.T) {
	e := &NoMatchingVersionError{Package: "mypkg", VersionRange: "^2.0.0"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") {
		t.Errorf("error message should contain pkg name: %q", msg)
	}
}

func TestResolveResult_Errors(t *testing.T) {
	r := ResolveResult{Errors: [][2]string{{"pkg1", "not found"}}}
	if r.OK() {
		t.Error("result with errors should not be OK")
	}
}

func TestResolveResult_Empty_OK(t *testing.T) {
	r := ResolveResult{}
	if !r.OK() {
		t.Error("empty result should be OK")
	}
}

func TestBuildDiagnostic_VerboseLevel(t *testing.T) {
	d := BuildDiagnostic{Level: "verbose", Message: "resolved to v1.2.3"}
	if d.Level != "verbose" {
		t.Errorf("Level = %q", d.Level)
	}
}

func TestResolvedPackage_Tags(t *testing.T) {
	p := ResolvedPackage{Tags: []string{"v1.0.0", "v1.1.0", "v1.2.0"}}
	if len(p.Tags) != 3 {
		t.Errorf("Tags len = %d", len(p.Tags))
	}
}
