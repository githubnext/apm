// Package dispatch defines the primitive dispatch registry.
// Maps each APM primitive type to its integrator class and integration methods.
// Mirrors src/apm_cli/integration/dispatch.py.
package dispatch

// PrimitiveDispatch describes how to integrate a single primitive type.
type PrimitiveDispatch struct {
	// IntegratorClass is the name of the integrator (used as a reference key).
	IntegratorClass string

	// IntegrateMethod is the method name for install (per-target or all-targets).
	IntegrateMethod string

	// SyncMethod is the method name for uninstall/removal.
	SyncMethod string

	// CounterKey is the key in the result counters dict (e.g., "agents").
	CounterKey string

	// MultiTarget indicates the integrator receives all targets at once.
	// Used by SkillIntegrator.
	MultiTarget bool
}

// DispatchTable maps primitive names to their dispatch configuration.
type DispatchTable map[string]PrimitiveDispatch

// DefaultDispatchTable returns the standard primitive dispatch table.
// This mirrors the _build_dispatch() function in the Python implementation.
func DefaultDispatchTable() DispatchTable {
	return DispatchTable{
		"prompts": {
			IntegratorClass: "PromptIntegrator",
			IntegrateMethod: "integrate_prompts_for_target",
			SyncMethod:      "sync_for_target",
			CounterKey:      "prompts",
			MultiTarget:     false,
		},
		"agents": {
			IntegratorClass: "AgentIntegrator",
			IntegrateMethod: "integrate_agents_for_target",
			SyncMethod:      "sync_for_target",
			CounterKey:      "agents",
			MultiTarget:     false,
		},
		"commands": {
			IntegratorClass: "CommandIntegrator",
			IntegrateMethod: "integrate_commands_for_target",
			SyncMethod:      "sync_for_target",
			CounterKey:      "commands",
			MultiTarget:     false,
		},
		"instructions": {
			IntegratorClass: "InstructionIntegrator",
			IntegrateMethod: "integrate_instructions_for_target",
			SyncMethod:      "sync_for_target",
			CounterKey:      "instructions",
			MultiTarget:     false,
		},
		"hooks": {
			IntegratorClass: "HookIntegrator",
			IntegrateMethod: "integrate_hooks_for_target",
			SyncMethod:      "sync_integration",
			CounterKey:      "hooks",
			MultiTarget:     false,
		},
		"skills": {
			IntegratorClass: "SkillIntegrator",
			IntegrateMethod: "integrate_package_skill",
			SyncMethod:      "sync_integration",
			CounterKey:      "skills",
			MultiTarget:     true,
		},
	}
}
