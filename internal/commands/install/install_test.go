package install

import (
	"strings"
	"testing"
)

func TestParseDependencyRefs_simple_name(t *testing.T) {
	entries := parseDependencyRefs([]string{"my-package"})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "my-package" {
		t.Errorf("expected name my-package, got %s", entries[0].Name)
	}
}

func TestParseDependencyRefs_org_repo(t *testing.T) {
	entries := parseDependencyRefs([]string{"myorg/myrepo"})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Org != "myorg" || entries[0].Repo != "myrepo" {
		t.Errorf("unexpected org/repo: %+v", entries[0])
	}
}

func TestParseDependencyRefs_host_org_repo(t *testing.T) {
	entries := parseDependencyRefs([]string{"github.com/myorg/myrepo"})
	if entries[0].Host != "github.com" {
		t.Errorf("expected host github.com, got %s", entries[0].Host)
	}
	if entries[0].Org != "myorg" {
		t.Errorf("expected org myorg, got %s", entries[0].Org)
	}
}

func TestParseDependencyRefs_with_ref(t *testing.T) {
	entries := parseDependencyRefs([]string{"myorg/myrepo@v1.2.3"})
	if entries[0].Ref != "v1.2.3" {
		t.Errorf("expected ref v1.2.3, got %s", entries[0].Ref)
	}
	if entries[0].Repo != "myrepo" {
		t.Errorf("expected repo myrepo, got %s", entries[0].Repo)
	}
}

func TestParseDependencyRefs_multiple(t *testing.T) {
	entries := parseDependencyRefs([]string{"pkg1", "pkg2@main", "org/repo"})
	if len(entries) != 3 {
		t.Errorf("expected 3, got %d", len(entries))
	}
}

func TestMergeDependencies_adds_new(t *testing.T) {
	existing := []DependencyEntry{{Name: "pkg1"}}
	additions := []DependencyEntry{{Name: "pkg2"}}
	result := mergeDependencies(existing, additions)
	if len(result) != 2 {
		t.Errorf("expected 2 after merge, got %d", len(result))
	}
}

func TestMergeDependencies_updates_existing(t *testing.T) {
	existing := []DependencyEntry{{Name: "pkg1", Ref: "v1.0.0"}}
	additions := []DependencyEntry{{Name: "pkg1", Ref: "v2.0.0"}}
	result := mergeDependencies(existing, additions)
	if len(result) != 1 {
		t.Errorf("expected 1, got %d", len(result))
	}
	if result[0].Ref != "v2.0.0" {
		t.Errorf("expected ref updated to v2.0.0, got %s", result[0].Ref)
	}
}

func TestMergeDependencies_empty_existing(t *testing.T) {
	additions := []DependencyEntry{{Name: "pkg1"}}
	result := mergeDependencies(nil, additions)
	if len(result) != 1 {
		t.Errorf("expected 1, got %d", len(result))
	}
}

func TestFormatInstallSummary_installed(t *testing.T) {
	r := &InstallResult{PackagesInstalled: 3, DurationSeconds: 1.5}
	got := FormatInstallSummary(r)
	if !strings.Contains(got, "Installed 3") {
		t.Errorf("expected installed count, got: %s", got)
	}
	if !strings.Contains(got, "[+]") {
		t.Errorf("expected [+] prefix, got: %s", got)
	}
}

func TestFormatInstallSummary_nothing_to_install(t *testing.T) {
	r := &InstallResult{DurationSeconds: 0.1}
	got := FormatInstallSummary(r)
	if !strings.Contains(got, "Nothing to install") {
		t.Errorf("expected 'Nothing to install', got: %s", got)
	}
}

func TestFormatInstallSummary_skipped(t *testing.T) {
	r := &InstallResult{PackagesInstalled: 1, PackagesSkipped: 2, DurationSeconds: 0.5}
	got := FormatInstallSummary(r)
	if !strings.Contains(got, "skipped") {
		t.Errorf("expected 'skipped', got: %s", got)
	}
}

func TestFormatInstallSummary_with_warnings(t *testing.T) {
	r := &InstallResult{PackagesInstalled: 1, Warnings: []string{"some warning"}, DurationSeconds: 1.0}
	got := FormatInstallSummary(r)
	if !strings.Contains(got, "[!]") {
		t.Errorf("expected [!] for warning, got: %s", got)
	}
}

func TestFormatInstallSummary_with_errors(t *testing.T) {
	r := &InstallResult{PackagesInstalled: 0, Errors: []string{"something failed"}, DurationSeconds: 1.0}
	got := FormatInstallSummary(r)
	if !strings.Contains(got, "[x]") {
		t.Errorf("expected [x] for error, got: %s", got)
	}
}
