// Package tokenmanager provides centralized token management for different AI runtimes
// and git platforms. It handles the complex token environment setup required by
// different AI CLI tools, each of which expects different environment variable names.
package tokenmanager

import (
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

// ADOBearerSource is the diagnostic source label for bearer-resolved tokens (AAD via az CLI).
const ADOBearerSource = "AAD_BEARER_AZ_CLI"

// DefaultCredentialTimeout is the default timeout for git credential fill operations.
const DefaultCredentialTimeout = 60

// MaxCredentialTimeout is the maximum allowed credential timeout.
const MaxCredentialTimeout = 180

// tokenPrecedence defines token precedence for different use cases.
var tokenPrecedence = map[string][]string{
	"copilot":              {"GITHUB_COPILOT_PAT", "GITHUB_TOKEN", "GITHUB_APM_PAT"},
	"models":               {"GITHUB_TOKEN", "GITHUB_APM_PAT"},
	"modules":              {"GITHUB_APM_PAT", "GITHUB_TOKEN", "GH_TOKEN"},
	"gitlab_modules":       {"GITLAB_APM_PAT", "GITLAB_TOKEN"},
	"generic_modules":      {},
	"ado_modules":          {"ADO_APM_PAT"},
	"artifactory_modules":  {"ARTIFACTORY_APM_TOKEN"},
}

// runtimeEnvVars defines runtime-specific environment variable mappings.
var runtimeEnvVars = map[string][]string{
	"copilot": {"GH_TOKEN", "GITHUB_PERSONAL_ACCESS_TOKEN"},
	"codex":   {"GITHUB_TOKEN"},
	"llm":     {"GITHUB_MODELS_KEY"},
}

// GitHubTokenManager manages GitHub token environment setup for different AI runtimes.
type GitHubTokenManager struct {
	PreserveExisting bool
	credentialCache  map[credentialKey]*string
}

type credentialKey struct {
	host string
	port *int
}

// New creates a new GitHubTokenManager.
func New(preserveExisting bool) *GitHubTokenManager {
	return &GitHubTokenManager{
		PreserveExisting: preserveExisting,
		credentialCache:  make(map[credentialKey]*string),
	}
}

// formatCredentialHost embeds a custom port into the git credential host field.
func formatCredentialHost(host string, port *int) string {
	if port != nil {
		return host + ":" + strconv.Itoa(*port)
	}
	return host
}

// sanitizeCredentialPath strips leading /, rejects control chars, allowlists URL schemes.
func sanitizeCredentialPath(path string) string {
	parsed, err := url.Parse(path)
	scheme := ""
	if err == nil {
		scheme = strings.ToLower(parsed.Scheme)
	}
	if scheme != "" {
		allowed := map[string]bool{"https": true, "http": true, "ssh": true}
		if !allowed[scheme] {
			return ""
		}
	}
	var cleaned string
	if scheme != "" && err == nil {
		cleaned = strings.TrimLeft(parsed.Path, "/")
	} else {
		cleaned = strings.TrimLeft(path, "/")
	}
	if cleaned == "" {
		return ""
	}
	for _, ch := range cleaned {
		if ch < 0x20 || ch == 0x7F || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			return ""
		}
	}
	return cleaned
}

// isValidCredentialToken validates that a credential-fill token looks like a real credential.
func isValidCredentialToken(token string) bool {
	if token == "" {
		return false
	}
	if len(token) > 1024 {
		return false
	}
	for _, ch := range []byte(token) {
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			return false
		}
	}
	prompts := []string{"Password for", "Username for", "password for", "username for"}
	for _, p := range prompts {
		if strings.Contains(token, p) {
			return false
		}
	}
	return true
}

// supportsGhCLIHost returns true when host should use gh CLI fallback.
func supportsGhCLIHost(host string) bool {
	if host == "" {
		return false
	}
	if githubhost.IsGitHubHostname(host) {
		return true
	}
	configuredHost := strings.ToLower(githubhost.DefaultHost())
	hostLower := strings.ToLower(host)
	if hostLower != configuredHost {
		return false
	}
	if configuredHost == "github.com" || strings.HasSuffix(configuredHost, ".ghe.com") {
		return false
	}
	if githubhost.IsAzureDevOpsHostname(configuredHost) {
		return false
	}
	return githubhost.IsValidFQDN(configuredHost)
}

// getCredentialTimeout returns the timeout for git credential fill.
func getCredentialTimeout() int {
	raw := strings.TrimSpace(os.Getenv("APM_GIT_CREDENTIAL_TIMEOUT"))
	if raw == "" {
		return DefaultCredentialTimeout
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val < 1 {
		return DefaultCredentialTimeout
	}
	if val > MaxCredentialTimeout {
		return MaxCredentialTimeout
	}
	return val
}

// ResolveCredentialFromGit resolves a credential from the git credential store.
func ResolveCredentialFromGit(host string, port *int, path string) *string {
	hostField := formatCredentialHost(host, port)
	lines := []string{"protocol=https", "host=" + hostField}
	if path != "" {
		sanitized := sanitizeCredentialPath(path)
		if sanitized != "" {
			lines = append(lines, "path="+sanitized)
		}
	}
	stdin := strings.Join(lines, "\n") + "\n\n"

	env := os.Environ()
	env = appendOrReplace(env, "GIT_TERMINAL_PROMPT", "0")
	if runtime.GOOS != "windows" {
		env = appendOrReplace(env, "GIT_ASKPASS", "")
	} else {
		env = appendOrReplace(env, "GIT_ASKPASS", "echo")
	}

	timeout := getCredentialTimeout()
	cmd := exec.Command("git", "credential", "fill")
	cmd.Env = env
	cmd.Stdin = strings.NewReader(stdin)
	done := make(chan struct{})
	var out []byte
	var runErr error
	go func() {
		out, runErr = cmd.Output()
		close(done)
	}()

	timer := make(chan struct{})
	go func() {
		select {
		case <-done:
		case <-timerAfter(timeout):
			cmd.Process.Kill() //nolint:errcheck
			close(timer)
			return
		}
	}()

	select {
	case <-done:
	case <-timer:
		return nil
	}

	if runErr != nil {
		return nil
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "password=") {
			token := line[len("password="):]
			if isValidCredentialToken(token) {
				return &token
			}
			return nil
		}
	}
	return nil
}

// ResolveCredentialFromGhCLI resolves a token from the active gh CLI account for host.
func ResolveCredentialFromGhCLI(host string) *string {
	if !supportsGhCLIHost(host) {
		return nil
	}
	env := os.Environ()
	env = appendOrReplace(env, "GH_PROMPT_DISABLED", "1")
	env = appendOrReplace(env, "GH_NO_UPDATE_NOTIFIER", "1")

	timeout := getCredentialTimeout()
	cmd := exec.Command("gh", "auth", "token", "--hostname", host)
	cmd.Env = env
	cmd.Stdin = strings.NewReader("")
	done := make(chan struct{})
	var out []byte
	var runErr error
	go func() {
		out, runErr = cmd.Output()
		close(done)
	}()

	timer := make(chan struct{})
	go func() {
		select {
		case <-done:
		case <-timerAfter(timeout):
			if cmd.Process != nil {
				cmd.Process.Kill() //nolint:errcheck
			}
			close(timer)
			return
		}
	}()

	select {
	case <-done:
	case <-timer:
		return nil
	}

	if runErr != nil {
		return nil
	}

	token := strings.TrimSpace(string(out))
	if isValidCredentialToken(token) {
		return &token
	}
	return nil
}

// SetupEnvironment sets up the complete token environment for all runtimes.
func (m *GitHubTokenManager) SetupEnvironment(env map[string]string) map[string]string {
	if env == nil {
		env = osEnvMap()
	}
	available := m.getAvailableTokens(env)
	m.setupCopilotTokens(env, available)
	m.setupCodexTokens(env, available)
	m.setupLLMTokens(env, available)
	return env
}

// GetTokenForPurpose gets the best available token for a specific purpose.
func (m *GitHubTokenManager) GetTokenForPurpose(purpose string, env map[string]string) (string, bool) {
	if env == nil {
		env = osEnvMap()
	}
	vars, ok := tokenPrecedence[purpose]
	if !ok {
		return "", false
	}
	for _, v := range vars {
		if t, exists := env[v]; exists && t != "" {
			return t, true
		}
	}
	return "", false
}

// GetTokenWithCredentialFallback gets a token, falling back to git credential helpers.
func (m *GitHubTokenManager) GetTokenWithCredentialFallback(purpose, host string, env map[string]string, port *int) (string, bool) {
	if tok, ok := m.GetTokenForPurpose(purpose, env); ok {
		return tok, true
	}
	key := credentialKey{host: host, port: port}
	if cached, exists := m.credentialCache[key]; exists {
		if cached != nil {
			return *cached, true
		}
		return "", false
	}
	if supportsGhCLIHost(host) {
		if t := ResolveCredentialFromGhCLI(host); t != nil {
			m.credentialCache[key] = t
			return *t, true
		}
	}
	t := ResolveCredentialFromGit(host, port, "")
	m.credentialCache[key] = t
	if t != nil {
		return *t, true
	}
	return "", false
}

// ValidateTokens validates that required tokens are available.
func (m *GitHubTokenManager) ValidateTokens(env map[string]string) (bool, string) {
	if env == nil {
		env = osEnvMap()
	}
	hasAny := false
	for _, purpose := range []string{"copilot", "models", "modules"} {
		if _, ok := m.GetTokenForPurpose(purpose, env); ok {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return false, "No tokens found. Set one of:\n- GITHUB_TOKEN (user-scoped PAT for GitHub Models)\n- GITHUB_APM_PAT (fine-grained PAT for APM modules on GitHub)\n- ADO_APM_PAT (PAT for APM modules on Azure DevOps)"
	}
	if _, ok := m.GetTokenForPurpose("models", env); !ok {
		if env["GITHUB_APM_PAT"] != "" {
			return true, "Warning: Only fine-grained PAT available. GitHub Models requires GITHUB_TOKEN (user-scoped PAT)"
		}
	}
	return true, "Token validation passed"
}

func (m *GitHubTokenManager) getAvailableTokens(env map[string]string) map[string]string {
	tokens := make(map[string]string)
	for _, vars := range tokenPrecedence {
		for _, v := range vars {
			if t, ok := env[v]; ok && t != "" {
				tokens[v] = t
			}
		}
	}
	return tokens
}

func (m *GitHubTokenManager) setupCopilotTokens(env, available map[string]string) {
	tok, ok := m.GetTokenForPurpose("copilot", available)
	if !ok {
		return
	}
	for _, v := range runtimeEnvVars["copilot"] {
		if m.PreserveExisting {
			if _, exists := env[v]; exists {
				continue
			}
		}
		env[v] = tok
	}
}

func (m *GitHubTokenManager) setupCodexTokens(env, available map[string]string) {
	if !(m.PreserveExisting && env["GITHUB_TOKEN"] != "") {
		if tok, ok := m.GetTokenForPurpose("models", available); ok {
			if env["GITHUB_TOKEN"] == "" {
				env["GITHUB_TOKEN"] = tok
			}
		}
	}
	if !(m.PreserveExisting && env["GITHUB_APM_PAT"] != "") {
		if t, ok := available["GITHUB_APM_PAT"]; ok && env["GITHUB_APM_PAT"] == "" {
			env["GITHUB_APM_PAT"] = t
		}
	}
}

func (m *GitHubTokenManager) setupLLMTokens(env, available map[string]string) {
	if m.PreserveExisting && env["GITHUB_MODELS_KEY"] != "" {
		return
	}
	if tok, ok := m.GetTokenForPurpose("models", available); ok {
		env["GITHUB_MODELS_KEY"] = tok
	}
}

// SetupRuntimeEnvironment sets up the complete runtime environment for all AI CLIs.
func SetupRuntimeEnvironment(env map[string]string) map[string]string {
	m := New(true)
	return m.SetupEnvironment(env)
}

// ValidateGitHubTokens validates GitHub token setup.
func ValidateGitHubTokens(env map[string]string) (bool, string) {
	m := New(true)
	return m.ValidateTokens(env)
}

// GetGitHubTokenForRuntime gets the appropriate GitHub token for a specific runtime.
func GetGitHubTokenForRuntime(runtime string, env map[string]string) (string, bool) {
	m := New(true)
	runtimeToPurpose := map[string]string{
		"copilot": "copilot",
		"codex":   "models",
		"llm":     "models",
	}
	purpose, ok := runtimeToPurpose[runtime]
	if !ok {
		return "", false
	}
	return m.GetTokenForPurpose(purpose, env)
}

// osEnvMap returns os.Environ as a map.
func osEnvMap() map[string]string {
	m := make(map[string]string)
	for _, kv := range os.Environ() {
		i := strings.IndexByte(kv, '=')
		if i < 0 {
			continue
		}
		m[kv[:i]] = kv[i+1:]
	}
	return m
}

func appendOrReplace(env []string, key, val string) []string {
	prefix := key + "="
	for i, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			env[i] = prefix + val
			return env
		}
	}
	return append(env, prefix+val)
}
