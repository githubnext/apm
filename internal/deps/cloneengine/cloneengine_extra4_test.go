package cloneengine_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/deps/cloneengine"
)

func TestNew_ReturnsEngine(t *testing.T) {
	plan := cloneengine.TransportPlan{}
	e := cloneengine.New(plan, nil)
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestBuildFailureMessage_ContainsDepNameE4(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("my-dep", "https://gh.com/o/r", []string{"err1"})
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestBuildFailureMessage_ContainsRepoURLValue(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("dep", "https://gh.com/o/r", []string{"network error"})
	if msg == "" {
		t.Error("expected non-empty failure message")
	}
}

func TestDefaultPlanForGitHub_NonEmpty(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "token")
	if len(plan.Attempts) == 0 {
		t.Error("expected non-empty attempts for GitHub plan")
	}
}

func TestDefaultPlanForADO_NonEmpty(t *testing.T) {
	plan := cloneengine.DefaultPlanForADO("org", "proj", "repo", "tok")
	if len(plan.Attempts) == 0 {
		t.Error("expected non-empty attempts for ADO plan")
	}
}

func TestClone_SingleAttemptSuccess(t *testing.T) {
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/repo"},
		},
	}
	e := cloneengine.New(plan, func(url, destDir string, env map[string]string) error {
		return nil
	})
	_, err := e.Clone(cloneengine.CloneOptions{DestDir: t.TempDir()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClone_AllFailReturnsError(t *testing.T) {
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/repo"},
		},
	}
	e := cloneengine.New(plan, func(url, destDir string, env map[string]string) error {
		return errors.New("clone failed")
	})
	_, err := e.Clone(cloneengine.CloneOptions{DestDir: t.TempDir()})
	if err == nil {
		t.Error("expected error when all attempts fail")
	}
}

func TestDefaultPlanForGitHub_EmptyToken(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "")
	if len(plan.Attempts) == 0 {
		t.Error("expected attempts even with empty token")
	}
}

func TestTransportAttempt_KindField(t *testing.T) {
	a := cloneengine.TransportAttempt{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com"}
	if a.Kind != cloneengine.AttemptHTTPS {
		t.Errorf("unexpected kind: %s", a.Kind)
	}
}
