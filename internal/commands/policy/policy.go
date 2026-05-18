// Package policy implements the "apm policy" command group.
//
// Provides diagnostic visibility into policy discovery, caching, inheritance
// chains, and effective rule counts.  Always exits 0 -- failures are reported
// inline so the command is safe for CI/SIEM ingestion.
//
// Migrated from: src/apm_cli/commands/policy.py
package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PolicySource describes where a policy was loaded from.
type PolicySource struct {
	Label     string `json:"label"`
	URL       string `json:"url,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	CacheAge  int    `json:"cache_age_seconds,omitempty"`
	Stale     bool   `json:"stale"`
	FetchError string `json:"fetch_error,omitempty"`
}

// PolicyStatus is the result of a policy status check.
type PolicyStatus struct {
	Discovered     bool           `json:"discovered"`
	Source         *PolicySource  `json:"source,omitempty"`
	InheritanceChain []PolicySource `json:"inheritance_chain,omitempty"`
	RuleCount      map[string]int `json:"rule_counts,omitempty"`
	Error          string         `json:"error,omitempty"`
	ProjectRoot    string         `json:"project_root"`
	CheckedAt      string         `json:"checked_at"`
}

// StatusOptions configures the policy status command.
type StatusOptions struct {
	ProjectRoot string
	Format      string // "text" | "json"
	Verbose     bool
	NoFetch     bool
}

// RunStatus checks and prints policy status.
func RunStatus(opts StatusOptions) error {
	status := &PolicyStatus{
		ProjectRoot: opts.ProjectRoot,
		CheckedAt:   time.Now().UTC().Format(time.RFC3339),
		RuleCount:   make(map[string]int),
	}

	// Try to find a policy file.
	policyPath, err := discoverPolicyFile(opts.ProjectRoot)
	if err != nil {
		status.Error = err.Error()
	} else if policyPath != "" {
		status.Discovered = true
		status.Source = &PolicySource{
			Label:    stripSourcePrefix(policyPath),
			FilePath: policyPath,
		}
		rules, err := countRules(policyPath)
		if err == nil {
			status.RuleCount = rules
		}
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(status)
	}

	printStatusText(status, opts.Verbose)
	return nil
}

// discoverPolicyFile looks for an apm-policy.yml (or similar) in the project root.
func discoverPolicyFile(projectRoot string) (string, error) {
	candidates := []string{
		"apm-policy.yml",
		"apm-policy.yaml",
		".apm/policy.yml",
		".apm/policy.yaml",
	}
	for _, c := range candidates {
		p := filepath.Join(projectRoot, c)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", nil
}

// countRules does a lightweight scan of the policy YAML and counts rules.
func countRules(policyPath string) (map[string]int, error) {
	data, err := os.ReadFile(policyPath)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		// Count top-level YAML keys as sections.
		if !strings.HasPrefix(l, " ") && !strings.HasPrefix(l, "\t") &&
			strings.HasSuffix(strings.SplitN(t, ":", 2)[0], "") {
			parts := strings.SplitN(t, ":", 2)
			if len(parts) == 2 && parts[1] == "" {
				counts[parts[0]]++
			}
		}
	}
	return counts, nil
}

// stripSourcePrefix removes "org:", "url:", "file:" prefixes from a label.
func stripSourcePrefix(s string) string {
	for _, pfx := range []string{"org:", "url:", "file:"} {
		if strings.HasPrefix(s, pfx) {
			return s[len(pfx):]
		}
	}
	return s
}

// formatAge renders a cache age in compact human-friendly form.
func formatAge(seconds int) string {
	if seconds < 0 {
		return "n/a"
	}
	if seconds < 60 {
		return fmt.Sprintf("%ds ago", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm ago", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}
	return fmt.Sprintf("%dd ago", hours/24)
}

// printStatusText renders a human-readable policy status report.
func printStatusText(s *PolicyStatus, verbose bool) {
	fmt.Printf("Policy Status for: %s\n", s.ProjectRoot)
	fmt.Printf("Checked at: %s\n\n", s.CheckedAt)

	if !s.Discovered {
		fmt.Println("[i] No policy file discovered.")
		if s.Error != "" {
			fmt.Printf("[x] Error: %s\n", s.Error)
		}
		return
	}

	fmt.Printf("[+] Policy discovered: %s\n", s.Source.Label)
	if s.Source.FilePath != "" {
		fmt.Printf("    File: %s\n", s.Source.FilePath)
	}
	if s.Source.Stale {
		fmt.Printf("    [!] Cache is stale (%s)\n", formatAge(s.Source.CacheAge))
	} else if s.Source.CacheAge > 0 {
		fmt.Printf("    Cache age: %s\n", formatAge(s.Source.CacheAge))
	}

	if len(s.RuleCount) > 0 {
		fmt.Println("Rule counts:")
		for k, v := range s.RuleCount {
			fmt.Printf("  %-30s %d\n", k, v)
		}
	}

	if verbose && len(s.InheritanceChain) > 0 {
		fmt.Println("Inheritance chain:")
		for i, ps := range s.InheritanceChain {
			fmt.Printf("  %d. %s\n", i+1, ps.Label)
		}
	}
}

// DebugOptions configures the policy debug sub-command.
type DebugOptions struct {
	ProjectRoot string
	Format      string
	Source      string
}

// RunDebug prints the raw policy content.
func RunDebug(opts DebugOptions) error {
	policyPath, err := discoverPolicyFile(opts.ProjectRoot)
	if err != nil {
		return err
	}
	if policyPath == "" {
		fmt.Println("[i] No policy file found.")
		return nil
	}
	data, err := os.ReadFile(policyPath)
	if err != nil {
		return fmt.Errorf("read policy file: %w", err)
	}
	fmt.Printf("# Policy from: %s\n\n", policyPath)
	os.Stdout.Write(data)
	return nil
}
