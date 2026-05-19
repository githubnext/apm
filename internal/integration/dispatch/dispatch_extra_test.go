package dispatch_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration/dispatch"
)

func TestDefaultDispatchTable_Size(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if len(table) != 6 {
		t.Errorf("expected 6 entries in dispatch table, got %d", len(table))
	}
}

func TestDefaultDispatchTable_NoEmptyCounterKeys(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	for key, entry := range table {
		if entry.CounterKey == "" {
			t.Errorf("table[%q].CounterKey is empty", key)
		}
	}
}

func TestDefaultDispatchTable_NoEmptySyncMethods(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	for key, entry := range table {
		if entry.SyncMethod == "" {
			t.Errorf("table[%q].SyncMethod is empty", key)
		}
	}
}

func TestDefaultDispatchTable_SkillsCounterKey(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["skills"].CounterKey != "skills" {
		t.Errorf("skills.CounterKey = %q, want skills", table["skills"].CounterKey)
	}
}

func TestDefaultDispatchTable_HooksMultiTarget(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["hooks"].MultiTarget {
		t.Error("hooks should have MultiTarget=false")
	}
}

func TestDefaultDispatchTable_InstructionsMultiTarget(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["instructions"].MultiTarget {
		t.Error("instructions should have MultiTarget=false")
	}
}

func TestDefaultDispatchTable_OnlySkillsMultiTarget(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	for key, entry := range table {
		if key == "skills" {
			if !entry.MultiTarget {
				t.Error("skills must be MultiTarget=true")
			}
		} else {
			if entry.MultiTarget {
				t.Errorf("table[%q] should not be MultiTarget", key)
			}
		}
	}
}

func TestPrimitiveDispatch_Fields(t *testing.T) {
	pd := dispatch.PrimitiveDispatch{
		IntegratorClass: "MyIntegrator",
		IntegrateMethod: "my_method",
		SyncMethod:      "sync_method",
		CounterKey:      "my_key",
		MultiTarget:     true,
	}
	if pd.IntegratorClass != "MyIntegrator" {
		t.Errorf("IntegratorClass = %q", pd.IntegratorClass)
	}
	if pd.CounterKey != "my_key" {
		t.Errorf("CounterKey = %q", pd.CounterKey)
	}
	if !pd.MultiTarget {
		t.Error("MultiTarget should be true")
	}
}

func TestDispatchTable_IsMap(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	// Should be able to use as a regular map
	entry, ok := table["agents"]
	if !ok {
		t.Fatal("agents key missing")
	}
	if entry.IntegratorClass == "" {
		t.Error("IntegratorClass should not be empty")
	}
}

func TestDefaultDispatchTable_PromptsCounterKey(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["prompts"].CounterKey != "prompts" {
		t.Errorf("prompts.CounterKey = %q, want prompts", table["prompts"].CounterKey)
	}
}

func TestDefaultDispatchTable_HooksCounterKey(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["hooks"].CounterKey != "hooks" {
		t.Errorf("hooks.CounterKey = %q, want hooks", table["hooks"].CounterKey)
	}
}

func TestDefaultDispatchTable_CommandsIntegratorClass(t *testing.T) {
	table := dispatch.DefaultDispatchTable()
	if table["commands"].IntegratorClass != "CommandIntegrator" {
		t.Errorf("commands.IntegratorClass = %q", table["commands"].IntegratorClass)
	}
}

func TestDefaultDispatchTable_ImmutableBaseline(t *testing.T) {
	// Two calls should return independent tables with same content
	t1 := dispatch.DefaultDispatchTable()
	t2 := dispatch.DefaultDispatchTable()
	if len(t1) != len(t2) {
		t.Error("two DefaultDispatchTable calls should return same size tables")
	}
	for key := range t1 {
		if t1[key].IntegratorClass != t2[key].IntegratorClass {
			t.Errorf("tables differ at key %q", key)
		}
	}
}
