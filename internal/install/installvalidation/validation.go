// Package installvalidation provides package existence and validation helpers.
// Migrated from src/apm_cli/install/validation.py
package installvalidation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TLSErrorPrefix is the prefix placed on errors raised for TLS verification failures.
const TLSErrorPrefix = "TLS verification failed"

// AuthenticationError indicates an authentication failure during package validation.
type AuthenticationError struct {
	Host    string
	Message string
}

func (e *AuthenticationError) Error() string {
	if e.Host != "" {
		return fmt.Sprintf("authentication failed for %s: %s", e.Host, e.Message)
	}
	return "authentication failed: " + e.Message
}

// TLSError wraps a TLS verification failure.
type TLSError struct {
	Host  string
	Cause error
}

func (e *TLSError) Error() string {
	msg := TLSErrorPrefix
	if e.Host != "" {
		msg += " for " + e.Host
	}
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}
	return msg
}

func (e *TLSError) Unwrap() error { return e.Cause }

// IsTLSFailure reports whether err (or any cause in its chain) is a TLS failure.
func IsTLSFailure(err error) bool {
	if err == nil {
		return false
	}
	var te *TLSError
	if errors.As(err, &te) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, TLSErrorPrefix) || strings.Contains(msg, "CERTIFICATE_VERIFY_FAILED")
}

// LocalPathMarkers are file/dir names that indicate an installable APM package.
var LocalPathMarkers = []string{"apm.yml", "apm.yaml", ".apm"}

// LocalPathFailureReason returns a human-readable message when a local-path dep fails.
func LocalPathFailureReason(localPath string) string {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Sprintf("local path %q does not exist", localPath)
	}
	for _, marker := range LocalPathMarkers {
		if _, err := os.Stat(filepath.Join(localPath, marker)); err == nil {
			return "" // found a marker; path is valid
		}
	}
	return fmt.Sprintf("local path %q exists but contains no apm.yml/.apm marker", localPath)
}

// LocalPathNoMarkersHint scans a directory for nested installable packages
// and returns a hint string for the user.
func LocalPathNoMarkersHint(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var candidates []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		for _, marker := range LocalPathMarkers {
			if _, err := os.Stat(filepath.Join(dir, e.Name(), marker)); err == nil {
				candidates = append(candidates, e.Name())
				break
			}
		}
		if len(candidates) >= 5 {
			break
		}
	}
	if len(candidates) == 0 {
		return ""
	}
	return fmt.Sprintf("Found installable packages in sub-directories: %s", strings.Join(candidates, ", "))
}

// PackageProber probes a package reference for reachability.
type PackageProber struct {
	AuthToken  string
	Host       string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// NewPackageProber creates a PackageProber with default settings.
func NewPackageProber(host, authToken string) *PackageProber {
	return &PackageProber{
		Host:       host,
		AuthToken:  authToken,
		Timeout:    15 * time.Second,
		HTTPClient: http.DefaultClient,
	}
}

// ProbeResult is the outcome of a package probe.
type ProbeResult struct {
	Reachable bool
	// Reason is set when Reachable is false.
	Reason string
	// IsAuthError is true when the failure is an authentication problem.
	IsAuthError bool
	// IsTLSError is true when the failure is a TLS verification problem.
	IsTLSError bool
}

// ProbeGitHubAPI checks whether owner/repo is accessible via the GitHub API.
func (p *PackageProber) ProbeGitHubAPI(owner, repo, ref string) ProbeResult {
	apiBase := "https://api.github.com"
	if p.Host != "github.com" {
		apiBase = "https://" + p.Host + "/api/v3"
	}
	url := fmt.Sprintf("%s/repos/%s/%s", apiBase, owner, repo)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ProbeResult{Reachable: false, Reason: err.Error()}
	}
	if p.AuthToken != "" {
		req.Header.Set("Authorization", "token "+p.AuthToken)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		if IsTLSFailure(err) {
			return ProbeResult{Reachable: false, Reason: err.Error(), IsTLSError: true}
		}
		return ProbeResult{Reachable: false, Reason: err.Error()}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return ProbeResult{Reachable: true}
	case http.StatusUnauthorized, http.StatusForbidden:
		return ProbeResult{Reachable: false, Reason: "authentication failed", IsAuthError: true}
	case http.StatusNotFound:
		return ProbeResult{Reachable: false, Reason: fmt.Sprintf("repository %s/%s not found", owner, repo)}
	default:
		return ProbeResult{Reachable: false, Reason: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}
}

// IsADOAuthFailureSignal reports whether an HTTP status or message looks like an ADO auth failure.
func IsADOAuthFailureSignal(statusCode int, body string) bool {
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return true
	}
	lower := strings.ToLower(body)
	return strings.Contains(lower, "tfs auth") ||
		strings.Contains(lower, "azure devops") ||
		strings.Contains(lower, "unauthorized")
}

// ValidatePackageExists is the main entry point: probe whether a package ref is reachable.
func ValidatePackageExists(
	pkg string,
	host string,
	authToken string,
	verbose bool,
) ProbeResult {
	prober := NewPackageProber(host, authToken)

	// Local path fast-path
	if strings.HasPrefix(pkg, "./") || strings.HasPrefix(pkg, "../") || filepath.IsAbs(pkg) {
		reason := LocalPathFailureReason(pkg)
		if reason == "" {
			return ProbeResult{Reachable: true}
		}
		return ProbeResult{Reachable: false, Reason: reason}
	}

	// Parse owner/repo from spec
	ref := ""
	spec := pkg
	if idx := strings.LastIndex(spec, "#"); idx >= 0 {
		ref = spec[idx+1:]
		spec = spec[:idx]
	}
	parts := strings.SplitN(spec, "/", 3)
	if len(parts) < 2 {
		return ProbeResult{Reachable: false, Reason: "invalid package spec: " + pkg}
	}
	owner := parts[0]
	repo := parts[1]

	return prober.ProbeGitHubAPI(owner, repo, ref)
}
