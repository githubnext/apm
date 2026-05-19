package discovery

import (
	"testing"
)

func TestPolicyFetchResultFields(t *testing.T) {
	r := PolicyFetchResult{
		Source:          "org:myorg/.github",
		Cached:          true,
		CacheAgeSeconds: 120,
		CacheStale:      false,
		Outcome:         "hit",
		RawBytesHash:    "sha256:abc",
		ExpectedHash:    "sha256:abc",
	}
	if r.Found() {
		t.Error("Found should be false without Policy set")
	}
	if r.Source != "org:myorg/.github" {
		t.Errorf("Source mismatch")
	}
	if !r.Cached {
		t.Error("Cached should be true")
	}
	if r.CacheAgeSeconds != 120 {
		t.Errorf("CacheAgeSeconds: got %d, want 120", r.CacheAgeSeconds)
	}
}

func TestPolicyFetchResultZeroValue(t *testing.T) {
	var r PolicyFetchResult
	if r.Found() {
		t.Error("zero value Found should be false")
	}
	if r.Cached {
		t.Error("zero value Cached should be false")
	}
	if r.CacheStale {
		t.Error("zero value CacheStale should be false")
	}
	if r.CacheAgeSeconds != 0 {
		t.Errorf("zero value CacheAgeSeconds should be 0")
	}
}

func TestSplitHashPinEmptyString(t *testing.T) {
	_, _, err := splitHashPin("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestSplitHashPinNoColon(t *testing.T) {
	// 64 hex chars with no colon treated as bare sha256
	validHex := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	algo, hex, err := splitHashPin(validHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if algo != "sha256" {
		t.Errorf("expected sha256 algo, got %q", algo)
	}
	if hex != validHex {
		t.Errorf("hex mismatch")
	}
}

func TestCacheConstants(t *testing.T) {
	if defaultCacheTTL <= 0 {
		t.Error("defaultCacheTTL should be positive")
	}
	if maxStaleTTL <= defaultCacheTTL {
		t.Error("maxStaleTTL should exceed defaultCacheTTL")
	}
	if cacheSchemaVersion == "" {
		t.Error("cacheSchemaVersion should not be empty")
	}
}
