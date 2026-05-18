package discovery

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// PolicyFetchResult.Found
// ---------------------------------------------------------------------------

func TestPolicyFetchResultFound(t *testing.T) {
	r := &PolicyFetchResult{}
	if r.Found() {
		t.Error("expected Found=false when Policy is nil")
	}
}

// ---------------------------------------------------------------------------
// splitHashPin
// ---------------------------------------------------------------------------

func TestSplitHashPinWithAlgo(t *testing.T) {
	validHex := strings.Repeat("a", 64)
	algo, hex, err := splitHashPin("sha256:" + validHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if algo != "sha256" || hex != validHex {
		t.Errorf("got algo=%q hex=%q", algo, hex)
	}
}

func TestSplitHashPinBareHex(t *testing.T) {
	validHex := strings.Repeat("b", 64)
	algo, hex, err := splitHashPin(validHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if algo != "sha256" || hex != validHex {
		t.Errorf("got algo=%q hex=%q", algo, hex)
	}
}

func TestSplitHashPinInvalidAlgo(t *testing.T) {
	_, _, err := splitHashPin("md5:" + strings.Repeat("c", 64))
	if err == nil {
		t.Error("expected error for unsupported algo")
	}
}

func TestSplitHashPinTooShort(t *testing.T) {
	_, _, err := splitHashPin("sha256:abc")
	if err == nil {
		t.Error("expected error for short hex")
	}
}

// ---------------------------------------------------------------------------
// verifyHashPin
// ---------------------------------------------------------------------------

func TestVerifyHashPinEmpty(t *testing.T) {
	result := verifyHashPin([]byte("content"), "", "src")
	if result != nil {
		t.Errorf("expected nil for empty pin, got %+v", result)
	}
}

func TestVerifyHashPinMatch(t *testing.T) {
	content := []byte("policy content")
	h := sha256.Sum256(content)
	pin := fmt.Sprintf("sha256:%x", h)
	result := verifyHashPin(content, pin, "src")
	if result != nil {
		t.Errorf("expected nil (match), got %+v", result)
	}
}

func TestVerifyHashPinMismatch(t *testing.T) {
	content := []byte("policy content")
	pin := "sha256:" + strings.Repeat("0", 64)
	result := verifyHashPin(content, pin, "src")
	if result == nil {
		t.Error("expected mismatch result")
	}
	if result.Outcome != "hash_mismatch" {
		t.Errorf("unexpected outcome: %q", result.Outcome)
	}
}

func TestVerifyHashPinInvalidPin(t *testing.T) {
	result := verifyHashPin([]byte("x"), "md5:abc", "src")
	if result == nil {
		t.Error("expected error result for invalid pin")
	}
}

// ---------------------------------------------------------------------------
// computeHashNormalized
// ---------------------------------------------------------------------------

func TestComputeHashNormalized(t *testing.T) {
	content := []byte("hello world")
	h := computeHashNormalized(content, "")
	if !strings.HasPrefix(h, "sha256:") {
		t.Errorf("expected sha256: prefix, got %q", h)
	}
	expected := fmt.Sprintf("sha256:%x", sha256.Sum256(content))
	if h != expected {
		t.Errorf("got %q want %q", h, expected)
	}
}

// ---------------------------------------------------------------------------
// parseRemoteURL
// ---------------------------------------------------------------------------

func TestParseRemoteURLHTTPS(t *testing.T) {
	cases := []struct {
		url, wantOrg, wantHost string
	}{
		{"https://github.com/myorg/myrepo.git", "myorg", "github.com"},
		{"https://github.com/myorg/myrepo", "myorg", "github.com"},
		{"https://myhost.ghe.com/contoso/project", "contoso", "myhost.ghe.com"},
	}
	for _, c := range cases {
		org, host, err := parseRemoteURL(c.url)
		if err != nil {
			t.Errorf("parseRemoteURL(%q): unexpected error: %v", c.url, err)
			continue
		}
		if org != c.wantOrg || host != c.wantHost {
			t.Errorf("parseRemoteURL(%q) = (%q, %q), want (%q, %q)", c.url, org, host, c.wantOrg, c.wantHost)
		}
	}
}

func TestParseRemoteURLSSH(t *testing.T) {
	org, host, err := parseRemoteURL("git@github.com:myorg/myrepo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org != "myorg" || host != "github.com" {
		t.Errorf("got org=%q host=%q", org, host)
	}
}

func TestParseRemoteURLEmpty(t *testing.T) {
	_, _, err := parseRemoteURL("")
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestParseRemoteURLInvalid(t *testing.T) {
	_, _, err := parseRemoteURL("not-a-valid-url")
	if err == nil {
		t.Error("expected error for unparseable URL")
	}
}

// ---------------------------------------------------------------------------
// loadFromFile
// ---------------------------------------------------------------------------

func TestLoadFromFileNotFound(t *testing.T) {
	r := loadFromFile("/nonexistent/path/file.yml", "")
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if r.Err == "" {
		t.Error("expected error message for missing file")
	}
}

func TestLoadFromFileValidPolicy(t *testing.T) {
	dir := t.TempDir()
	policyContent := "version: 1\nrules: []\n"
	p := filepath.Join(dir, "apm-policy.yml")
	if err := os.WriteFile(p, []byte(policyContent), 0o644); err != nil {
		t.Fatal(err)
	}
	r := loadFromFile(p, "")
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	// Should have no error even if policy is minimal.
	if r.Outcome == "malformed" {
		t.Errorf("unexpected malformed outcome: %s", r.Err)
	}
}

func TestLoadFromFileHashMismatch(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "apm-policy.yml")
	_ = os.WriteFile(p, []byte("version: 1\n"), 0o644)
	badPin := "sha256:" + strings.Repeat("0", 64)
	r := loadFromFile(p, badPin)
	if r == nil {
		t.Fatal("expected non-nil")
	}
	if r.Outcome != "hash_mismatch" {
		t.Errorf("expected hash_mismatch, got %q", r.Outcome)
	}
}

// ---------------------------------------------------------------------------
// cacheKey
// ---------------------------------------------------------------------------

func TestCacheKeyDeterministic(t *testing.T) {
	k1 := cacheKey("org/repo")
	k2 := cacheKey("org/repo")
	if k1 != k2 {
		t.Errorf("cacheKey should be deterministic: %q vs %q", k1, k2)
	}
	k3 := cacheKey("other/repo")
	if k1 == k3 {
		t.Error("different inputs should produce different keys")
	}
}

// ---------------------------------------------------------------------------
// DiscoverPolicy (file override)
// ---------------------------------------------------------------------------

func TestDiscoverPolicyFromFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.yml")
	_ = os.WriteFile(p, []byte("version: 1\nrules: []\n"), 0o644)
	r := DiscoverPolicy(dir, p, true, "")
	if r == nil {
		t.Fatal("expected result")
	}
	if !strings.HasPrefix(r.Source, "file:") {
		t.Errorf("expected file: source, got %q", r.Source)
	}
}
