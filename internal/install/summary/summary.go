// Package summary provides post-install summary rendering helpers.
package summary

import "fmt"

// SummaryResult holds data for a rendered install summary line.
type SummaryResult struct {
ApmCount      int
McpCount      int
Errors        int
StalesCleaned int
ElapsedSecs   float64
}

// FormatSummary returns the install summary line as a string.
func FormatSummary(r SummaryResult) string {
base := fmt.Sprintf("Installed %d APM package(s), %d MCP server(s)", r.ApmCount, r.McpCount)
if r.Errors > 0 {
base += fmt.Sprintf(", %d error(s)", r.Errors)
}
if r.StalesCleaned > 0 {
base += fmt.Sprintf(", cleaned %d stale artifact(s)", r.StalesCleaned)
}
if r.ElapsedSecs > 0 {
base += fmt.Sprintf(" in %.1fs", r.ElapsedSecs)
}
return base + "."
}

// HasCriticalSecurityError returns true when the diagnostic collector signals a critical security finding.
func HasCriticalSecurityError(hasCriticalSecurity bool, force bool) bool {
return !force && hasCriticalSecurity
}
