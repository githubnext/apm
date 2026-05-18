package cloneengine_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/deps/cloneengine"
)

func TestNew_DefaultAction(t *testing.T) {
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/repo.git", Label: "https"},
		},
	}
	eng := cloneengine.New(plan, nil)
	if eng == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestNew_CustomAction(t *testing.T) {
	called := false
	action := func(url, dest string, env map[string]string) error {
		called = true
		return nil
	}
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptGitHub, URL: "https://github.com/org/repo.git", Label: "github"},
		},
	}
	eng := cloneengine.New(plan, action)
	_, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("custom action should have been called")
	}
}

func TestClone_EmptyPlan(t *testing.T) {
	eng := cloneengine.New(cloneengine.TransportPlan{}, func(url, dest string, env map[string]string) error {
		return nil
	})
	_, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err == nil {
		t.Fatal("expected error for empty plan")
	}
}

func TestClone_AllAttemptsFail(t *testing.T) {
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/a.git", Label: "a"},
			{Kind: cloneengine.AttemptSSH, URL: "ssh://example.com/b.git", Label: "b"},
		},
	}
	eng := cloneengine.New(plan, func(url, dest string, env map[string]string) error {
		return errors.New("transport failure")
	})
	_, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err == nil {
		t.Fatal("expected error when all attempts fail")
	}
}

func TestClone_SecondAttemptSucceeds(t *testing.T) {
	callCount := 0
	action := func(url, dest string, env map[string]string) error {
		callCount++
		if callCount == 1 {
			return errors.New("first attempt fails")
		}
		return nil
	}
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/a.git", Label: "first"},
			{Kind: cloneengine.AttemptGitHub, URL: "https://github.com/org/repo.git", Label: "second"},
		},
	}
	eng := cloneengine.New(plan, action)
	idx, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err != nil {
		t.Fatalf("expected success on second attempt, got: %v", err)
	}
	if idx != 1 {
		t.Fatalf("expected attempt index 1, got %d", idx)
	}
}

func TestAttemptKindConstants(t *testing.T) {
	kinds := []cloneengine.AttemptKind{
		cloneengine.AttemptHTTPS,
		cloneengine.AttemptSSH,
		cloneengine.AttemptGitHub,
		cloneengine.AttemptADO,
	}
	seen := map[cloneengine.AttemptKind]bool{}
	for _, k := range kinds {
		if seen[k] {
			t.Fatalf("duplicate AttemptKind: %s", k)
		}
		seen[k] = true
	}
}
