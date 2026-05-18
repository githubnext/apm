package cloneengine_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/cloneengine"
)

func TestBuildFailureMessage_WithErrors(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("my-dep", "https://github.com/org/repo", []string{"auth failure", "timeout"})
	if msg == "" {
		t.Error("expected non-empty failure message")
	}
	if !strings.Contains(msg, "my-dep") {
		t.Errorf("expected dep name in message, got: %s", msg)
	}
}

func TestBuildFailureMessage_Empty(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("dep", "url", nil)
	if msg == "" {
		t.Error("expected non-empty message even with no errors")
	}
}

func TestDefaultPlanForGitHub_NoToken(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "")
	if len(plan.Attempts) == 0 {
		t.Error("expected at least one attempt in plan")
	}
}

func TestDefaultPlanForGitHub_WithToken(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "mytoken")
	if len(plan.Attempts) == 0 {
		t.Error("expected at least one attempt in plan with token")
	}
	hasHTTPSAttempt := false
	for _, a := range plan.Attempts {
		if a.Kind == cloneengine.AttemptHTTPS {
			hasHTTPSAttempt = true
		}
	}
	if !hasHTTPSAttempt {
		t.Error("expected HTTPS attempt when token provided")
	}
}

func TestDefaultPlanForADO_Basic(t *testing.T) {
	plan := cloneengine.DefaultPlanForADO("org", "project", "repo", "adotoken")
	if len(plan.Attempts) == 0 {
		t.Error("expected at least one attempt for ADO plan")
	}
}

func TestTransportAttempt_Fields(t *testing.T) {
	attempt := cloneengine.TransportAttempt{
		Kind:  cloneengine.AttemptHTTPS,
		URL:   "https://github.com/org/repo.git",
		Label: "https-fallback",
	}
	if attempt.Kind != cloneengine.AttemptHTTPS {
		t.Errorf("Kind: %q", attempt.Kind)
	}
	if attempt.Label != "https-fallback" {
		t.Errorf("Label: %q", attempt.Label)
	}
}

func TestCloneOptions_Fields(t *testing.T) {
	opts := cloneengine.CloneOptions{
		DestDir: "/tmp/dest",
		Verbose: true,
	}
	if opts.DestDir != "/tmp/dest" {
		t.Errorf("DestDir: %q", opts.DestDir)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestClone_FirstAttemptSucceeds(t *testing.T) {
	called := 0
	action := func(url, dest string, env map[string]string) error {
		called++
		return nil
	}
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://example.com/r.git", Label: "first"},
			{Kind: cloneengine.AttemptSSH, URL: "git@example.com:r.git", Label: "second"},
		},
	}
	eng := cloneengine.New(plan, action)
	idx, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 0 {
		t.Errorf("expected idx=0, got %d", idx)
	}
	if called != 1 {
		t.Errorf("expected 1 call, got %d", called)
	}
}

func TestClone_ActionReceivesURL(t *testing.T) {
	var receivedURL string
	action := func(url, dest string, env map[string]string) error {
		receivedURL = url
		return nil
	}
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://expected.com/repo.git", Label: "test"},
		},
	}
	eng := cloneengine.New(plan, action)
	_, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(receivedURL, "expected.com") {
		t.Errorf("unexpected URL: %q", receivedURL)
	}
}

func TestClone_AuthFailureFallsThrough(t *testing.T) {
	// Simulate auth failure on first attempt, success on second
	callN := 0
	action := func(url, dest string, env map[string]string) error {
		callN++
		if callN == 1 {
			return errors.New("remote: Repository not found")
		}
		return nil
	}
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptSSH, URL: "git@github.com:org/repo.git", Label: "ssh"},
			{Kind: cloneengine.AttemptHTTPS, URL: "https://github.com/org/repo.git", Label: "https"},
		},
	}
	eng := cloneengine.New(plan, action)
	idx, err := eng.Clone(cloneengine.CloneOptions{DestDir: "/tmp/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 {
		t.Errorf("expected idx=1, got %d", idx)
	}
}
