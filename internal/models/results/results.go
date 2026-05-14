// Package results defines typed result containers for APM operations.
package results

// InstallResult is the result of an APM install operation.
type InstallResult struct {
InstalledCount      int
PromptsIntegrated   int
AgentsIntegrated    int
PackageTypes        map[string]string // dep_key -> type string
}

// PrimitiveCounts holds counts of primitives in a package.
type PrimitiveCounts struct {
Prompts      int
Agents       int
Instructions int
Skills       int
Hooks        int
Commands     int
}
