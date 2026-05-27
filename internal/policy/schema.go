// Package policy provides schema types for apm-policy.yml.
// It mirrors the Python apm_cli.policy.schema module.
package policy

// PolicyCache holds cache configuration for remote policy resolution.
type PolicyCache struct {
	TTL int // seconds, default 3600
}

// DefaultPolicyCache returns the default cache config.
func DefaultPolicyCache() PolicyCache {
	return PolicyCache{TTL: 3600}
}

// RequireResolution is the resolution mode for required packages.
type RequireResolution string

const (
	RequireResolutionProjectWins RequireResolution = "project-wins"
	RequireResolutionPolicyWins  RequireResolution = "policy-wins"
	RequireResolutionBlock       RequireResolution = "block"
)

// DependencyPolicy rules governing which APM dependencies are permitted.
type DependencyPolicy struct {
	Allow             []string          // nil = no opinion; empty = explicitly nothing
	Deny              []string          // nil = no opinion; empty = explicit empty
	Require           []string          // nil = no opinion; empty = explicit empty
	RequireResolution RequireResolution // default: project-wins
	MaxDepth          int               // default: 50
}

// DefaultDependencyPolicy returns the default dependency policy.
func DefaultDependencyPolicy() DependencyPolicy {
	return DependencyPolicy{
		RequireResolution: RequireResolutionProjectWins,
		MaxDepth:          50,
	}
}

// EffectiveDeny returns the resolved deny list (nil -> empty slice).
func (d DependencyPolicy) EffectiveDeny() []string {
	if d.Deny == nil {
		return []string{}
	}
	return d.Deny
}

// EffectiveRequire returns the resolved require list (nil -> empty slice).
func (d DependencyPolicy) EffectiveRequire() []string {
	if d.Require == nil {
		return []string{}
	}
	return d.Require
}

// McpTransportPolicy describes allowed MCP transport protocols.
type McpTransportPolicy struct {
	Allow []string // e.g. stdio, sse, http, streamable-http
}

// McpPolicy rules governing MCP server references.
type McpPolicy struct {
	Allow     []string           // nil = no opinion
	Deny      []string           // nil = no opinion
	Transport McpTransportPolicy
	AllowSelf *bool // nil = no opinion
}

// CompilationPolicy governs compilation targets and strategies.
type CompilationPolicy struct {
	AllowedTargets    []string // nil = no opinion
	AllowedStrategies []string // nil = no opinion
}

// ScriptsPolicy governs allowed/denied scripts.
type ScriptsPolicy struct {
	Allow []string // nil = no opinion
	Deny  []string // nil = no opinion
}

// OutcomeAction is the action taken when a policy check fails.
type OutcomeAction string

const (
	OutcomeBlock OutcomeAction = "block"
	OutcomeWarn  OutcomeAction = "warn"
	OutcomeAllow OutcomeAction = "allow"
)

// OutcomeRouting maps check names to their outcome action.
type OutcomeRouting struct {
	Default OutcomeAction
	Checks  map[string]OutcomeAction
}

// DefaultOutcomeRouting returns routing that blocks on failure.
func DefaultOutcomeRouting() OutcomeRouting {
	return OutcomeRouting{Default: OutcomeBlock}
}

// ActionFor returns the outcome action for a given check name.
func (o OutcomeRouting) ActionFor(checkName string) OutcomeAction {
	if action, ok := o.Checks[checkName]; ok {
		return action
	}
	return o.Default
}

// PolicyDocument is the top-level apm-policy.yml structure.
type PolicyDocument struct {
	Version     string
	Extends     []string
	Cache       PolicyCache
	Dependencies DependencyPolicy
	Mcp         McpPolicy
	Compilation CompilationPolicy
	Scripts     ScriptsPolicy
	Outcomes    OutcomeRouting
}

// DefaultPolicyDocument returns a policy document with all defaults.
func DefaultPolicyDocument() PolicyDocument {
	return PolicyDocument{
		Cache:        DefaultPolicyCache(),
		Dependencies: DefaultDependencyPolicy(),
		Outcomes:     DefaultOutcomeRouting(),
	}
}
