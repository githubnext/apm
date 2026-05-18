package dispatch_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration/dispatch"
)

func TestDefaultDispatchTable_Keys(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	expected := []string{"prompts", "agents", "commands", "instructions", "hooks", "skills"}
	for _, key := range expected {
		if _, ok := table[key]; !ok {
			t.Errorf("missing key %q in dispatch table", key)
		}
	}
}

func TestDefaultDispatchTable_SkillsMultiTarget(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	skills := table["skills"]
	if !skills.MultiTarget {
		t.Error("skills should have MultiTarget=true")
	}
}

func TestDefaultDispatchTable_PromptsNotMultiTarget(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["prompts"].MultiTarget {
		t.Error("prompts should have MultiTarget=false")
	}
}

func TestDefaultDispatchTable_CounterKeys(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	cases := map[string]string{
		"prompts":      "prompts",
		"agents":       "agents",
		"commands":     "commands",
		"instructions": "instructions",
		"hooks":        "hooks",
		"skills":       "skills",
	}
	for key, wantCounter := range cases {
		if got := table[key].CounterKey; got != wantCounter {
			t.Errorf("table[%q].CounterKey=%q, want %q", key, got, wantCounter)
		}
	}
}

func TestDefaultDispatchTable_IntegratorClasses(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["prompts"].IntegratorClass != "PromptIntegrator" {
		t.Errorf("unexpected: %q", table["prompts"].IntegratorClass)
	}
	if table["skills"].IntegratorClass != "SkillIntegrator" {
		t.Errorf("unexpected: %q", table["skills"].IntegratorClass)
	}
}

func TestDefaultDispatchTable_IntegrateMethods(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	cases := map[string]string{
		"prompts":      "integrate_prompts_for_target",
		"agents":       "integrate_agents_for_target",
		"commands":     "integrate_commands_for_target",
		"instructions": "integrate_instructions_for_target",
		"hooks":        "integrate_hooks_for_target",
		"skills":       "integrate_package_skill",
	}
	for key, want := range cases {
		if got := table[key].IntegrateMethod; got != want {
			t.Errorf("table[%q].IntegrateMethod=%q, want %q", key, got, want)
		}
	}
}

func TestDefaultDispatchTable_SyncMethods(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	perTarget := map[string]bool{"prompts": true, "agents": true, "commands": true, "instructions": true}
	for key, entry := range table {
		if perTarget[key] {
			if entry.SyncMethod != "sync_for_target" {
				t.Errorf("table[%q].SyncMethod=%q, want sync_for_target", key, entry.SyncMethod)
			}
		} else {
			if entry.SyncMethod != "sync_integration" {
				t.Errorf("table[%q].SyncMethod=%q, want sync_integration", key, entry.SyncMethod)
			}
		}
	}
}

func TestDefaultDispatchTable_AllIntegratorClassesSet(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	for key, entry := range table {
		if entry.IntegratorClass == "" {
			t.Errorf("table[%q].IntegratorClass is empty", key)
		}
	}
}

func TestDefaultDispatchTable_AllIntegrateMethodsSet(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	for key, entry := range table {
		if entry.IntegrateMethod == "" {
			t.Errorf("table[%q].IntegrateMethod is empty", key)
		}
	}
}

func TestDefaultDispatchTable_AgentsIntegratorClass(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["agents"].IntegratorClass != "AgentIntegrator" {
		t.Errorf("agents IntegratorClass=%q, want AgentIntegrator", table["agents"].IntegratorClass)
	}
}

func TestDefaultDispatchTable_HooksIntegratorClass(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["hooks"].IntegratorClass != "HookIntegrator" {
		t.Errorf("hooks IntegratorClass=%q, want HookIntegrator", table["hooks"].IntegratorClass)
	}
}

func TestDefaultDispatchTable_InstructionsIntegratorClass(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["instructions"].IntegratorClass != "InstructionIntegrator" {
		t.Errorf("instructions IntegratorClass=%q, want InstructionIntegrator", table["instructions"].IntegratorClass)
	}
}
