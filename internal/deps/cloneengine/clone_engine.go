// Package cloneengine drives a transport-plan-driven clone execution.
//
// Each TransportAttempt is a self-contained recipe (URL scheme, auth
// scheme, label) that the engine renders into a concrete URL + git env,
// hands to the caller-provided clone action, and -- on auth/transport
// failure -- rolls forward to the next attempt or applies an in-attempt
// ADO bearer fallback.
//
// Migrated from: src/apm_cli/deps/clone_engine.py
package cloneengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AttemptKind identifies the transport / auth scheme for one clone attempt.
type AttemptKind string

const (
	AttemptHTTPS  AttemptKind = "https"
	AttemptSSH    AttemptKind = "ssh"
	AttemptGitHub AttemptKind = "github-token"
	AttemptADO    AttemptKind = "ado-token"
)

// TransportAttempt is one self-contained clone recipe.
type TransportAttempt struct {
	Kind    AttemptKind
	URL     string
	Token   string
	Label   string
	GitEnv  map[string]string
}

// TransportPlan is an ordered sequence of TransportAttempts to try.
type TransportPlan struct {
	Attempts []TransportAttempt
}

// CloneAction is the function the engine calls to perform a git clone.
type CloneAction func(url, destDir string, env map[string]string) error

// CloneOptions configures one engine run.
type CloneOptions struct {
	// DestDir is the directory to clone into.
	DestDir string
	// Verbose enables progress output.
	Verbose bool
	// Timeout is the per-attempt timeout (0 = no limit).
	// Depth limits the clone depth (0 = full).
	Depth int
	// Branch clones a specific branch.
	Branch string
}

// CloneEngine drives a TransportPlan to completion.
type CloneEngine struct {
	plan   TransportPlan
	action CloneAction
}

// New creates a CloneEngine. If action is nil, a default git-based action is used.
func New(plan TransportPlan, action CloneAction) *CloneEngine {
	if action == nil {
		action = defaultCloneAction
	}
	return &CloneEngine{plan: plan, action: action}
}

// Clone tries each attempt in order until one succeeds.
// Returns the index of the successful attempt and nil, or an error wrapping
// all attempt errors.
func (e *CloneEngine) Clone(opts CloneOptions) (int, error) {
	if len(e.plan.Attempts) == 0 {
		return -1, errors.New("no transport attempts in plan")
	}

	var errs []string
	for i, attempt := range e.plan.Attempts {
		if opts.Verbose {
			fmt.Printf("[>] Trying %s (%s)\n", attempt.Label, attempt.Kind)
		}

		url := attempt.URL
		env := mergeEnv(attempt.GitEnv)

		// Inject token into HTTPS URL or via git credential helper env.
		if attempt.Token != "" {
			switch attempt.Kind {
			case AttemptHTTPS, AttemptGitHub, AttemptADO:
				url = injectToken(url, attempt.Token)
			case AttemptSSH:
				env["GIT_SSH_COMMAND"] = sshCommand(attempt.Token)
			}
		}

		if err := e.action(url, opts.DestDir, env); err != nil {
			errs = append(errs, fmt.Sprintf("attempt %d (%s): %v", i+1, attempt.Label, err))
			if isAuthFailure(err) {
				continue // try next attempt
			}
			// Non-auth failure: try ADO bearer fallback if applicable.
			if attempt.Kind == AttemptADO && attempt.Token == "" {
				if bearerToken := os.Getenv("AZURE_ACCESS_TOKEN"); bearerToken != "" {
					env2 := mergeEnv(env)
					env2["GIT_TOKEN"] = bearerToken
					if err2 := e.action(url, opts.DestDir, env2); err2 == nil {
						return i, nil
					}
				}
			}
			continue
		}
		return i, nil
	}

	return -1, fmt.Errorf("all %d transport attempts failed:\n  %s",
		len(e.plan.Attempts), strings.Join(errs, "\n  "))
}

// BuildFailureMessage formats a human-readable clone failure message.
func BuildFailureMessage(depName, repoURL string, errs []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Failed to clone %s from %s.\n", depName, repoURL))
	sb.WriteString("Tried the following transports:\n")
	for _, e := range errs {
		sb.WriteString(fmt.Sprintf("  - %s\n", e))
	}
	sb.WriteString("\nCommon fixes:\n")
	sb.WriteString("  - Ensure GITHUB_TOKEN (or GH_TOKEN) is set for private repos\n")
	sb.WriteString("  - For ADO repos, set AZURE_ACCESS_TOKEN\n")
	sb.WriteString("  - Check network / firewall restrictions\n")
	return sb.String()
}

// ---------------------------------------------------------
// Transport plan builders
// ---------------------------------------------------------

// DefaultPlanForGitHub builds the standard HTTPS + SSH transport plan for GitHub.
func DefaultPlanForGitHub(owner, repo, token string) TransportPlan {
	httpsURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	sshURL := fmt.Sprintf("git@github.com:%s/%s.git", owner, repo)

	attempts := []TransportAttempt{
		{Kind: AttemptHTTPS, URL: httpsURL, Token: token, Label: "HTTPS+token"},
		{Kind: AttemptSSH, URL: sshURL, Label: "SSH"},
		{Kind: AttemptHTTPS, URL: httpsURL, Label: "HTTPS (unauthenticated)"},
	}
	return TransportPlan{Attempts: attempts}
}

// DefaultPlanForADO builds the ADO transport plan.
func DefaultPlanForADO(org, project, repo, token string) TransportPlan {
	httpsURL := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", org, project, repo)
	attempts := []TransportAttempt{
		{Kind: AttemptADO, URL: httpsURL, Token: token, Label: "ADO HTTPS+token"},
		{Kind: AttemptADO, URL: httpsURL, Label: "ADO bearer fallback"},
	}
	return TransportPlan{Attempts: attempts}
}

// ---------------------------------------------------------
// Helpers
// ---------------------------------------------------------

func defaultCloneAction(url, destDir string, env map[string]string) error {
	args := []string{"clone", "--depth=1", url, destDir}
	cmd := exec.Command("git", args...)
	cmd.Env = buildEnv(env)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func injectToken(rawURL, token string) string {
	if token == "" {
		return rawURL
	}
	// Insert token as x-access-token user in the URL.
	for _, scheme := range []string{"https://", "http://"} {
		if strings.HasPrefix(rawURL, scheme) {
			return scheme + "x-access-token:" + token + "@" + rawURL[len(scheme):]
		}
	}
	return rawURL
}

func sshCommand(token string) string {
	// For SSH key-based auth the token is treated as a key path.
	if filepath.IsAbs(token) || strings.HasPrefix(token, "~") {
		return fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no", token)
	}
	return "ssh -o StrictHostKeyChecking=no"
}

func isAuthFailure(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, sig := range []string{
		"authentication failed",
		"could not read username",
		"invalid credentials",
		"403",
		"401",
		"permission denied (publickey)",
	} {
		if strings.Contains(msg, sig) {
			return true
		}
	}
	return false
}

func mergeEnv(extra map[string]string) map[string]string {
	out := make(map[string]string, len(extra))
	for k, v := range extra {
		out[k] = v
	}
	return out
}

func buildEnv(extra map[string]string) []string {
	env := os.Environ()
	for k, v := range extra {
		env = append(env, k+"="+v)
	}
	return env
}
