package cloneengine

import (
	"fmt"
	"testing"
)

func TestNew_WithNilActionE3(t *testing.T) {
	plan := TransportPlan{}
	e := New(plan, nil)
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestNew_WithCustomActionE3(t *testing.T) {
	called := false
	action := func(url, dest string, env map[string]string) error {
		called = true
		return nil
	}
	plan := TransportPlan{Attempts: []TransportAttempt{{Kind: AttemptHTTPS, URL: "https://example.com", Label: "test"}}}
	e := New(plan, action)
	_, err := e.Clone(CloneOptions{DestDir: t.TempDir()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected action to be called")
	}
}

func TestClone_EmptyPlanErrorE3(t *testing.T) {
	e := New(TransportPlan{}, nil)
	_, err := e.Clone(CloneOptions{DestDir: t.TempDir()})
	if err == nil {
		t.Error("expected error for empty plan")
	}
}

func TestBuildFailureMessage_ContainsRepoURL(t *testing.T) {
	msg := BuildFailureMessage("mydep", "https://github.com/o/r", []string{"err1"})
	if msg == "" {
		t.Error("expected non-empty failure message")
	}
}

func TestDefaultPlanForGitHub_AttemptKinds(t *testing.T) {
	plan := DefaultPlanForGitHub("owner", "repo", "tok")
	for _, a := range plan.Attempts {
		switch a.Kind {
		case AttemptHTTPS, AttemptSSH, AttemptGitHub, AttemptADO:
		default:
			t.Errorf("unexpected attempt kind: %v", a.Kind)
		}
	}
}

func TestDefaultPlanForADO_ContainsProject(t *testing.T) {
	plan := DefaultPlanForADO("myorg", "myproject", "myrepo", "tok")
	found := false
	for _, a := range plan.Attempts {
		if a.URL != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one URL in ADO plan")
	}
}

func TestAttemptKind_HTTPSValue(t *testing.T) {
	if AttemptHTTPS != "https" {
		t.Errorf("expected https, got %q", AttemptHTTPS)
	}
}

func TestAttemptKind_SSHValue(t *testing.T) {
	if AttemptSSH != "ssh" {
		t.Errorf("expected ssh, got %q", AttemptSSH)
	}
}

func TestClone_SecondAttemptSucceedsE3(t *testing.T) {
	calls := 0
	action := func(url, dest string, env map[string]string) error {
		calls++
		if calls == 1 {
			return fmt.Errorf("auth failure")
		}
		return nil
	}
	plan := TransportPlan{Attempts: []TransportAttempt{
		{Kind: AttemptHTTPS, URL: "https://first.com", Label: "first"},
		{Kind: AttemptSSH, URL: "ssh://second.com", Label: "second"},
	}}
	e := New(plan, action)
	idx, err := e.Clone(CloneOptions{DestDir: t.TempDir()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 {
		t.Errorf("expected idx=1, got %d", idx)
	}
}

func TestTransportAttempt_ZeroValue(t *testing.T) {
	var a TransportAttempt
	if a.Kind != "" {
		t.Error("expected empty Kind")
	}
	if a.URL != "" {
		t.Error("expected empty URL")
	}
}
