// Package schema defines frozen data models for the apm-policy.yml schema.
package schema

// PolicyCache holds cache configuration for remote policy resolution.
type PolicyCache struct {
TTL int // seconds, default 3600
}

// DependencyPolicy defines rules governing which APM dependencies are permitted.
type DependencyPolicy struct {
Allow             []string
Deny              []string
Require           []string
RequireResolution string // project-wins | policy-wins | block
MaxDepth          int    // default 50
}

// McpTransportPolicy defines allowed MCP transport protocols.
type McpTransportPolicy struct {
Allow []string // stdio, sse, http, streamable-http
}

// McpPolicy defines rules governing MCP server references.
type McpPolicy struct {
Allow           []string
Deny            []string
Transport       McpTransportPolicy
SelfDefined     string // deny | warn | allow
TrustTransitive bool
}

// CompilationTargetPolicy defines allowed compilation targets.
type CompilationTargetPolicy struct {
Allow   []string // vscode, claude, all
Enforce string
}

// CompilationStrategyPolicy defines compilation strategy constraints.
type CompilationStrategyPolicy struct {
Enforce string // distributed | single-file
}

// CompilationPolicy bundles target and strategy policies.
type CompilationPolicy struct {
Targets  CompilationTargetPolicy
Strategy CompilationStrategyPolicy
}

// ApmPolicy is the root policy object parsed from apm-policy.yml.
type ApmPolicy struct {
Version     string
Remote      string
Cache       PolicyCache
Deps        DependencyPolicy
MCP         McpPolicy
Compilation CompilationPolicy
}

// DefaultDependencyPolicy returns a DependencyPolicy with sensible defaults.
func DefaultDependencyPolicy() DependencyPolicy {
return DependencyPolicy{
RequireResolution: "project-wins",
MaxDepth:          50,
}
}
