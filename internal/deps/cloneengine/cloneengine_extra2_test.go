package cloneengine_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/cloneengine"
)

func TestClone_AllAttemptsFailExtra2(t *testing.T) {
	plan := cloneengine.TransportPlan{
		Attempts: []cloneengine.TransportAttempt{
			{Kind: cloneengine.AttemptHTTPS, URL: "https://github.com/x/y"},
			{Kind: cloneengine.AttemptSSH, URL: "git@github.com:x/y.git"},
		},
	}
	eng := cloneengine.New(plan, func(url, dest string, env map[string]string) error {
		return errors.New("fail: " + url)
	})
	_, err := eng.Clone(cloneengine.CloneOptions{DestDir: t.TempDir()})
	if err == nil {
		t.Error("expected error when all attempts fail")
	}
}

func TestBuildFailureMessage_ContainsDepName(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("owner/repo", "https://github.com/owner/repo", []string{"auth failed"})
	if !strings.Contains(msg, "owner/repo") {
		t.Errorf("expected dep name in failure message, got %q", msg)
	}
}

func TestBuildFailureMessage_ContainsErrors(t *testing.T) {
	msg := cloneengine.BuildFailureMessage("x/y", "https://x", []string{"err1", "err2"})
	if !strings.Contains(msg, "err1") || !strings.Contains(msg, "err2") {
		t.Errorf("expected errors in failure message, got %q", msg)
	}
}

func TestDefaultPlanForGitHub_AttemptCount(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "")
	if len(plan.Attempts) == 0 {
		t.Error("expected at least one transport attempt")
	}
}

func TestDefaultPlanForGitHub_ContainsHTTPS(t *testing.T) {
	plan := cloneengine.DefaultPlanForGitHub("owner", "repo", "")
	found := false
	for _, a := range plan.Attempts {
		if strings.Contains(a.URL, "https") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected HTTPS attempt in default plan")
	}
}

func TestDefaultPlanForADO_ContainsOrg(t *testing.T) {
	plan := cloneengine.DefaultPlanForADO("myorg", "myproject", "myrepo", "")
	found := false
	for _, a := range plan.Attempts {
		if strings.Contains(a.URL, "myorg") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected org name in ADO plan URL")
	}
}

func TestTransportPlan_EmptyAttempts(t *testing.T) {
	plan := cloneengine.TransportPlan{Attempts: nil}
	if len(plan.Attempts) != 0 {
		t.Errorf("expected 0 attempts, got %d", len(plan.Attempts))
	}
}
