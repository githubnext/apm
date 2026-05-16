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
