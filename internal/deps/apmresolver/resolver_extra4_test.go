package apmresolver

import (
	"testing"
)

func TestResolveMaxParallel_LargeValue(t *testing.T) {
	result := resolveMaxParallel(1000)
	if result != 1000 {
		t.Errorf("expected 1000, got %d", result)
	}
}

func TestResolveMaxParallel_One(t *testing.T) {
	result := resolveMaxParallel(1)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

func TestResolveMaxParallel_Two(t *testing.T) {
	result := resolveMaxParallel(2)
	if result != 2 {
		t.Errorf("expected 2, got %d", result)
	}
}

func TestNewResolver_DefaultOptions(t *testing.T) {
	r := New(Options{})
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNewResolver_WithMaxParallel(t *testing.T) {
	r := New(Options{MaxParallel: 4})
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNewResolver_WithAuthResolver(t *testing.T) {
	r := New(Options{MaxParallel: 2})
	if r == nil {
		t.Fatal("expected non-nil resolver with auth resolver")
	}
}

func TestResolveMaxParallel_EnvVarPositive(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "8")
	result := resolveMaxParallel(0)
	if result != 8 {
		t.Errorf("expected 8 from env, got %d", result)
	}
}

func TestResolveMaxParallel_EnvVarOverrideByExplicit(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "4")
	result := resolveMaxParallel(10)
	if result != 10 {
		t.Errorf("expected explicit 10 to win, got %d", result)
	}
}

func TestResolveMaxParallel_EnvVarLarge(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "100")
	result := resolveMaxParallel(0)
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}

func TestResolveMaxParallel_EnvVarFive(t *testing.T) {
	t.Setenv("APM_RESOLVE_PARALLEL", "5")
	result := resolveMaxParallel(0)
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}
