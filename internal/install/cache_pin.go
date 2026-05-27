// Package install: cache-pin marker for drift-replay correctness.
//
// When apm install populates apm_modules/<owner>/<repo>/ from a specific
// lockfile pin, it drops a small JSON marker (.apm-pin) at the package root
// recording the resolved_commit that produced the cache contents.
//
// apm audit drift-replay verifies the marker matches the lockfile's
// resolved_commit before diffing. This catches shared CI runner hazards and
// lockfile-bumps without re-running apm install.
//
// Schema (v1):
//
//	{"schema_version": 1, "resolved_commit": "<git-sha-or-similar>"}
package install

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MarkerFilename is the name of the cache-pin file dropped inside each
// package install directory. Mirrors MARKER_FILENAME in cache_pin.py.
const MarkerFilename = ".apm-pin"

// CachePinSchemaVersion is the current marker schema version.
const CachePinSchemaVersion = 1

// CachePinError is returned when the cache pin is missing, malformed, or stale.
// The orchestrator (drift replay) catches this and translates it to a
// CacheMissError with the same message.
type CachePinError struct {
	msg string
}

func (e *CachePinError) Error() string { return e.msg }

// newCachePinError creates a CachePinError with the given message.
func newCachePinError(format string, args ...any) *CachePinError {
	return &CachePinError{msg: fmt.Sprintf(format, args...)}
}

// cachePinMarker is the JSON payload written to MarkerFilename.
type cachePinMarker struct {
	SchemaVersion   int    `json:"schema_version"`
	ResolvedCommit  string `json:"resolved_commit"`
}

// WriteMarker writes the cache-pin marker file to installPath.
//
// Idempotent: overwrites any prior marker. Silent on filesystem errors because
// a missing marker is non-fatal for apm install -- it is detected at
// drift-replay verify time.
func WriteMarker(installPath, resolvedCommit string) {
	info, err := os.Stat(installPath)
	if err != nil || !info.IsDir() {
		return
	}
	payload := cachePinMarker{
		SchemaVersion:  CachePinSchemaVersion,
		ResolvedCommit: resolvedCommit,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	marker := filepath.Join(installPath, MarkerFilename)
	_ = os.WriteFile(marker, data, 0o644)
}

// VerifyMarker verifies that the marker at installPath matches expectedCommit.
//
// Returns a *CachePinError on any of: marker absent, unreadable, malformed JSON,
// unsupported schema_version, missing resolved_commit, or commit mismatch.
func VerifyMarker(installPath, expectedCommit string) error {
	marker := filepath.Join(installPath, MarkerFilename)
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		return newCachePinError(
			"cache pin marker missing at %s -- cache pre-dates supply-chain hardening",
			marker,
		)
	}
	raw, err := os.ReadFile(marker)
	if err != nil {
		return newCachePinError("cache pin marker unreadable at %s: %v", marker, err)
	}
	var payload cachePinMarker
	if err := json.Unmarshal(raw, &payload); err != nil {
		return newCachePinError("cache pin marker at %s is not valid JSON: %v", marker, err)
	}
	if payload.SchemaVersion != CachePinSchemaVersion {
		return newCachePinError(
			"cache pin marker at %s has unsupported schema_version %d; expected %d",
			marker, payload.SchemaVersion, CachePinSchemaVersion,
		)
	}
	if payload.ResolvedCommit == "" {
		return newCachePinError("cache pin marker at %s is missing resolved_commit", marker)
	}
	if payload.ResolvedCommit != expectedCommit {
		return newCachePinError(
			"cache pin mismatch at %s: marker says %q, lockfile expects %q",
			marker, payload.ResolvedCommit, expectedCommit,
		)
	}
	return nil
}
