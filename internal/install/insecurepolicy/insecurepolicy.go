// Package insecurepolicy validates HTTP dependency policy for apm install.
// Mirrors src/apm_cli/install/insecure_policy.py.
package insecurepolicy

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// InsecureDependencyPolicyError is returned when HTTP dep policy blocks the install.
type InsecureDependencyPolicyError struct {
	Message string
}

func (e *InsecureDependencyPolicyError) Error() string { return e.Message }

// InsecureDependencyInfo holds resolved details for one insecure dependency.
type InsecureDependencyInfo struct {
	URL          string
	IsTransitive bool
	IntroducedBy string
}

// fqdnRe is a minimal FQDN validator matching the Python is_valid_fqdn logic.
var fqdnRe = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)

// IsValidFQDN returns true for valid fully-qualified domain names.
func IsValidFQDN(host string) bool {
	return fqdnRe.MatchString(strings.ToLower(strings.TrimSpace(host)))
}

// NormalizeAllowInsecureHost validates and normalises a hostname passed via
// --allow-insecure-host.
func NormalizeAllowInsecureHost(hostname string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(hostname))
	if !IsValidFQDN(normalized) {
		return "", fmt.Errorf("invalid hostname %q. Use a bare hostname like 'mirror.example.com'", hostname)
	}
	return normalized, nil
}

// GetInsecureDependencyHost extracts the hostname from an InsecureDependencyInfo URL.
func GetInsecureDependencyHost(info InsecureDependencyInfo) string {
	u, err := url.Parse(info.URL)
	if err != nil || u.Hostname() == "" {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

// FormatInsecureDependencyRequirements renders the canonical remediation message.
func FormatInsecureDependencyRequirements(
	u string,
	missingDepAllow bool,
	missingCLIFlag bool,
) string {
	lines := []string{
		fmt.Sprintf("%s -- HTTP dependency (unencrypted)", u),
		"To install:",
	}
	step := 1
	if missingDepAllow {
		lines = append(lines, fmt.Sprintf("  %d. Set allow_insecure: true on the dep in apm.yml", step))
		step++
	}
	if missingCLIFlag {
		lines = append(lines, fmt.Sprintf("  %d. Pass --allow-insecure to apm install", step))
	}
	return strings.Join(lines, "\n")
}

// FormatInsecureDependencyWarning renders install-time warning text.
func FormatInsecureDependencyWarning(info InsecureDependencyInfo) string {
	msg := fmt.Sprintf("Insecure HTTP fetch (unencrypted): %s", info.URL)
	if info.IsTransitive && info.IntroducedBy != "" {
		msg = fmt.Sprintf("%s (transitive, introduced by %s)", msg, info.IntroducedBy)
	}
	return msg
}

// GetAllowedTransitiveInsecureHosts builds the hostname allowlist for transitive deps.
func GetAllowedTransitiveInsecureHosts(
	infos []InsecureDependencyInfo,
	allowInsecure bool,
	allowInsecureHosts []string,
) map[string]bool {
	allowed := map[string]bool{}
	for _, h := range allowInsecureHosts {
		allowed[h] = true
	}
	if !allowInsecure {
		return allowed
	}
	for _, info := range infos {
		if info.IsTransitive {
			continue
		}
		if h := GetInsecureDependencyHost(info); h != "" {
			allowed[h] = true
		}
	}
	return allowed
}

// GuardTransitiveInsecureDependencies blocks transitive insecure deps from
// unapproved hosts. Returns an error when policy is violated.
func GuardTransitiveInsecureDependencies(
	infos []InsecureDependencyInfo,
	allowInsecure bool,
	allowInsecureHosts []string,
) error {
	var transitive []InsecureDependencyInfo
	for _, info := range infos {
		if info.IsTransitive {
			transitive = append(transitive, info)
		}
	}
	if len(transitive) == 0 {
		return nil
	}

	allowed := GetAllowedTransitiveInsecureHosts(infos, allowInsecure, allowInsecureHosts)
	blockedSet := map[string]bool{}
	for _, info := range transitive {
		h := GetInsecureDependencyHost(info)
		if h != "" && !allowed[h] {
			blockedSet[h] = true
		}
	}
	if len(blockedSet) == 0 {
		return nil
	}

	var blocked []string
	for h := range blockedSet {
		blocked = append(blocked, h)
	}
	sort.Strings(blocked)

	var flagParts []string
	for _, h := range blocked {
		flagParts = append(flagParts, "--allow-insecure-host "+h)
	}
	msg := fmt.Sprintf(
		"Re-run with %s to allow transitive HTTP dependencies from unapproved host(s): %s.",
		strings.Join(flagParts, " "),
		strings.Join(blocked, ", "),
	)
	return &InsecureDependencyPolicyError{Message: msg}
}
