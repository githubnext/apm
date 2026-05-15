// Package conflictdetector handles MCP server configuration conflict detection.
// Migrated from src/apm_cli/core/conflict_detector.py
package conflictdetector

import "strings"

// ServerConfig represents an MCP server configuration entry.
type ServerConfig map[string]interface{}

// MCPConflictDetector detects and resolves MCP server configuration conflicts.
type MCPConflictDetector struct {
	// GetExistingServers returns the current set of server configs (name -> config).
	GetExistingServers func() map[string]ServerConfig
	// ResolveCanonicalName maps a server reference to its canonical config name.
	ResolveCanonicalName func(ref string) (string, error)
	// FindServerByReference looks up a server in the registry by reference.
	FindServerByReference func(ref string) (map[string]interface{}, error)
}

// New creates a MCPConflictDetector with the supplied callbacks.
func New(
	getServers func() map[string]ServerConfig,
	resolveCanon func(ref string) (string, error),
	findServer func(ref string) (map[string]interface{}, error),
) *MCPConflictDetector {
	return &MCPConflictDetector{
		GetExistingServers:    getServers,
		ResolveCanonicalName:  resolveCanon,
		FindServerByReference: findServer,
	}
}

// ServerExistsResult is the outcome of a ServerExists check.
type ServerExistsResult struct {
	Exists        bool
	ConflictName  string
	ConflictUUID  string
}

// CheckServerExists reports whether serverRef already exists in the configuration.
func (d *MCPConflictDetector) CheckServerExists(serverRef string) ServerExistsResult {
	existing := d.GetExistingServers()

	// Try UUID-based lookup via registry first
	if d.FindServerByReference != nil {
		if info, err := d.FindServerByReference(serverRef); err == nil && info != nil {
			if uuid, ok := info["id"].(string); ok && uuid != "" {
				for name, cfg := range existing {
					if val, ok := cfg["id"].(string); ok && val == uuid {
						return ServerExistsResult{Exists: true, ConflictName: name, ConflictUUID: uuid}
					}
				}
			}
		}
	}

	// Fall back to canonical name comparison
	canonical := d.canonicalName(serverRef)
	if _, ok := existing[canonical]; ok {
		return ServerExistsResult{Exists: true, ConflictName: canonical}
	}
	for name := range existing {
		if name == canonical {
			continue
		}
		existingCanon := d.canonicalName(name)
		if existingCanon == canonical {
			return ServerExistsResult{Exists: true, ConflictName: name}
		}
	}
	return ServerExistsResult{Exists: false}
}

func (d *MCPConflictDetector) canonicalName(ref string) string {
	if d.ResolveCanonicalName != nil {
		if name, err := d.ResolveCanonicalName(ref); err == nil {
			return name
		}
	}
	// Fallback: lowercase last path component
	parts := strings.Split(strings.TrimSuffix(ref, "/"), "/")
	if len(parts) == 0 {
		return strings.ToLower(ref)
	}
	return strings.ToLower(parts[len(parts)-1])
}

// GetExistingServerConfigs returns the current server configurations.
func (d *MCPConflictDetector) GetExistingServerConfigs() map[string]ServerConfig {
	if d.GetExistingServers == nil {
		return map[string]ServerConfig{}
	}
	return d.GetExistingServers()
}

// GetCanonicalServerName resolves a reference to its canonical config name.
func (d *MCPConflictDetector) GetCanonicalServerName(ref string) string {
	return d.canonicalName(ref)
}

// FindConflicts returns all existing server names that conflict with serverRef.
func (d *MCPConflictDetector) FindConflicts(serverRef string) []string {
	result := d.CheckServerExists(serverRef)
	if !result.Exists {
		return nil
	}
	return []string{result.ConflictName}
}
