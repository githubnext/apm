package dispatch

import "testing"

func TestPrimitiveDispatch_ZeroValue(t *testing.T) {
	pd := PrimitiveDispatch{}
	if pd.IntegratorClass != "" || pd.IntegrateMethod != "" || pd.SyncMethod != "" {
		t.Error("zero value should have empty strings")
	}
	if pd.MultiTarget {
		t.Error("zero value MultiTarget should be false")
	}
}

func TestDefaultDispatchTable_AgentsEntry(t *testing.T) {
	dt := DefaultDispatchTable()
	agents, ok := dt["agents"]
	if !ok {
		t.Fatal("agents entry missing")
	}
	if agents.IntegratorClass != "AgentIntegrator" {
		t.Errorf("unexpected integratorClass: %q", agents.IntegratorClass)
	}
	if agents.MultiTarget {
		t.Error("agents should not be multi-target")
	}
}

func TestDefaultDispatchTable_SkillsMultiTargetTrue(t *testing.T) {
	dt := DefaultDispatchTable()
	skills, ok := dt["skills"]
	if !ok {
		t.Fatal("skills entry missing")
	}
	if !skills.MultiTarget {
		t.Error("skills should be multi-target")
	}
}

func TestDefaultDispatchTable_AllHaveCounterKey(t *testing.T) {
	dt := DefaultDispatchTable()
	for k, v := range dt {
		if v.CounterKey == "" {
			t.Errorf("entry %q has empty CounterKey", k)
		}
	}
}

func TestDefaultDispatchTable_InstructionsIntegratorClass(t *testing.T) {
	dt := DefaultDispatchTable()
	instr, ok := dt["instructions"]
	if !ok {
		t.Fatal("instructions entry missing")
	}
	if instr.IntegratorClass != "InstructionIntegrator" {
		t.Errorf("unexpected class: %q", instr.IntegratorClass)
	}
}

func TestDefaultDispatchTable_PromptsIntegrateMethod(t *testing.T) {
	dt := DefaultDispatchTable()
	prompts, ok := dt["prompts"]
	if !ok {
		t.Fatal("prompts entry missing")
	}
	if prompts.IntegrateMethod != "integrate_prompts_for_target" {
		t.Errorf("unexpected method: %q", prompts.IntegrateMethod)
	}
}

func TestDefaultDispatchTable_IndependentCopies(t *testing.T) {
	dt1 := DefaultDispatchTable()
	dt2 := DefaultDispatchTable()
	dt1["extra"] = PrimitiveDispatch{CounterKey: "extra"}
	if _, ok := dt2["extra"]; ok {
		t.Error("modifying one table should not affect another")
	}
}
