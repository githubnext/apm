// Package cachepin provides cache-pin marker functionality for drift-replay correctness.
//
// When apm install populates apm_modules/<owner>/<repo>/ from a specific lockfile
// pin, it drops a small JSON marker (.apm-pin) at the package root recording the
// resolved_commit that produced the cache contents.
//
// apm audit drift-replay verifies the marker matches the lockfile's resolved_commit
// BEFORE diffing.
//
// Schema (v1):
//
//	{"schema_version": 1, "resolved_commit": "<git-sha-or-similar>"}
package cachepin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// MarkerFilename is the name of the cache-pin marker file.
const MarkerFilename = ".apm-pin"

// SchemaVersion is the current schema version.
const SchemaVersion = 1

// CachePinError is raised when the cache pin is missing, malformed, or stale.
type CachePinError struct {
	Msg string
}

func (e *CachePinError) Error() string { return e.Msg }

// IsCachePinError reports whether err is a CachePinError.
func IsCachePinError(err error) bool {
	var t *CachePinError
	return errors.As(err, &t)
}

type markerPayload struct {
	SchemaVersion   int    `json:"schema_version"`
	ResolvedCommit  string `json:"resolved_commit"`
}

// WriteMarker writes the cache-pin marker file to installPath.
//
// Idempotent: overwrites any prior marker. Failures are silent because
// they are non-fatal for apm install itself.
func WriteMarker(installPath, resolvedCommit string) {
	info, err := os.Stat(installPath)
	if err != nil || !info.IsDir() {
		return
	}
	payload := markerPayload{SchemaVersion: SchemaVersion, ResolvedCommit: resolvedCommit}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	markerPath := filepath.Join(installPath, MarkerFilename)
	_ = os.WriteFile(markerPath, data, 0o644)
}

// VerifyMarker verifies the marker at installPath matches expectedCommit.
//
// Returns CachePinError on any of: marker file absent, unreadable, malformed
// JSON, unsupported schema_version, missing resolved_commit field, or commit
// mismatch.
func VerifyMarker(installPath, expectedCommit string) error {
	markerPath := filepath.Join(installPath, MarkerFilename)
	data, err := os.ReadFile(markerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &CachePinError{Msg: fmt.Sprintf("cache-pin marker missing at %s (run apm install to refresh)", installPath)}
		}
		return &CachePinError{Msg: fmt.Sprintf("cannot read cache-pin marker at %s: %v", markerPath, err)}
	}

	var payload markerPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return &CachePinError{Msg: fmt.Sprintf("cache-pin marker at %s is malformed JSON: %v", markerPath, err)}
	}

	if payload.SchemaVersion != SchemaVersion {
		return &CachePinError{Msg: fmt.Sprintf("cache-pin marker at %s has unsupported schema_version %d (expected %d)", markerPath, payload.SchemaVersion, SchemaVersion)}
	}

	if payload.ResolvedCommit == "" {
		return &CachePinError{Msg: fmt.Sprintf("cache-pin marker at %s is missing resolved_commit field", markerPath)}
	}

	if payload.ResolvedCommit != expectedCommit {
		return &CachePinError{Msg: fmt.Sprintf("cache-pin marker mismatch at %s: marker=%s expected=%s (run apm install to refresh)", markerPath, payload.ResolvedCommit, expectedCommit)}
	}

	return nil
}
