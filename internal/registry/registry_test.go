package registry_test

import (
	"testing"

	"github.com/githubnext/apm/internal/registry"
)

// TestParityDefaultRegistryURL verifies the default registry URL.
func TestParityDefaultRegistryURL(t *testing.T) {
	if registry.DefaultRegistryURL != "https://api.mcp.github.com" {
		t.Errorf("unexpected default registry URL: %s", registry.DefaultRegistryURL)
	}
}

// TestParityRegistryV01Prefix verifies the API version prefix.
func TestParityRegistryV01Prefix(t *testing.T) {
	if registry.V01Prefix != "/v0.1" {
		t.Errorf("unexpected v0.1 prefix: %s", registry.V01Prefix)
	}
}

// TestParityServerNotFoundError verifies error message format.
func TestParityServerNotFoundError(t *testing.T) {
	err := &registry.ServerNotFoundError{
		ServerName:  "my-server",
		RegistryURL: "https://api.mcp.github.com",
	}
	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
	if !containsStr(msg, "my-server") {
		t.Errorf("expected server name in error: %s", msg)
	}
	if !containsStr(msg, "https://api.mcp.github.com") {
		t.Errorf("expected registry URL in error: %s", msg)
	}
	if !containsStr(msg, "v0.1") {
		t.Errorf("expected v0.1 hint in error: %s", msg)
	}
}

// TestParityRegistryError verifies error with status code.
func TestParityRegistryError(t *testing.T) {
	err := &registry.RegistryError{StatusCode: 503, Message: "service unavailable"}
	if !containsStr(err.Error(), "503") {
		t.Errorf("expected status code in error: %s", err.Error())
	}
}

// TestParityRegistryErrorNoCode verifies error without status code.
func TestParityRegistryErrorNoCode(t *testing.T) {
	err := &registry.RegistryError{Message: "timeout"}
	if err.Error() != "timeout" {
		t.Errorf("expected plain message, got: %s", err.Error())
	}
}

// TestParityInstallStatusString verifies install status string labels.
func TestParityInstallStatusString(t *testing.T) {
	cases := []struct {
		status registry.InstallStatus
		want   string
	}{
		{registry.StatusNotInstalled, "not-installed"},
		{registry.StatusInstalled, "installed"},
		{registry.StatusConflict, "conflict"},
		{registry.StatusOutdated, "outdated"},
	}
	for _, c := range cases {
		if c.status.String() != c.want {
			t.Errorf("status %d: expected %s, got %s", c.status, c.want, c.status.String())
		}
	}
}

// TestParityParseServerReferenceSimple verifies simple name parsing.
func TestParityParseServerReferenceSimple(t *testing.T) {
	ref := registry.ParseServerReference("my-server")
	if ref.Repo != "my-server" {
		t.Errorf("expected repo=my-server, got %s", ref.Repo)
	}
	if ref.Owner != "" {
		t.Errorf("expected empty owner, got %s", ref.Owner)
	}
	if ref.Version != "" {
		t.Errorf("expected empty version, got %s", ref.Version)
	}
}

// TestParityParseServerReferenceOwnerRepo verifies owner/repo parsing.
func TestParityParseServerReferenceOwnerRepo(t *testing.T) {
	ref := registry.ParseServerReference("acme/my-server")
	if ref.Owner != "acme" {
		t.Errorf("expected owner=acme, got %s", ref.Owner)
	}
	if ref.Repo != "my-server" {
		t.Errorf("expected repo=my-server, got %s", ref.Repo)
	}
}

// TestParityParseServerReferenceWithVersion verifies version extraction.
func TestParityParseServerReferenceWithVersion(t *testing.T) {
	ref := registry.ParseServerReference("acme/my-server@1.2.3")
	if ref.Version != "1.2.3" {
		t.Errorf("expected version=1.2.3, got %s", ref.Version)
	}
	if ref.Owner != "acme" {
		t.Errorf("expected owner=acme, got %s", ref.Owner)
	}
	if ref.Repo != "my-server" {
		t.Errorf("expected repo=my-server, got %s", ref.Repo)
	}
}

// TestParityServerReferenceIsVersioned verifies IsVersioned().
func TestParityServerReferenceIsVersioned(t *testing.T) {
	versioned := registry.ParseServerReference("owner/repo@1.0.0")
	if !versioned.IsVersioned() {
		t.Error("expected IsVersioned=true")
	}
	plain := registry.ParseServerReference("owner/repo")
	if plain.IsVersioned() {
		t.Error("expected IsVersioned=false")
	}
}

// TestParitySemVerString verifies version string formatting.
func TestParitySemVerString(t *testing.T) {
	v := registry.SemVer{Major: 1, Minor: 2, Patch: 3}
	if v.String() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", v.String())
	}
}

// TestParitySemVerStringPre verifies pre-release formatting.
func TestParitySemVerStringPre(t *testing.T) {
	v := registry.SemVer{Major: 2, Minor: 0, Patch: 0, Pre: "beta.1"}
	if v.String() != "2.0.0-beta.1" {
		t.Errorf("expected 2.0.0-beta.1, got %s", v.String())
	}
}

// TestParitySemVerCompareEqual verifies equal version comparison.
func TestParitySemVerCompareEqual(t *testing.T) {
	a := registry.SemVer{Major: 1, Minor: 2, Patch: 3}
	b := registry.SemVer{Major: 1, Minor: 2, Patch: 3}
	if a.Compare(b) != 0 {
		t.Errorf("expected 0, got %d", a.Compare(b))
	}
}

// TestParitySemVerCompareLess verifies less-than comparison.
func TestParitySemVerCompareLess(t *testing.T) {
	a := registry.SemVer{Major: 1, Minor: 0, Patch: 0}
	b := registry.SemVer{Major: 2, Minor: 0, Patch: 0}
	if a.Compare(b) != -1 {
		t.Errorf("expected -1, got %d", a.Compare(b))
	}
}

// TestParitySemVerCompareGreater verifies greater-than comparison.
func TestParitySemVerCompareGreater(t *testing.T) {
	a := registry.SemVer{Major: 2, Minor: 1, Patch: 0}
	b := registry.SemVer{Major: 2, Minor: 0, Patch: 5}
	if a.Compare(b) != 1 {
		t.Errorf("expected 1, got %d", a.Compare(b))
	}
}

// TestParitySemVerPreReleaseLower verifies pre-release sorts lower than release.
func TestParitySemVerPreReleaseLower(t *testing.T) {
	preRelease := registry.SemVer{Major: 1, Minor: 0, Patch: 0, Pre: "alpha"}
	release := registry.SemVer{Major: 1, Minor: 0, Patch: 0}
	if preRelease.Compare(release) != -1 {
		t.Errorf("expected pre-release < release, got %d", preRelease.Compare(release))
	}
}

// TestParityServerEntryFields verifies ServerEntry fields are accessible.
func TestParityServerEntryFields(t *testing.T) {
	e := registry.ServerEntry{
		Name:        "my-server",
		Description: "A test server",
		Source:      "github",
		Repository:  "owner/repo",
		Version:     "1.0.0",
		Tags:        []string{"mcp", "testing"},
	}
	if e.Name != "my-server" {
		t.Errorf("unexpected name: %s", e.Name)
	}
	if len(e.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(e.Tags))
	}
}

// TestParitySearchResultFields verifies SearchResult fields.
func TestParitySearchResultFields(t *testing.T) {
	r := registry.SearchResult{
		Servers: []registry.ServerEntry{{Name: "s1"}, {Name: "s2"}},
		Total:   2,
	}
	if r.Total != 2 {
		t.Errorf("expected total=2, got %d", r.Total)
	}
	if len(r.Servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(r.Servers))
	}
}

// TestParityConflictEntry verifies ConflictEntry fields.
func TestParityConflictEntry(t *testing.T) {
	c := registry.ConflictEntry{
		ServerName:   "my-server",
		Integrations: []string{"claude.json", "vscode.json"},
	}
	if len(c.Integrations) != 2 {
		t.Errorf("expected 2 integrations, got %d", len(c.Integrations))
	}
}

// TestParityDefaultTimeouts verifies timeout constants are reasonable.
func TestParityDefaultTimeouts(t *testing.T) {
	if registry.DefaultConnectTimeout <= 0 {
		t.Error("connect timeout must be positive")
	}
	if registry.DefaultReadTimeout <= 0 {
		t.Error("read timeout must be positive")
	}
	if registry.DefaultReadTimeout <= registry.DefaultConnectTimeout {
		t.Error("read timeout should be greater than connect timeout")
	}
}

func containsStr(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
