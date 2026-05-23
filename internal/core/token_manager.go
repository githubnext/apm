// Package core provides core APM CLI functionality.
// token_manager.go mirrors src/apm_cli/core/token_manager.py.
package core

import (
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

// tokenCacheKey identifies a host+port pair for credential caching.
type tokenCacheKey struct {
	host string
	port int // 0 means no port
}

// TokenPrecedence maps purpose to ordered env var names.
// Mirrors GitHubTokenManager.TOKEN_PRECEDENCE in token_manager.py.
var TokenPrecedence = map[string][]string{
	"copilot":              {"GITHUB_COPILOT_PAT", "GITHUB_TOKEN", "GITHUB_APM_PAT"},
	"models":               {"GITHUB_TOKEN", "GITHUB_APM_PAT"},
	"modules":              {"GITHUB_APM_PAT", "GITHUB_TOKEN", "GH_TOKEN"},
	"gitlab_modules":       {"GITLAB_APM_PAT", "GITLAB_TOKEN"},
	"generic_modules":      {},
	"ado_modules":          {"ADO_APM_PAT"},
	"artifactory_modules":  {"ARTIFACTORY_APM_TOKEN"},
}

// RuntimeEnvVars maps runtime to env var names to set.
var RuntimeEnvVars = map[string][]string{
	"copilot": {"GH_TOKEN", "GITHUB_PERSONAL_ACCESS_TOKEN"},
	"codex":   {"GITHUB_TOKEN"},
	"llm":     {"GITHUB_MODELS_KEY"},
}

const (
	adoBearerSource           = "AAD_BEARER_AZ_CLI"
	defaultCredentialTimeout  = 60
	maxCredentialTimeout      = 180
)

// GitHubTokenManager manages GitHub token environment setup for different AI runtimes.
// Mirrors GitHubTokenManager in token_manager.py.
type GitHubTokenManager struct {
	PreserveExisting bool
	mu               sync.Mutex
	credentialCache  map[tokenCacheKey]*string // *string so nil = "not found", "" = cached none
}

// NewGitHubTokenManager constructs a manager with preserve_existing=true (Python default).
func NewGitHubTokenManager() *GitHubTokenManager {
	return &GitHubTokenManager{
		PreserveExisting: true,
		credentialCache:  make(map[tokenCacheKey]*string),
	}
}

// credentialTimeout returns the timeout in seconds for git credential fill.
// Configurable via APM_GIT_CREDENTIAL_TIMEOUT.
func credentialTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv("APM_GIT_CREDENTIAL_TIMEOUT"))
	if raw == "" {
		return time.Duration(defaultCredentialTimeout) * time.Second
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return time.Duration(defaultCredentialTimeout) * time.Second
	}
	if v < 1 {
		v = 1
	}
	if v > maxCredentialTimeout {
		v = maxCredentialTimeout
	}
	return time.Duration(v) * time.Second
}

// isValidCredentialToken validates that a credential-fill token is not garbage.
func isValidCredentialToken(token string) bool {
	if token == "" || len(token) > 1024 {
		return false
	}
	for _, c := range token {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			return false
		}
	}
	promptFragments := []string{
		"Password for", "Username for", "password for", "username for",
	}
	for _, f := range promptFragments {
		if strings.Contains(token, f) {
			return false
		}
	}
	return true
}

// sanitizeCredentialPath strips leading '/', rejects control chars, allowlists URL schemes.
// Mirrors _sanitize_credential_path in token_manager.py.
func sanitizeCredentialPath(path string) string {
	parsed, err := url.Parse(path)
	if err != nil {
		return ""
	}
	scheme := strings.ToLower(parsed.Scheme)
	var cleaned string
	if scheme != "" {
		if scheme != "https" && scheme != "http" && scheme != "ssh" {
			return ""
		}
		cleaned = strings.TrimLeft(parsed.Path, "/")
	} else {
		cleaned = strings.TrimLeft(path, "/")
	}
	if cleaned == "" {
		return ""
	}
	for _, ch := range cleaned {
		if ch < 0x20 || ch == 0x7F || unicode.IsSpace(ch) {
			return ""
		}
	}
	return cleaned
}

// formatCredentialHost embeds a non-standard port into the host field per gitcredentials(7).
func formatCredentialHost(host string, port int) string {
	if port != 0 {
		return host + ":" + strconv.Itoa(port)
	}
	return host
}

// ResolveCredentialFromGit queries git credential fill for a token.
// Mirrors GitHubTokenManager.resolve_credential_from_git.
func ResolveCredentialFromGit(host string, port int, path string) (string, bool) {
	hostField := formatCredentialHost(host, port)
	lines := []string{"protocol=https", "host=" + hostField}
	if path != "" {
		if sanitized := sanitizeCredentialPath(path); sanitized != "" {
			lines = append(lines, "path="+sanitized)
		}
	}
	stdin := strings.Join(lines, "\n") + "\n\n"

	env := os.Environ()
	env = append(env, "GIT_TERMINAL_PROMPT=0")
	if runtime.GOOS != "windows" {
		env = append(env, "GIT_ASKPASS=")
	} else {
		env = append(env, "GIT_ASKPASS=echo")
	}

	cmd := exec.Command("git", "credential", "fill")
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Env = env

	done := make(chan struct{})
	var out []byte
	var runErr error
	go func() {
		defer close(done)
		out, runErr = cmd.Output()
	}()
	select {
	case <-done:
	case <-time.After(credentialTimeout()):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", false
	}
	if runErr != nil {
		return "", false
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "password=") {
			token := line[len("password="):]
			if isValidCredentialToken(token) {
				return token, true
			}
			return "", false
		}
	}
	return "", false
}

// ResolveCredentialFromGHCLI resolves a token from the active gh CLI account.
// Mirrors GitHubTokenManager.resolve_credential_from_gh_cli.
func ResolveCredentialFromGHCLI(host string) (string, bool) {
	if !githubhost.SupportGHCLIHost(host) {
		return "", false
	}
	env := os.Environ()
	env = append(env, "GH_PROMPT_DISABLED=1", "GH_NO_UPDATE_NOTIFIER=1")

	cmd := exec.Command("gh", "auth", "token", "--hostname", host)
	cmd.Env = env
	cmd.Stdin = nil

	done := make(chan struct{})
	var out []byte
	var runErr error
	go func() {
		defer close(done)
		out, runErr = cmd.Output()
	}()
	select {
	case <-done:
	case <-time.After(credentialTimeout()):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", false
	}
	if runErr != nil {
		return "", false
	}
	token := strings.TrimSpace(string(out))
	if isValidCredentialToken(token) {
		return token, true
	}
	return "", false
}

// GetTokenForPurpose returns the first available env var for the given purpose.
// Mirrors GitHubTokenManager.get_token_for_purpose.
func (m *GitHubTokenManager) GetTokenForPurpose(purpose string, env map[string]string) (string, bool) {
	vars, ok := TokenPrecedence[purpose]
	if !ok {
		return "", false
	}
	for _, v := range vars {
		if token := env[v]; token != "" {
			return token, true
		}
	}
	return "", false
}

// GetTokenWithCredentialFallback tries env vars, then gh CLI, then git credential fill.
// Mirrors GitHubTokenManager.get_token_with_credential_fallback.
func (m *GitHubTokenManager) GetTokenWithCredentialFallback(
	purpose, host string, env map[string]string, port int,
) (string, bool) {
	if token, ok := m.GetTokenForPurpose(purpose, env); ok {
		return token, true
	}

	key := tokenCacheKey{host: host, port: port}
	m.mu.Lock()
	if cached, hit := m.credentialCache[key]; hit {
		m.mu.Unlock()
		if cached == nil || *cached == "" {
			return "", false
		}
		return *cached, true
	}
	m.mu.Unlock()

	var result string
	if githubhost.SupportGHCLIHost(host) {
		if token, ok := ResolveCredentialFromGHCLI(host); ok {
			result = token
		}
	}
	if result == "" {
		if token, ok := ResolveCredentialFromGit(host, port, ""); ok {
			result = token
		}
	}

	m.mu.Lock()
	if result != "" {
		m.credentialCache[key] = &result
	} else {
		empty := ""
		m.credentialCache[key] = &empty
	}
	m.mu.Unlock()

	if result != "" {
		return result, true
	}
	return "", false
}

// OSEnvMap returns os.Environ() as a map[string]string.
func OSEnvMap() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		idx := strings.Index(e, "=")
		if idx < 0 {
			continue
		}
		m[e[:idx]] = e[idx+1:]
	}
	return m
}

// ValidateTokens checks that at least one useful token is available.
// Mirrors GitHubTokenManager.validate_tokens.
func (m *GitHubTokenManager) ValidateTokens(env map[string]string) (bool, string) {
	if env == nil {
		env = OSEnvMap()
	}
	hasAny := false
	for _, purpose := range []string{"copilot", "models", "modules"} {
		if _, ok := m.GetTokenForPurpose(purpose, env); ok {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return false, "No tokens found. Set one of:\n" +
			"- GITHUB_TOKEN (user-scoped PAT for GitHub Models)\n" +
			"- GITHUB_APM_PAT (fine-grained PAT for APM modules on GitHub)\n" +
			"- ADO_APM_PAT (PAT for APM modules on Azure DevOps)"
	}
	if _, ok := m.GetTokenForPurpose("models", env); !ok {
		if env["GITHUB_APM_PAT"] != "" {
			return true, "Warning: Only fine-grained PAT available. GitHub Models requires GITHUB_TOKEN (user-scoped PAT)"
		}
	}
	return true, "Token validation passed"
}

// IdentifyEnvSource returns the name of the first env var that matched for purpose.
// Mirrors AuthResolver._identify_env_source.
func (m *GitHubTokenManager) IdentifyEnvSource(purpose string) string {
	for _, v := range TokenPrecedence[purpose] {
		if os.Getenv(v) != "" {
			return v
		}
	}
	return "env"
}

// ADOBearerSource is the diagnostic source label for ADO bearer tokens.
const ADOBearerSource = adoBearerSource
