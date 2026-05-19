package apmresolver

import (
	"os"
	"testing"
)

func TestParseApmYMLDeps_WhitespaceOnlyLines(t *testing.T) {
	content := "dependencies:\n   \n  - owner/pkg\n   \n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 1 {
		t.Errorf("expected 1 dep, got %d", len(deps))
	}
}

func TestParseApmYMLDeps_CommentAfterDash(t *testing.T) {
	content := "dependencies:\n  - owner/pkg # my comment\n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].RepoURL == "" {
		t.Error("expected non-empty RepoURL after stripping comment")
	}
}

func TestParseApmYMLDeps_QuotedDepRef(t *testing.T) {
	content := "dependencies:\n  - \"owner/repo-q\"\n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d: %v", len(deps), deps)
	}
}

func TestResolveMaxParallel_EnvVarNegative(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "-5")
	p := resolveMaxParallel(0)
	if p <= 0 {
		t.Errorf("negative env var should fall back to default, got %d", p)
	}
}

func TestResolveMaxParallel_EnvVarZero(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "0")
	p := resolveMaxParallel(0)
	if p <= 0 {
		t.Errorf("zero env var should fall back to default, got %d", p)
	}
}

func TestResolveMaxParallel_ExplicitPositive(t *testing.T) {
	p := resolveMaxParallel(7)
	if p != 7 {
		t.Errorf("expected 7, got %d", p)
	}
}

func TestNew_AllOptionsSet(t *testing.T) {
	downloadCalled := false
	fn := func(ref interface{}, a, b, c string) string {
		downloadCalled = true
		return ""
	}
	_ = fn
	r := New(Options{
		MaxDepth:      20,
		ApmModulesDir: "/tmp/apm_modules",
		MaxParallel:   3,
	})
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
	_ = downloadCalled
}

func TestNew_DownloadFnCalledMapsInit(t *testing.T) {
	r := New(Options{MaxDepth: 5})
	if r.downloadedPackages == nil {
		t.Error("downloadedPackages map should be initialized")
	}
	if r.callbackFailures == nil {
		t.Error("callbackFailures map should be initialized")
	}
}

func TestResolveDependencies_NoApmYML(t *testing.T) {
	dir := t.TempDir()
	r := New(Options{})
	graph := r.ResolveDependencies(dir)
	if graph == nil {
		t.Error("expected non-nil dependency graph even without apm.yml")
	}
}

func TestResolveDependencies_EmptyApmYML(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/apm.yml", []byte(""), 0o600)
	r := New(Options{})
	graph := r.ResolveDependencies(dir)
	if graph == nil {
		t.Error("expected non-nil graph for empty apm.yml")
	}
}

func TestParseApmYMLDeps_ManyDeps(t *testing.T) {
	content := "dependencies:\n"
	for i := 0; i < 10; i++ {
		content += "  - owner/pkg" + string(rune('a'+i)) + "\n"
	}
	deps := parseApmYMLDeps(content)
	if len(deps) != 10 {
		t.Errorf("expected 10 deps, got %d", len(deps))
	}
}
