package discovery

import (
	"testing"
)

// ---------------------------------------------------------------------------
// parseRemoteURL additional cases
// ---------------------------------------------------------------------------

func TestParseRemoteURLSCPLike(t *testing.T) {
	org, host, err := parseRemoteURL("git@github.com:myorg/myrepo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "github.com" {
		t.Errorf("expected host=github.com, got %q", host)
	}
	if org != "myorg" {
		t.Errorf("expected org=myorg, got %q", org)
	}
}

func TestParseRemoteURLHTTPSWithToken(t *testing.T) {
	org, host, err := parseRemoteURL("https://token:x-oauth@github.com/corp/repo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "github.com" {
		t.Errorf("expected github.com, got %q", host)
	}
	if org != "corp" {
		t.Errorf("expected org=corp, got %q", org)
	}
}

func TestParseRemoteURLSubOrg(t *testing.T) {
	org, host, err := parseRemoteURL("https://github.com/parent/suborg/repo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "github.com" {
		t.Errorf("expected github.com, got %q", host)
	}
	_ = org
}

// ---------------------------------------------------------------------------
// PolicyFetchResult additional fields and behaviors
// ---------------------------------------------------------------------------

func TestPolicyFetchResultFoundFalse(t *testing.T) {
	r := &PolicyFetchResult{}
	if r.Found() {
		t.Error("zero value should not be Found")
	}
}

func TestPolicyFetchResultErrField(t *testing.T) {
	r := &PolicyFetchResult{Err: "not found"}
	if r.Err != "not found" {
		t.Errorf("expected Err='not found', got %q", r.Err)
	}
}

func TestPolicyFetchResultCachedField(t *testing.T) {
	r := &PolicyFetchResult{Cached: true, CacheAgeSeconds: 300}
	if !r.Cached {
		t.Error("expected Cached=true")
	}
	if r.CacheAgeSeconds != 300 {
		t.Errorf("expected CacheAgeSeconds=300, got %d", r.CacheAgeSeconds)
	}
}

func TestPolicyFetchResultOutcomeField(t *testing.T) {
	r := &PolicyFetchResult{Outcome: "absent"}
	if r.Outcome != "absent" {
		t.Errorf("expected Outcome=absent, got %q", r.Outcome)
	}
}

func TestPolicyFetchResultRawBytesHash(t *testing.T) {
	r := &PolicyFetchResult{RawBytesHash: "sha256:abc123"}
	if r.RawBytesHash != "sha256:abc123" {
		t.Errorf("expected sha256:abc123, got %q", r.RawBytesHash)
	}
}

// ---------------------------------------------------------------------------
// splitHashPin edge cases
// ---------------------------------------------------------------------------

func TestSplitHashPinMultiColon(t *testing.T) {
	// sha256 hex must be exactly 64 chars
	sha256hex := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	_, _, err := splitHashPin("sha256:" + sha256hex)
	if err != nil {
		t.Fatalf("unexpected error for valid sha256: %v", err)
	}
}

func TestSplitHashPinLowercase(t *testing.T) {
	_, _, err := splitHashPin("sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err != nil {
		t.Fatalf("unexpected error for valid sha256: %v", err)
	}
}

// ---------------------------------------------------------------------------
// verifyHashPin with empty pin (no-op)
// ---------------------------------------------------------------------------

func TestVerifyHashPinEmptyPinIsNil(t *testing.T) {
	r := verifyHashPin([]byte("content"), "", "test-source")
	if r != nil {
		t.Error("expected nil result for empty expectedHash")
	}
}

// ---------------------------------------------------------------------------
// computeHashNormalized
// ---------------------------------------------------------------------------

func TestComputeHashNormalizedSHA256(t *testing.T) {
	got := computeHashNormalized([]byte("hello"), "sha256:")
	if got == "" {
		t.Error("expected non-empty hash for sha256")
	}
}

func TestComputeHashNormalizedSHA512(t *testing.T) {
	got := computeHashNormalized([]byte("hello"), "sha512:")
	if got == "" {
		t.Error("expected non-empty hash for sha512")
	}
}
