// auth.go mirrors src/apm_cli/core/auth.py.
// Provides AuthResolver, HostInfo, AuthContext, and BearerFallbackOutcome.
package core

import (
	"os"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

// HostInfo is an immutable description of a remote Git host.
// Mirrors the HostInfo dataclass in auth.py.
type HostInfo struct {
	Host          string
	Kind          string // "github" | "ghe_cloud" | "ghes" | "ado" | "gitlab" | "generic"
	HasPublicRepos bool
	APIBase       string
	Port          int // 0 = default port
}

// DisplayName returns "host:port" when a non-default port is set, else bare host.
func (h HostInfo) DisplayName() string {
	wellKnown := map[int]bool{443: true, 80: true, 22: true}
	if h.Port != 0 && !wellKnown[h.Port] {
		return h.Host + ":" + itoa(h.Port)
	}
	return h.Host
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	// simple int-to-string without importing strconv at package level
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// AuthContext holds resolved authentication for a single (host, org) pair.
// Mirrors the AuthContext dataclass in auth.py.
type AuthContext struct {
	Token      string // empty string = no token
	Source     string // e.g. "GITHUB_APM_PAT", "none"
	TokenType  string // "fine-grained", "classic", "oauth", "github-app", "unknown"
	HostInfo   HostInfo
	GitEnv     map[string]string
	AuthScheme string // "basic" | "bearer"
}

// BearerFallbackOutcome carries the result of execute_with_bearer_fallback.
type BearerFallbackOutcome struct {
	// Outcome is the final result (caller-defined, stored as interface{}).
	Outcome         interface{}
	BearerAttempted bool
}

// authCacheKey is the map key for the AuthResolver cache.
type authCacheKey struct {
	host string
	port int
	org  string
}

// AuthResolver is the single source of truth for auth resolution.
// Mirrors AuthResolver in auth.py.
type AuthResolver struct {
	tokenManager *GitHubTokenManager
	cache        map[authCacheKey]*AuthContext
	mu           sync.Mutex

	verboseAuthLoggedHosts map[string]bool
	stalePATWarnedHosts    map[string]bool
}

// NewAuthResolver constructs an AuthResolver with a default token manager.
func NewAuthResolver() *AuthResolver {
	return &AuthResolver{
		tokenManager:           NewGitHubTokenManager(),
		cache:                  make(map[authCacheKey]*AuthContext),
		verboseAuthLoggedHosts: make(map[string]bool),
		stalePATWarnedHosts:    make(map[string]bool),
	}
}

// NewAuthResolverWithManager constructs an AuthResolver with a provided token manager.
func NewAuthResolverWithManager(tm *GitHubTokenManager) *AuthResolver {
	r := NewAuthResolver()
	r.tokenManager = tm
	return r
}

// ClassifyHost returns a HostInfo for the given host and port.
// Mirrors AuthResolver.classify_host.
func ClassifyHost(host string, port int) HostInfo {
	h := strings.ToLower(host)

	if h == "github.com" {
		return HostInfo{Host: host, Kind: "github", HasPublicRepos: true,
			APIBase: "https://api.github.com", Port: port}
	}
	if strings.HasSuffix(h, ".ghe.com") {
		return HostInfo{Host: host, Kind: "ghe_cloud", HasPublicRepos: false,
			APIBase: "https://" + host + "/api/v3", Port: port}
	}
	if githubhost.IsAzureDevOpsHostname(host) {
		return HostInfo{Host: host, Kind: "ado", HasPublicRepos: true,
			APIBase: "https://dev.azure.com", Port: port}
	}

	// GHES: GITHUB_HOST is set to a non-github.com, non-ghe.com FQDN.
	ghesHost := strings.ToLower(os.Getenv("GITHUB_HOST"))
	if ghesHost != "" && ghesHost == h &&
		ghesHost != "github.com" && ghesHost != "gitlab.com" &&
		!strings.HasSuffix(ghesHost, ".ghe.com") &&
		githubhost.IsValidFQDN(ghesHost) {
		return HostInfo{Host: host, Kind: "ghes", HasPublicRepos: true,
			APIBase: "https://" + host + "/api/v3", Port: port}
	}

	// GitLab (SaaS + env-configured self-managed) -- after GHES per spec.
	if githubhost.IsGitLabHostname(host) {
		apiBase := "https://gitlab.com/api/v4"
		if h != "gitlab.com" {
			apiBase = "https://" + host + "/api/v4"
		}
		return HostInfo{Host: host, Kind: "gitlab", HasPublicRepos: true,
			APIBase: apiBase, Port: port}
	}

	// Generic FQDN.
	return HostInfo{Host: host, Kind: "generic", HasPublicRepos: true,
		APIBase: "https://" + host + "/api/v3", Port: port}
}

// DetectTokenType classifies a token string by its prefix.
// Mirrors AuthResolver.detect_token_type.
func DetectTokenType(token string) string {
	switch {
	case strings.HasPrefix(token, "github_pat_"):
		return "fine-grained"
	case strings.HasPrefix(token, "ghp_"):
		return "classic"
	case strings.HasPrefix(token, "ghu_"):
		return "oauth"
	case strings.HasPrefix(token, "gho_"):
		return "oauth"
	case strings.HasPrefix(token, "ghs_"):
		return "github-app"
	case strings.HasPrefix(token, "ghr_"):
		return "github-app"
	default:
		return "unknown"
	}
}

// Resolve resolves auth for (host, port, org). Cached and thread-safe.
// Mirrors AuthResolver.resolve.
func (r *AuthResolver) Resolve(host string, org string, port int) *AuthContext {
	hostLower := strings.ToLower(host)
	orgLower := strings.ToLower(org)
	key := authCacheKey{host: hostLower, port: port, org: orgLower}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cached := r.cache[key]; cached != nil {
		return cached
	}

	hostInfo := ClassifyHost(host, port)
	token, source, scheme := r.resolveToken(hostInfo, org)
	tokenType := "unknown"
	if token != "" {
		tokenType = DetectTokenType(token)
	}
	gitEnv := r.buildGitEnv(token, scheme, hostInfo.Kind)

	ctx := &AuthContext{
		Token:      token,
		Source:     source,
		TokenType:  tokenType,
		HostInfo:   hostInfo,
		GitEnv:     gitEnv,
		AuthScheme: scheme,
	}
	r.cache[key] = ctx
	return ctx
}

// purposeForHost maps host kind to token purpose.
func purposeForHost(info HostInfo) string {
	switch info.Kind {
	case "ado":
		return "ado_modules"
	case "gitlab":
		return "gitlab_modules"
	case "generic":
		return "generic_modules"
	default:
		return "modules"
	}
}

// orgToEnvSuffix converts an org name to upper-case env-var suffix with hyphens as underscores.
func orgToEnvSuffix(org string) string {
	return strings.ToUpper(strings.ReplaceAll(org, "-", "_"))
}

// resolveToken walks the token resolution chain. Returns (token, source, scheme).
// Mirrors AuthResolver._resolve_token.
func (r *AuthResolver) resolveToken(info HostInfo, org string) (string, string, string) {
	// ADO: PAT -> none (bearer is fetched lazily in try_with_fallback)
	if info.Kind == "ado" {
		if pat := os.Getenv("ADO_APM_PAT"); pat != "" {
			return pat, "ADO_APM_PAT", "basic"
		}
		return "", "none", "basic"
	}

	// 1. Per-org PAT (GitHub-class only).
	if org != "" && (info.Kind == "github" || info.Kind == "ghe_cloud" || info.Kind == "ghes") {
		envName := "GITHUB_APM_PAT_" + orgToEnvSuffix(org)
		if token := os.Getenv(envName); token != "" {
			return token, envName, "basic"
		}
	}

	// 2. Global env vars by host class.
	purpose := purposeForHost(info)
	env := OSEnvMap()
	if token, ok := r.tokenManager.GetTokenForPurpose(purpose, env); ok {
		source := r.tokenManager.IdentifyEnvSource(purpose)
		return token, source, "basic"
	}

	// 3. gh CLI.
	if token, ok := ResolveCredentialFromGHCLI(info.Host); ok {
		return token, "gh-auth-token", "basic"
	}

	// 4. Git credential helper (not for ADO).
	if info.Kind != "ado" {
		if token, ok := ResolveCredentialFromGit(info.Host, info.Port, ""); ok {
			return token, "git-credential-fill", "basic"
		}
	}

	return "", "none", "basic"
}

// buildGitEnv constructs a process env for git subcommands.
// Mirrors AuthResolver._build_git_env.
func (r *AuthResolver) buildGitEnv(token, scheme, hostKind string) map[string]string {
	env := OSEnvMap()
	env["GIT_TERMINAL_PROMPT"] = "0"
	env["GIT_ASKPASS"] = "echo"

	if scheme == "bearer" && token != "" && hostKind == "ado" {
		delete(env, "GIT_TOKEN")
		for k, v := range githubhost.BuildADOBearerGitEnv(token) {
			env[k] = v
		}
	} else if token != "" {
		env["GIT_TOKEN"] = token
	}
	return env
}

// GitLabRESTHeaders builds HTTP headers for GitLab REST API v4 calls.
// Mirrors AuthResolver.gitlab_rest_headers.
func GitLabRESTHeaders(token string, oauthBearer bool) map[string]string {
	if token == "" {
		return map[string]string{}
	}
	if oauthBearer {
		return map[string]string{"Authorization": "Bearer " + token}
	}
	return map[string]string{"PRIVATE-TOKEN": token}
}
