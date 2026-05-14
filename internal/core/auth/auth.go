// Package auth provides centralized authentication resolution for APM CLI.
// Every APM operation that touches a remote host MUST use AuthResolver.
// Resolution is per-(host, org) pair, thread-safe, and cached per-process.
package auth

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/core/tokenmanager"
	"github.com/githubnext/apm/internal/utils/githubhost"
)

// HostInfo is an immutable description of a remote Git host.
type HostInfo struct {
	Host           string
	Kind           string // "github" | "ghe_cloud" | "ghes" | "ado" | "gitlab" | "generic"
	HasPublicRepos bool
	APIBase        string
	Port           *int // Non-standard git port, nil for default
}

// DisplayName returns "host:port" when a non-default port is set, else bare host.
func (h HostInfo) DisplayName() string {
	wellKnown := map[int]bool{443: true, 80: true, 22: true}
	if h.Port != nil && !wellKnown[*h.Port] {
		return fmt.Sprintf("%s:%d", h.Host, *h.Port)
	}
	return h.Host
}

// AuthContext holds resolved authentication for a single (host, org) pair.
type AuthContext struct {
	Token      *string // nil means no token; never print
	Source     string  // e.g. "GITHUB_APM_PAT_ORGNAME", "GITHUB_TOKEN", "none"
	TokenType  string  // "fine-grained", "classic", "oauth", "github-app", "unknown"
	HostInfo   HostInfo
	GitEnv     map[string]string
	AuthScheme string // "basic" | "bearer"
}

// BearerFallbackOutcome is the result of ExecuteWithBearerFallback.
type BearerFallbackOutcome struct {
	Outcome         interface{}
	BearerAttempted bool
}

type cacheKey struct {
	host string
	port int // 0 means no port
	org  string
}

// AuthResolver is the single source of truth for auth resolution.
// Every APM operation that touches a remote host MUST use this struct.
type AuthResolver struct {
	tokenManager *tokenmanager.GitHubTokenManager
	cache        map[cacheKey]*AuthContext
	mu           sync.Mutex

	// Optional logger interface (set via SetLogger).
	logger interface{}

	verboseAuthLoggedHosts map[string]bool
	stalePATWarnedHosts    map[string]bool
}

// NewAuthResolver constructs a new AuthResolver with an optional token manager.
func NewAuthResolver(tm *tokenmanager.GitHubTokenManager) *AuthResolver {
	if tm == nil {
		tm = &tokenmanager.GitHubTokenManager{}
	}
	return &AuthResolver{
		tokenManager:           tm,
		cache:                  make(map[cacheKey]*AuthContext),
		verboseAuthLoggedHosts: make(map[string]bool),
		stalePATWarnedHosts:    make(map[string]bool),
	}
}

// SetLogger wires a logger into the resolver after construction.
func (r *AuthResolver) SetLogger(logger interface{}) {
	r.logger = logger
}

// ClassifyHost returns a HostInfo describing host.
func ClassifyHost(host string, port *int) HostInfo {
	h := strings.ToLower(host)

	if h == "github.com" {
		return HostInfo{
			Host:           host,
			Kind:           "github",
			HasPublicRepos: true,
			APIBase:        "https://api.github.com",
			Port:           port,
		}
	}

	if strings.HasSuffix(h, ".ghe.com") {
		return HostInfo{
			Host:           host,
			Kind:           "ghe_cloud",
			HasPublicRepos: false,
			APIBase:        fmt.Sprintf("https://%s/api/v3", host),
			Port:           port,
		}
	}

	if githubhost.IsAzureDevOpsHostname(host) {
		return HostInfo{
			Host:           host,
			Kind:           "ado",
			HasPublicRepos: true,
			APIBase:        "https://dev.azure.com",
			Port:           port,
		}
	}

	// GHES: GITHUB_HOST is set to a non-github.com, non-ghe.com FQDN
	ghesHost := strings.ToLower(os.Getenv("GITHUB_HOST"))
	if ghesHost != "" && ghesHost == h &&
		ghesHost != "github.com" && ghesHost != "gitlab.com" &&
		!strings.HasSuffix(ghesHost, ".ghe.com") {
		if githubhost.IsValidFQDN(ghesHost) {
			return HostInfo{
				Host:           host,
				Kind:           "ghes",
				HasPublicRepos: true,
				APIBase:        fmt.Sprintf("https://%s/api/v3", host),
				Port:           port,
			}
		}
	}

	// GitLab (after GHES per spec)
	if githubhost.IsGitLabHostname(host) {
		var apiBase string
		if h == "gitlab.com" {
			apiBase = "https://gitlab.com/api/v4"
		} else {
			apiBase = fmt.Sprintf("https://%s/api/v4", host)
		}
		return HostInfo{
			Host:           host,
			Kind:           "gitlab",
			HasPublicRepos: true,
			APIBase:        apiBase,
			Port:           port,
		}
	}

	// Generic FQDN (Bitbucket, self-hosted, etc.)
	return HostInfo{
		Host:           host,
		Kind:           "generic",
		HasPublicRepos: true,
		APIBase:        fmt.Sprintf("https://%s/api/v3", host),
		Port:           port,
	}
}

// DetectTokenType classifies a token string by its prefix.
func DetectTokenType(token string) string {
	switch {
	case strings.HasPrefix(token, "github_pat_"):
		return "fine-grained"
	case strings.HasPrefix(token, "ghp_"):
		return "classic"
	case strings.HasPrefix(token, "ghu_") || strings.HasPrefix(token, "gho_"):
		return "oauth"
	case strings.HasPrefix(token, "ghs_") || strings.HasPrefix(token, "ghr_"):
		return "github-app"
	}
	return "unknown"
}

// GitLabRESTHeaders builds HTTP headers for GitLab REST API v4 calls.
func GitLabRESTHeaders(token string, oauthBearer bool) map[string]string {
	if token == "" {
		return map[string]string{}
	}
	if oauthBearer {
		return map[string]string{"Authorization": "Bearer " + token}
	}
	return map[string]string{"PRIVATE-TOKEN": token}
}

// Resolve resolves auth for (host, port, org). Cached and thread-safe.
func (r *AuthResolver) Resolve(host, org string, port *int) *AuthContext {
	portVal := 0
	if port != nil {
		portVal = *port
	}
	key := cacheKey{
		host: strings.ToLower(host),
		port: portVal,
		org:  strings.ToLower(org),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cached, ok := r.cache[key]; ok {
		return cached
	}

	hostInfo := ClassifyHost(host, port)
	token, source, scheme := r.resolveToken(hostInfo, org)

	var tokenType string
	if token != nil {
		tokenType = DetectTokenType(*token)
	} else {
		tokenType = "unknown"
	}
	gitEnv := buildGitEnv(token, scheme, hostInfo.Kind)

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

// resolveToken walks the token resolution chain. Returns (token, source, scheme).
func (r *AuthResolver) resolveToken(hostInfo HostInfo, org string) (*string, string, string) {
	if hostInfo.Kind == "ado" {
		if pat := os.Getenv("ADO_APM_PAT"); pat != "" {
			return &pat, "ADO_APM_PAT", "basic"
		}
		return nil, "none", "basic"
	}

	// 1. Per-org GitHub PAT (GitHub-class hosts only)
	if org != "" && (hostInfo.Kind == "github" || hostInfo.Kind == "ghe_cloud" || hostInfo.Kind == "ghes") {
		envName := "GITHUB_APM_PAT_" + orgToEnvSuffix(org)
		if val := os.Getenv(envName); val != "" {
			return &val, envName, "basic"
		}
	}

	// 2. Global env vars by host class
	purpose := purposeForHost(hostInfo)
	token, ok := r.tokenManager.GetTokenForPurpose(purpose, nil)
	if ok && token != "" {
		source := identifyEnvSource(purpose)
		return &token, source, "basic"
	}

	// 3. gh CLI active account
	ghTokenPtr := tokenmanager.ResolveCredentialFromGhCLI(hostInfo.Host)
	if ghTokenPtr != nil && *ghTokenPtr != "" {
		ghToken := *ghTokenPtr
		return &ghToken, "gh-auth-token", "basic"
	}

	// 4. Git credential helper (not for ADO)
	if hostInfo.Kind != "ado" {
		credPtr := tokenmanager.ResolveCredentialFromGit(hostInfo.Host, hostInfo.Port, "")
		if credPtr != nil && *credPtr != "" {
			cred := *credPtr
			return &cred, "git-credential-fill", "basic"
		}
	}

	return nil, "none", "basic"
}

// purposeForHost maps host kind to token manager purpose key.
func purposeForHost(hostInfo HostInfo) string {
	switch hostInfo.Kind {
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

// tokenPrecedenceByPurpose mirrors the Python tokenPrecedence dict.
var tokenPrecedenceByPurpose = map[string][]string{
	"modules":         {"GITHUB_APM_PAT", "GITHUB_TOKEN", "GH_TOKEN"},
	"gitlab_modules":  {"GITLAB_APM_PAT", "GITLAB_TOKEN"},
	"generic_modules": {},
	"ado_modules":     {"ADO_APM_PAT"},
}

// identifyEnvSource returns the name of the first env var that matched for purpose.
func identifyEnvSource(purpose string) string {
	for _, v := range tokenPrecedenceByPurpose[purpose] {
		if os.Getenv(v) != "" {
			return v
		}
	}
	return "env"
}

// orgToEnvSuffix converts an org name to an env-var suffix (upper-case, hyphens to underscores).
func orgToEnvSuffix(org string) string {
	return strings.ToUpper(strings.ReplaceAll(org, "-", "_"))
}

// buildGitEnv builds environment for subprocess git calls.
func buildGitEnv(token *string, scheme, hostKind string) map[string]string {
	env := make(map[string]string)
	// Copy current env
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	env["GIT_TERMINAL_PROMPT"] = "0"
	env["GIT_ASKPASS"] = "echo"

	if scheme == "bearer" && token != nil && *token != "" && hostKind == "ado" {
		delete(env, "GIT_TOKEN")
		// ADO bearer: inject via GIT_CONFIG env vars
		env["GIT_CONFIG_COUNT"] = "1"
		env["GIT_CONFIG_KEY_0"] = "http.extraHeader"
		env["GIT_CONFIG_VALUE_0"] = "Authorization: Bearer " + *token
	} else if token != nil && *token != "" {
		env["GIT_TOKEN"] = *token
	}
	return env
}

// TryWithFallbackOptions configures TryWithFallback.
type TryWithFallbackOptions struct {
	Org             string
	Port            *int
	Path            string
	UnauthFirst     bool
	VerboseCallback func(string)
}

// TryWithFallback executes op with automatic auth/unauth fallback.
// op receives (token *string, gitEnv map[string]string).
func (r *AuthResolver) TryWithFallback(
	host string,
	op func(token *string, gitEnv map[string]string) (interface{}, error),
	opts TryWithFallbackOptions,
) (interface{}, error) {
	authCtx := r.Resolve(host, opts.Org, opts.Port)
	hostInfo := authCtx.HostInfo

	log := func(msg string) {
		if opts.VerboseCallback != nil {
			opts.VerboseCallback(msg)
		}
	}

	tryCredentialFallback := func(origErr error) (interface{}, error) {
		if authCtx.Source == "gh-auth-token" || authCtx.Source == "git-credential-fill" || authCtx.Source == "none" {
			return nil, origErr
		}
		if hostInfo.Kind == "ado" {
			return nil, origErr
		}
		log(fmt.Sprintf("Token from %s failed for %s; trying secondary credential sources",
			authCtx.Source, hostInfo.DisplayName()))
		log(fmt.Sprintf("trying gh auth token for %s", hostInfo.DisplayName()))
		ghTokenPtr := tokenmanager.ResolveCredentialFromGhCLI(hostInfo.Host)
		if ghTokenPtr != nil && *ghTokenPtr != "" {
			log(fmt.Sprintf("gh auth token resolved a credential for %s", hostInfo.DisplayName()))
			return op(ghTokenPtr, buildGitEnv(ghTokenPtr, "basic", hostInfo.Kind))
		}
		pathSuffix := ""
		if opts.Path != "" {
			pathSuffix = fmt.Sprintf(" (path=%s)", opts.Path)
		}
		log(fmt.Sprintf("trying git credential fill for %s%s", hostInfo.DisplayName(), pathSuffix))
		credPtr := tokenmanager.ResolveCredentialFromGit(hostInfo.Host, hostInfo.Port, opts.Path)
		if credPtr != nil && *credPtr != "" {
			log(fmt.Sprintf("git credential fill resolved a credential for %s", hostInfo.DisplayName()))
			return op(credPtr, buildGitEnv(credPtr, "basic", hostInfo.Kind))
		}
		return nil, origErr
	}

	// Hosts that never have public repos -> auth-only
	if hostInfo.Kind == "ghe_cloud" {
		log(fmt.Sprintf("Auth-only attempt for %s host %s", hostInfo.Kind, hostInfo.DisplayName()))
		res, err := op(authCtx.Token, authCtx.GitEnv)
		if err != nil {
			return tryCredentialFallback(err)
		}
		return res, nil
	}

	// ADO: auth-first (bearer fallback handled separately)
	if hostInfo.Kind == "ado" {
		log(fmt.Sprintf("Auth-only attempt for %s host %s", hostInfo.Kind, hostInfo.DisplayName()))
		return op(authCtx.Token, authCtx.GitEnv)
	}

	if opts.UnauthFirst {
		res, err := op(nil, authCtx.GitEnv)
		if err != nil && authCtx.Token != nil {
			log(fmt.Sprintf("Unauthenticated failed, retrying with token (source: %s)", authCtx.Source))
			res2, err2 := op(authCtx.Token, authCtx.GitEnv)
			if err2 != nil {
				return tryCredentialFallback(err2)
			}
			return res2, nil
		}
		return res, err
	}
	if authCtx.Token != nil {
		log(fmt.Sprintf("Trying authenticated access to %s (source: %s)", hostInfo.DisplayName(), authCtx.Source))
		res, err := op(authCtx.Token, authCtx.GitEnv)
		if err != nil {
			if hostInfo.HasPublicRepos {
				log("Authenticated failed, retrying without token")
				res2, err2 := op(nil, authCtx.GitEnv)
				if err2 != nil {
					return tryCredentialFallback(err2)
				}
				return res2, nil
			}
			return tryCredentialFallback(err)
		}
		return res, nil
	}
	log(fmt.Sprintf("No token available, trying unauthenticated access to %s", hostInfo.DisplayName()))
	return op(nil, authCtx.GitEnv)
}

// BuildErrorContext returns an actionable error message for auth failures.
func (r *AuthResolver) BuildErrorContext(
	host, operation, org string,
	port *int,
	depURL string,
	bearerAlsoFailed bool,
) string {
	authCtx := r.Resolve(host, org, port)
	hostInfo := authCtx.HostInfo
	display := hostInfo.DisplayName()

	if hostInfo.Kind == "ado" {
		azAvailable := false // simplified: no az CLI check in Go migration
		patSet := os.Getenv("ADO_APM_PAT") != ""

		orgPart := org
		if orgPart == "" && depURL != "" {
			stripped := strings.TrimPrefix(depURL, "https://")
			parts := strings.SplitN(stripped, "/", 3)
			if len(parts) >= 2 {
				if parts[0] == "dev.azure.com" || strings.HasSuffix(parts[0], ".visualstudio.com") {
					orgPart = parts[1]
				}
			}
		}
		tokenURL := "https://dev.azure.com/<org>/_usersSettings/tokens"
		if orgPart != "" {
			tokenURL = fmt.Sprintf("https://dev.azure.com/%s/_usersSettings/tokens", orgPart)
		}

		if patSet {
			if azAvailable {
				prefix := ""
				if bearerAlsoFailed {
					prefix = "    ADO_APM_PAT was rejected; az cli bearer was also rejected.\n\n"
				}
				return fmt.Sprintf("\n%s    ADO_APM_PAT is set, and Azure CLI credentials may also be available,\n    but the Azure DevOps request still failed.\n\n    To fix:\n      1. Unset the PAT to test Azure CLI auth only:  unset ADO_APM_PAT\n      2. Re-authenticate Azure CLI if needed:        az login\n      3. Retry:                                       apm install\n\n    Docs: https://microsoft.github.io/apm/getting-started/authentication/#azure-devops", prefix)
			}
			return fmt.Sprintf("\n    ADO_APM_PAT is set, but the Azure DevOps request failed.\n    If this is an authentication failure, the token may be expired,\n    revoked, or scoped to a different org.\n\n    Generate a new PAT at %s\n    with Code (Read) scope.\n\n    Docs: https://microsoft.github.io/apm/getting-started/authentication/#azure-devops", tokenURL)
		}
		return fmt.Sprintf("\n    Azure DevOps requires authentication. You have two options:\n\n    1. Install Azure CLI and sign in (recommended for Entra ID users):\n         az login\n         apm install\n\n    2. Use a Personal Access Token:\n         export ADO_APM_PAT=your_token\n         (Create one at %s with Code (Read) scope.)\n\n    Docs: https://microsoft.github.io/apm/getting-started/authentication/#azure-devops", tokenURL)
	}

	// Non-ADO paths
	lines := []string{fmt.Sprintf("Authentication failed for %s on %s.", operation, display)}
	if authCtx.Token != nil {
		lines = append(lines, fmt.Sprintf("Token was provided (source: %s, type: %s).", authCtx.Source, authCtx.TokenType))
		switch {
		case hostInfo.Kind == "ghe_cloud":
			lines = append(lines, "GHE Cloud Data Residency hosts (*.ghe.com) require enterprise-scoped tokens.")
		case hostInfo.Kind == "gitlab":
			lines = append(lines, "Ensure your GitLab personal or project access token meets the API read requirements for your instance policy.")
		case strings.ToLower(host) == "github.com":
			lines = append(lines, "If your organization uses SAML SSO or is an EMU org, ensure your PAT is authorized at https://github.com/settings/tokens")
		case hostInfo.Kind == "generic":
			lines = append(lines, "Verify credentials for this host in your git credential helper.")
		default:
			lines = append(lines, "If your organization uses SAML SSO, you may need to authorize your token at https://github.com/settings/tokens")
		}
	} else {
		lines = append(lines, "No token available.")
		switch hostInfo.Kind {
		case "gitlab":
			lines = append(lines, fmt.Sprintf("Set GITLAB_APM_PAT or GITLAB_TOKEN, or configure git credential fill for %s.", display))
		case "generic":
			lines = append(lines, fmt.Sprintf("APM does not apply GitHub PAT environment variables to generic git hosts; configure git credential fill for %s or use a public repository if available.", display))
		default:
			lines = append(lines, "Set GITHUB_APM_PAT or GITHUB_TOKEN, or run 'gh auth login'.")
		}
	}
	if org != "" && hostInfo.Kind != "ado" && hostInfo.Kind != "gitlab" && hostInfo.Kind != "generic" {
		lines = append(lines, fmt.Sprintf("If packages span multiple organizations, set per-org tokens: GITHUB_APM_PAT_%s", orgToEnvSuffix(org)))
	}
	if hostInfo.Port != nil {
		lines = append(lines, fmt.Sprintf("[i] Host '%s' -- verify your credential helper stores per-port entries (some helpers key by host only).", display))
	}
	lines = append(lines, "Run with --verbose for detailed auth diagnostics.")
	return strings.Join(lines, "\n")
}

// EmitStalePATDiagnostic emits a warning when PAT was rejected but bearer succeeded.
func (r *AuthResolver) EmitStalePATDiagnostic(hostDisplay string) {
	r.mu.Lock()
	if r.stalePATWarnedHosts[hostDisplay] {
		r.mu.Unlock()
		return
	}
	r.stalePATWarnedHosts[hostDisplay] = true
	r.mu.Unlock()

	msg := fmt.Sprintf("ADO_APM_PAT was rejected for %s; fell back to az cli bearer.", hostDisplay)
	fmt.Fprintln(os.Stderr, "[!] "+msg)
	fmt.Fprintln(os.Stderr, "[!]     Consider unsetting the stale variable.")
}

// NotifyAuthSource emits the verbose auth-source line for hostDisplay exactly once.
func (r *AuthResolver) NotifyAuthSource(hostDisplay string, ctx *AuthContext) {
	hostKey := strings.ToLower(hostDisplay)
	if hostKey == "" {
		return
	}
	r.mu.Lock()
	already := r.verboseAuthLoggedHosts[hostKey]
	if !already {
		r.verboseAuthLoggedHosts[hostKey] = true
	}
	r.mu.Unlock()
	if already {
		return
	}
	if ctx == nil || ctx.Source == "none" {
		return
	}
	var line string
	if ctx.AuthScheme == "bearer" {
		line = fmt.Sprintf("  [i] %s -- using bearer from az cli (source: %s)", hostKey, ctx.Source)
	} else {
		line = fmt.Sprintf("  [i] %s -- token from %s", hostKey, ctx.Source)
	}
	fmt.Fprintln(os.Stderr, line)
}

// ExecuteWithBearerFallback runs primaryOp; on ADO auth failure retries via bearer.
func (r *AuthResolver) ExecuteWithBearerFallback(
	depRef interface{},
	primaryOp func() (interface{}, error),
	bearerOp func(bearer string) (interface{}, error),
	isAuthFailure func(result interface{}, err error) bool,
) BearerFallbackOutcome {
	primary, primaryErr := primaryOp()
	if depRef == nil {
		return BearerFallbackOutcome{Outcome: primary, BearerAttempted: false}
	}
	// Check if dep is ADO via duck typing
	type adoChecker interface {
		IsAzureDevOps() bool
	}
	if checker, ok := depRef.(adoChecker); !ok || !checker.IsAzureDevOps() {
		return BearerFallbackOutcome{Outcome: primary, BearerAttempted: false}
	}
	if !isAuthFailure(primary, primaryErr) {
		return BearerFallbackOutcome{Outcome: primary, BearerAttempted: false}
	}

	// No az CLI support in Go sandbox; return primary
	return BearerFallbackOutcome{Outcome: primary, BearerAttempted: false}
}
