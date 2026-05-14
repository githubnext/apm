// Package depreference provides the DependencyReference model -- the core
// dependency representation and parsing layer for the APM CLI.
//
// Migrated from: src/apm_cli/models/dependency/reference.py
package depreference

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/githubnext/apm/internal/utils/githubhost"
	"github.com/githubnext/apm/internal/utils/pathsecurity"
)

// defaultSchemePorts maps URI schemes to their default ports so that
// redundant explicit ports (https://host:443/...) can be stripped.
var defaultSchemePorts = map[string]int{
	"https": 443,
	"http":  80,
	"ssh":   22,
}

// VirtualPackageType classifies a virtual (sub-repo) package.
type VirtualPackageType int

const (
	VirtualPackageFile         VirtualPackageType = iota // Individual file (*.prompt.md etc.)
	VirtualPackageSubdirectory                            // Subdirectory package
)

// virtualFileExtensions lists the file extensions recognised as virtual FILE packages.
var virtualFileExtensions = []string{
	".prompt.md",
	".instructions.md",
	".chatmode.md",
	".agent.md",
}

// removedCollectionExtensions lists legacy collection-manifest extensions that
// are rejected at parse time with a migration message.
var removedCollectionExtensions = []string{
	".collection.yml",
	".collection.yaml",
}

// gitlabVirtualRootSegments is the set of first-path segments that, on
// GitLab, often start an in-repo virtual layout.
var gitlabVirtualRootSegments = map[string]bool{
	"prompts":      true,
	"instructions": true,
	"collections":  true,
}

// scpLikeRE matches SCP-style SSH URLs: <user>@<host>:<path>
// Mirrors the Python SCP_LIKE_RE used in cache/url_normalize.
var scpLikeRE = regexp.MustCompile(
	`^(?P<user>[^@]+)@(?P<host>[^:]+):(?P<path>.+)$`,
)

// aliasRE validates alias strings.
var aliasRE = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// adoRepoRE validates org/project/repo paths for Azure DevOps.
var adoRepoRE = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._\- ]+/[a-zA-Z0-9._\- ]+$`)

// DependencyReference is the central model for an APM dependency.
//
// Fields mirror the Python DependencyReference dataclass exactly.
type DependencyReference struct {
	RepoURL    string // e.g. "owner/repo" or "org/project/repo" for ADO
	Host       string // Optional host; empty means default (github.com)
	Port       int    // Non-standard SSH/HTTPS port; 0 means default
	// ExplicitScheme is the user-stated transport: "ssh", "https", "http",
	// or "" for shorthand notation.
	ExplicitScheme string
	Reference      string // e.g. "main", "v1.0.0", "abc123"
	Alias          string // Optional alias for the dependency
	VirtualPath    string // Path for virtual packages
	IsVirtual      bool   // True if this is a virtual package

	// Azure DevOps specific fields
	ADOOrganization string
	ADOProject      string
	ADORepo         string

	// Local path dependency
	IsLocal   bool
	LocalPath string // Original local path string

	// Monorepo parent inheritance
	IsParentRepoInheritance bool

	ArtifactoryPrefix string // e.g. "artifactory/github"

	// HTTP insecure dependency
	IsInsecure    bool
	AllowInsecure bool

	// SKILL_BUNDLE subset selection
	SkillSubset []string // sorted skill names, nil = all
}

// VirtualType returns the type of virtual package, or -1 if not virtual.
func (d *DependencyReference) VirtualType() VirtualPackageType {
	if !d.IsVirtual || d.VirtualPath == "" {
		return -1
	}
	for _, ext := range virtualFileExtensions {
		if strings.HasSuffix(d.VirtualPath, ext) {
			return VirtualPackageFile
		}
	}
	return VirtualPackageSubdirectory
}

// IsVirtualFile returns true when this is a virtual file package.
func (d *DependencyReference) IsVirtualFile() bool {
	return d.VirtualType() == VirtualPackageFile
}

// IsVirtualSubdirectory returns true when this is a virtual subdirectory package.
func (d *DependencyReference) IsVirtualSubdirectory() bool {
	return d.VirtualType() == VirtualPackageSubdirectory
}

// IsArtifactory returns true when this reference points to a JFrog Artifactory VCS repo.
func (d *DependencyReference) IsArtifactory() bool {
	return d.ArtifactoryPrefix != ""
}

// IsAzureDevOps returns true when this reference points to Azure DevOps.
func (d *DependencyReference) IsAzureDevOps() bool {
	return d.Host != "" && githubhost.IsAzureDevOpsHostname(d.Host)
}

// GetVirtualPackageName generates a package name for a virtual package.
//
//	owner/repo/prompts/code-review.prompt.md -> repo-code-review
//	owner/repo/collections/project-planning  -> repo-project-planning
func (d *DependencyReference) GetVirtualPackageName() string {
	if !d.IsVirtual || d.VirtualPath == "" {
		parts := strings.Split(d.RepoURL, "/")
		return parts[len(parts)-1]
	}
	repoParts := strings.Split(d.RepoURL, "/")
	repoName := "package"
	if len(repoParts) > 0 {
		repoName = repoParts[len(repoParts)-1]
	}
	pathParts := strings.Split(d.VirtualPath, "/")
	last := pathParts[len(pathParts)-1]
	for _, ext := range virtualFileExtensions {
		if strings.HasSuffix(last, ext) {
			last = last[:len(last)-len(ext)]
			break
		}
	}
	return repoName + "-" + last
}

// IsLocalPath returns true when dep_str looks like a local filesystem path.
func IsLocalPath(depStr string) bool {
	s := strings.TrimSpace(depStr)
	if strings.HasPrefix(s, "//") {
		return false
	}
	for _, pfx := range []string{"./", "../", "/", "~/", `~\`, `.\`, `..\`} {
		if strings.HasPrefix(s, pfx) {
			return true
		}
	}
	// Windows absolute path: drive letter + colon + separator
	if runtime.GOOS == "windows" || (len(s) >= 3 &&
		((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z')) &&
		s[1] == ':' && (s[2] == '\\' || s[2] == '/')) {
		return len(s) >= 3
	}
	return false
}

// GetUniqueKey returns a key for deduplication.
func (d *DependencyReference) GetUniqueKey() string {
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	if d.IsVirtual && d.VirtualPath != "" {
		return d.RepoURL + "/" + d.VirtualPath
	}
	return d.RepoURL
}

// effectiveHost returns d.Host or the default host (github.com).
func (d *DependencyReference) effectiveHost() string {
	if d.Host != "" {
		return d.Host
	}
	return githubhost.DefaultHost()
}

// hostLabel returns host:port or host.
func (d *DependencyReference) hostLabel() string {
	h := d.effectiveHost()
	if d.Port != 0 {
		return fmt.Sprintf("%s:%d", h, d.Port)
	}
	return h
}

// ToCanonical returns the canonical scheme-free identity string.
func (d *DependencyReference) ToCanonical() string {
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	host := d.effectiveHost()
	isDefault := strings.EqualFold(host, githubhost.DefaultHost())
	hl := d.hostLabel()

	var result string
	switch {
	case isDefault && d.Port == 0 && d.ArtifactoryPrefix == "":
		result = d.RepoURL
	case d.ArtifactoryPrefix != "":
		result = hl + "/" + d.ArtifactoryPrefix + "/" + d.RepoURL
	default:
		result = hl + "/" + d.RepoURL
	}
	if d.IsVirtual && d.VirtualPath != "" {
		result = result + "/" + d.VirtualPath
	}
	if d.Reference != "" {
		result = result + "#" + d.Reference
	}
	return result
}

// GetIdentity returns the identity (canonical without ref/alias).
func (d *DependencyReference) GetIdentity() string {
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	host := d.effectiveHost()
	isDefault := strings.EqualFold(host, githubhost.DefaultHost())
	hl := d.hostLabel()

	var result string
	switch {
	case isDefault && d.Port == 0 && d.ArtifactoryPrefix == "":
		result = d.RepoURL
	case d.ArtifactoryPrefix != "":
		result = hl + "/" + d.ArtifactoryPrefix + "/" + d.RepoURL
	default:
		result = hl + "/" + d.RepoURL
	}
	if d.IsVirtual && d.VirtualPath != "" {
		result = result + "/" + d.VirtualPath
	}
	return result
}

// GetCanonicalDependencyString is host-blind (filesystem-layout) canonical string.
func (d *DependencyReference) GetCanonicalDependencyString() string {
	return d.GetUniqueKey()
}

// GetInstallPath returns the canonical filesystem path under apm_modules_dir.
func (d *DependencyReference) GetInstallPath(apmModulesDir string) (string, error) {
	if d.IsLocal && d.LocalPath != "" {
		pkgDirName := filepath.Base(d.LocalPath)
		if pkgDirName == "" || pkgDirName == "." || pkgDirName == ".." {
			return "", fmt.Errorf("local path %q does not resolve to a named directory", d.LocalPath)
		}
		if err := pathsecurity.ValidatePathSegments(pkgDirName, "local package path", true, false); err != nil {
			return "", err
		}
		result := filepath.Join(apmModulesDir, "_local", pkgDirName)
		return pathsecurity.EnsurePathWithin(result, apmModulesDir)
	}

	repoParts := strings.Split(d.RepoURL, "/")
	if err := pathsecurity.ValidatePathSegments(d.RepoURL, "repo_url", false, false); err != nil {
		return "", err
	}
	if d.VirtualPath != "" {
		if err := pathsecurity.ValidatePathSegments(d.VirtualPath, "virtual_path", false, false); err != nil {
			return "", err
		}
	}

	var result string
	if d.IsVirtual {
		if d.IsVirtualSubdirectory() {
			if d.IsAzureDevOps() && len(repoParts) >= 3 {
				result = filepath.Join(apmModulesDir, repoParts[0], repoParts[1], repoParts[2], d.VirtualPath)
			} else if len(repoParts) >= 2 {
				parts := append(repoParts, strings.Split(d.VirtualPath, "/")...)
				result = filepath.Join(append([]string{apmModulesDir}, parts...)...)
			}
		} else {
			pkgName := d.GetVirtualPackageName()
			if d.IsAzureDevOps() && len(repoParts) >= 3 {
				result = filepath.Join(apmModulesDir, repoParts[0], repoParts[1], pkgName)
			} else if len(repoParts) >= 2 {
				result = filepath.Join(apmModulesDir, repoParts[0], pkgName)
			}
		}
	} else if d.IsAzureDevOps() && len(repoParts) >= 3 {
		result = filepath.Join(apmModulesDir, repoParts[0], repoParts[1], repoParts[2])
	} else if len(repoParts) >= 2 {
		result = filepath.Join(append([]string{apmModulesDir}, repoParts...)...)
	}

	if result == "" {
		result = filepath.Join(append([]string{apmModulesDir}, repoParts...)...)
	}

	return pathsecurity.EnsurePathWithin(result, apmModulesDir)
}

// ToGitHubURL converts to a full repository HTTPS URL.
func (d *DependencyReference) ToGitHubURL() string {
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	host := d.effectiveHost()
	netloc := host
	if d.Port != 0 {
		netloc = fmt.Sprintf("%s:%d", host, d.Port)
	}
	scheme := "https"
	if d.IsInsecure {
		scheme = "http"
	}
	if d.IsAzureDevOps() {
		proj := url.PathEscape(d.ADOProject)
		repo := url.PathEscape(d.ADORepo)
		return fmt.Sprintf("https://%s/%s/%s/_git/%s", netloc, d.ADOOrganization, proj, repo)
	}
	if d.ArtifactoryPrefix != "" {
		return fmt.Sprintf("%s://%s/%s/%s", scheme, netloc, d.ArtifactoryPrefix, d.RepoURL)
	}
	return fmt.Sprintf("%s://%s/%s", scheme, netloc, d.RepoURL)
}

// ToCloneURL is the same as ToGitHubURL for most purposes.
func (d *DependencyReference) ToCloneURL() string {
	return d.ToGitHubURL()
}

// GetDisplayName returns the alias, local path, virtual name, or repo URL.
func (d *DependencyReference) GetDisplayName() string {
	if d.Alias != "" {
		return d.Alias
	}
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	if d.IsVirtual {
		return d.GetVirtualPackageName()
	}
	return d.RepoURL
}

// String returns a human-readable representation.
func (d *DependencyReference) String() string {
	if d.IsLocal && d.LocalPath != "" {
		return d.LocalPath
	}
	var result string
	if d.Host != "" {
		hl := d.hostLabel()
		if d.ArtifactoryPrefix != "" {
			result = hl + "/" + d.ArtifactoryPrefix + "/" + d.RepoURL
		} else {
			result = hl + "/" + d.RepoURL
		}
	} else {
		result = d.RepoURL
	}
	if d.VirtualPath != "" {
		result += "/" + d.VirtualPath
	}
	if d.Reference != "" {
		result += "#" + d.Reference
	}
	if d.Alias != "" {
		result += "@" + d.Alias
	}
	return result
}

// ----- Parsing helpers -----

// parseSCPURL parses an SCP-shorthand SSH URL (user@host:path).
// Returns (host, port, repoURL, reference, alias, true) or ("","",…, false).
func parseSCPURL(depStr string) (host string, port int, repoURL, reference, alias string, ok bool) {
	m := scpLikeRE.FindStringSubmatch(depStr)
	if m == nil {
		return
	}
	sshRepo := m[3]
	if strings.Contains(sshRepo, "@") {
		idx := strings.LastIndex(sshRepo, "@")
		alias = strings.TrimSpace(sshRepo[idx+1:])
		sshRepo = sshRepo[:idx]
	}
	if strings.Contains(sshRepo, "#") {
		idx := strings.LastIndex(sshRepo, "#")
		reference = strings.TrimSpace(sshRepo[idx+1:])
		sshRepo = sshRepo[:idx]
	}
	if strings.HasSuffix(sshRepo, ".git") {
		sshRepo = sshRepo[:len(sshRepo)-4]
	}
	repoURL = strings.TrimSpace(sshRepo)
	if err := pathsecurity.ValidatePathSegments(repoURL, "SSH repository path", true, false); err != nil {
		return
	}
	host = m[2]
	ok = true
	return
}

// parseSSHProtocolURL parses ssh:// URLs.
func parseSSHProtocolURL(rawURL string) (host string, port int, repoURL, reference, alias string, ok bool) {
	if !strings.HasPrefix(rawURL, "ssh://") {
		return
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	host = u.Hostname()
	if p, err2 := parsePortInt(u.Port()); err2 == nil && p != 0 {
		port = p
		if port == defaultSchemePorts["ssh"] {
			port = 0
		}
	}
	path := strings.TrimPrefix(u.Path, "/")
	fragment := u.Fragment
	if fragment != "" {
		if strings.Contains(fragment, "@") {
			idx := strings.LastIndex(fragment, "@")
			reference = strings.TrimSpace(fragment[:idx])
			alias = strings.TrimSpace(fragment[idx+1:])
		} else {
			reference = strings.TrimSpace(fragment)
		}
	}
	if alias == "" && strings.Contains(path, "@") {
		idx := strings.LastIndex(path, "@")
		alias = strings.TrimSpace(path[idx+1:])
		path = path[:idx]
	}
	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}
	repoURL = strings.TrimSpace(path)
	if err2 := pathsecurity.ValidatePathSegments(repoURL, "SSH repository path", true, false); err2 != nil {
		return
	}
	ok = true
	return
}

func parsePortInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	var p int
	_, err := fmt.Sscanf(s, "%d", &p)
	return p, err
}

// hasVirtualExt returns true if any segment ends in a virtual file extension.
func hasVirtualExt(segments []string) bool {
	for _, seg := range segments {
		for _, ext := range virtualFileExtensions {
			if strings.HasSuffix(seg, ext) {
				return true
			}
		}
	}
	return false
}

// gitlabSegmentCount computes how many path segments belong to the GitLab
// project path vs the virtual package suffix.
func gitlabSegmentCount(segs []string, hasVirtExt, hasCollection bool) int {
	n := len(segs)
	if n < 2 {
		return n
	}
	if hasCollection {
		for i, s := range segs {
			if s == "collections" && i >= 2 {
				return i
			}
		}
		return n
	}
	if hasVirtExt {
		for i, seg := range segs {
			if i >= 2 && gitlabVirtualRootSegments[seg] {
				return i
			}
		}
		if n == 3 {
			return 2
		}
		if n == 4 {
			return 3
		}
		if n >= 5 {
			return 3
		}
		return 2
	}
	return n
}

// detectVirtualPackage scans a dependency string for virtual package indicators.
// Returns (isVirtual, virtualPath, validatedHost, error).
func detectVirtualPackage(depStr string) (bool, string, string, error) {
	temp := depStr
	if idx := strings.LastIndex(temp, "#"); idx >= 0 {
		temp = temp[:idx]
	}

	lower := strings.ToLower(temp)
	for _, pfx := range []string{"git@", "https://", "http://", "ssh://"} {
		if strings.HasPrefix(lower, pfx) {
			return false, "", "", nil
		}
	}

	check := temp
	var validatedHost string

	if strings.Contains(check, "/") {
		firstSeg := strings.SplitN(check, "/", 2)[0]
		if strings.Contains(firstSeg, ".") {
			testURL := "https://" + check
			u, err := url.Parse(testURL)
			if err == nil && u.Hostname() != "" && githubhost.IsSupportedGitHost(u.Hostname()) {
				validatedHost = u.Hostname()
				check = strings.SplitN(check, "/", 2)[1]
			} else if err == nil {
				return false, "", "", fmt.Errorf("invalid Git host: %s", firstSeg)
			}
		} else if strings.HasPrefix(check, "gh/") {
			check = check[3:]
		}
	}

	pathSegments := filterEmpty(strings.Split(check, "/"))

	isADO := validatedHost != "" && githubhost.IsAzureDevOpsHostname(validatedHost)
	isGenericHost := validatedHost != "" && !githubhost.IsGitHubHostname(validatedHost) && !githubhost.IsAzureDevOpsHostname(validatedHost)
	isGitLabHost := validatedHost != "" && githubhost.IsGitLabHostname(validatedHost)

	if isADO {
		for i, s := range pathSegments {
			if s == "_git" {
				pathSegments = append(pathSegments[:i], pathSegments[i+1:]...)
				break
			}
		}
	}

	isArtifactory := isGenericHost && githubhost.IsArtifactoryPath(pathSegments)

	var minBaseSegments int
	switch {
	case isADO:
		if validatedHost != "" && githubhost.IsVisualStudioLegacyHostname(validatedHost) {
			minBaseSegments = 2
		} else {
			minBaseSegments = 3
		}
	case isArtifactory:
		minBaseSegments = 4
	case isGenericHost:
		hv := hasVirtualExt(pathSegments)
		hc := contains(pathSegments, "collections")
		if isGitLabHost {
			minBaseSegments = gitlabSegmentCount(pathSegments, hv, hc)
		} else if hv || hc {
			minBaseSegments = 2
		} else {
			minBaseSegments = len(pathSegments)
		}
	default:
		minBaseSegments = 2
	}

	if len(pathSegments) >= minBaseSegments+1 {
		vPath := strings.Join(pathSegments[minBaseSegments:], "/")
		if err := pathsecurity.ValidatePathSegments(vPath, "virtual path", false, false); err != nil {
			return false, "", validatedHost, err
		}
		for _, ext := range removedCollectionExtensions {
			if strings.HasSuffix(vPath, ext) {
				return false, "", validatedHost, fmt.Errorf(
					".collection.yml is no longer supported. Convert %q to an apm.yml with a 'dependencies' section", vPath)
			}
		}
		for _, ext := range virtualFileExtensions {
			if strings.HasSuffix(vPath, ext) {
				return true, vPath, validatedHost, nil
			}
		}
		last := vPath
		if idx := strings.LastIndex(vPath, "/"); idx >= 0 {
			last = vPath[idx+1:]
		}
		if strings.Contains(last, ".") {
			return false, "", validatedHost, fmt.Errorf(
				"invalid virtual package path %q: individual files must end with a recognized extension", vPath)
		}
		return true, vPath, validatedHost, nil
	}

	return false, "", validatedHost, nil
}

func filterEmpty(ss []string) []string {
	out := ss[:0]
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// validateURLRepoPath validates and normalises the repo path from a parsed URL.
// Returns (repoURL, virtualPath, error).
func validateURLRepoPath(u *url.URL) (string, string, error) {
	hostname := u.Hostname()
	if !githubhost.IsSupportedGitHost(hostname) {
		return "", "", fmt.Errorf("invalid Git host: %s", hostname)
	}

	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		return "", "", fmt.Errorf("repository path cannot be empty")
	}
	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}

	pathParts := make([]string, 0)
	for _, p := range strings.Split(path, "/") {
		pathParts = append(pathParts, urlUnescape(p))
	}
	// Remove _git segment (Azure DevOps)
	for i, p := range pathParts {
		if p == "_git" {
			pathParts = append(pathParts[:i], pathParts[i+1:]...)
			break
		}
	}

	isADO := githubhost.IsAzureDevOpsHostname(hostname)
	var urlVirtualPath string

	if isADO {
		isVSLegacy := githubhost.IsVisualStudioLegacyHostname(hostname)
		minParts := 3
		if isVSLegacy {
			minParts = 2
		}
		if len(pathParts) < minParts {
			return "", "", fmt.Errorf("invalid Azure DevOps repository path: expected 'org/project/repo', got %q", path)
		}
		if len(pathParts) > minParts {
			adoVirtual := strings.Join(pathParts[minParts:], "/")
			if err := pathsecurity.ValidatePathSegments(adoVirtual, "virtual path", false, false); err != nil {
				return "", "", err
			}
			for _, ext := range removedCollectionExtensions {
				if strings.HasSuffix(adoVirtual, ext) {
					return "", "", fmt.Errorf(".collection.yml is no longer supported for %q", adoVirtual)
				}
			}
			isFile := false
			for _, ext := range virtualFileExtensions {
				if strings.HasSuffix(adoVirtual, ext) {
					isFile = true
					break
				}
			}
			if !isFile {
				last := adoVirtual
				if idx := strings.LastIndex(adoVirtual, "/"); idx >= 0 {
					last = adoVirtual[idx+1:]
				}
				if strings.Contains(last, ".") {
					return "", "", fmt.Errorf("invalid virtual package path %q", adoVirtual)
				}
			}
			urlVirtualPath = adoVirtual
			pathParts = pathParts[:minParts]
		}
		if isVSLegacy {
			vsOrg := strings.SplitN(hostname, ".", 2)[0]
			pathParts = append([]string{vsOrg}, pathParts...)
		}
	} else {
		if len(pathParts) < 2 {
			return "", "", fmt.Errorf("invalid repository path: expected at least 'user/repo', got %q", path)
		}
		for _, pp := range pathParts {
			for _, ext := range virtualFileExtensions {
				if strings.HasSuffix(pp, ext) {
					return "", "", fmt.Errorf("invalid repository path %q: contains a virtual file extension; use dict format with 'path:' for virtual packages", path)
				}
			}
		}
	}

	isADOPath := githubhost.IsAzureDevOpsHostname(hostname)
	allowedPattern := `^[a-zA-Z0-9._-]+$`
	if isADOPath {
		allowedPattern = `^[a-zA-Z0-9._\- ]+$`
	}
	allowedRE := regexp.MustCompile(allowedPattern)

	if err := pathsecurity.ValidatePathSegments(strings.Join(pathParts, "/"), "repository URL path", true, false); err != nil {
		return "", "", err
	}
	for _, part := range pathParts {
		if !allowedRE.MatchString(part) {
			return "", "", fmt.Errorf("invalid repository path component: %s", part)
		}
	}

	return strings.Join(pathParts, "/"), urlVirtualPath, nil
}

func urlUnescape(s string) string {
	out, err := url.PathUnescape(s)
	if err != nil {
		return s
	}
	return out
}

// resolveVirtualShorthandRepo strips the virtual suffix from a shorthand repo_url.
// Returns (host, repoURL).
func resolveVirtualShorthandRepo(repoURL, validatedHost, virtualPath string) (string, string) {
	parts := filterEmpty(strings.Split(repoURL, "/"))
	// Remove _git
	for i, p := range parts {
		if p == "_git" {
			parts = append(parts[:i], parts[i+1:]...)
			break
		}
	}

	host := ""
	if len(parts) >= 3 && githubhost.IsSupportedGitHost(parts[0]) {
		host = parts[0]
		if githubhost.IsAzureDevOpsHostname(parts[0]) {
			if githubhost.IsVisualStudioLegacyHostname(parts[0]) {
				if len(parts) >= 4 {
					repoURL = strings.Join(parts[1:3], "/")
				}
			} else {
				if len(parts) >= 5 {
					repoURL = strings.Join(parts[1:4], "/")
				}
			}
		} else if githubhost.IsArtifactoryPath(parts[1:]) {
			prefix, owner, repo := githubhost.ParseArtifactoryPath(parts[1:])
			if owner != "" && repo != "" {
				_ = prefix
				repoURL = owner + "/" + repo
			}
		} else if githubhost.IsGitLabHostname(parts[0]) && virtualPath != "" {
			vParts := filterEmpty(strings.Split(virtualPath, "/"))
			tail := len(vParts)
			if tail > 0 && len(parts) > 1+tail {
				repoURL = strings.Join(parts[1:len(parts)-tail], "/")
			} else {
				repoURL = strings.Join(parts[1:], "/")
			}
		} else {
			repoURL = strings.Join(parts[1:3], "/")
		}
	} else if len(parts) >= 2 {
		if host == "" {
			host = githubhost.DefaultHost()
		}
		if validatedHost != "" && githubhost.IsAzureDevOpsHostname(validatedHost) {
			if len(parts) >= 4 {
				repoURL = strings.Join(parts[:3], "/")
			}
		} else {
			repoURL = strings.Join(parts[:2], "/")
		}
	}
	return host, repoURL
}

// resolveShorthandToParsedURL converts a shorthand to a *url.URL.
func resolveShorthandToParsedURL(repoURL, host string) (*url.URL, string, error) {
	parts := filterEmpty(strings.Split(repoURL, "/"))
	for i, p := range parts {
		if p == "_git" {
			parts = append(parts[:i], parts[i+1:]...)
			break
		}
	}

	var userRepo string
	if len(parts) >= 3 && githubhost.IsSupportedGitHost(parts[0]) {
		host = parts[0]
		if githubhost.IsVisualStudioLegacyHostname(host) && len(parts) >= 3 {
			userRepo = strings.Join(parts[1:3], "/")
		} else if githubhost.IsAzureDevOpsHostname(host) && len(parts) >= 4 {
			userRepo = strings.Join(parts[1:4], "/")
		} else if !githubhost.IsGitHubHostname(host) && !githubhost.IsAzureDevOpsHostname(host) {
			if githubhost.IsArtifactoryPath(parts[1:]) {
				_, owner, repo := githubhost.ParseArtifactoryPath(parts[1:])
				if owner != "" && repo != "" {
					userRepo = owner + "/" + repo
				} else {
					userRepo = strings.Join(parts[1:], "/")
				}
			} else {
				userRepo = strings.Join(parts[1:], "/")
			}
		} else {
			userRepo = strings.Join(parts[1:], "/")
		}
	} else if len(parts) >= 2 && !strings.Contains(parts[0], ".") {
		if host == "" {
			host = githubhost.DefaultHost()
		}
		if githubhost.IsAzureDevOpsHostname(host) && len(parts) >= 3 {
			userRepo = strings.Join(parts[:3], "/")
		} else if host != "" && !githubhost.IsGitHubHostname(host) && !githubhost.IsAzureDevOpsHostname(host) {
			userRepo = strings.Join(parts, "/")
		} else {
			userRepo = strings.Join(parts[:2], "/")
		}
	} else {
		return nil, "", fmt.Errorf("use 'user/repo' or 'github.com/user/repo' format")
	}

	if userRepo == "" || !strings.Contains(userRepo, "/") {
		return nil, "", fmt.Errorf("invalid repository format: %s", repoURL)
	}

	uParts := strings.Split(userRepo, "/")
	isADOHost := host != "" && githubhost.IsAzureDevOpsHostname(host)

	if isADOHost {
		minADOParts := 3
		if githubhost.IsVisualStudioLegacyHostname(host) {
			minADOParts = 2
		}
		if len(uParts) < minADOParts {
			return nil, "", fmt.Errorf("invalid Azure DevOps repository format: %s", repoURL)
		}
	} else if len(uParts) < 2 {
		return nil, "", fmt.Errorf("invalid repository format: %s", repoURL)
	}

	if err := pathsecurity.ValidatePathSegments(strings.Join(uParts, "/"), "repository path", false, false); err != nil {
		return nil, "", err
	}

	allowedPattern := `^[a-zA-Z0-9._-]+$`
	if isADOHost {
		allowedPattern = `^[a-zA-Z0-9._\- ]+$`
	}
	allowedRE := regexp.MustCompile(allowedPattern)
	for _, part := range uParts {
		stripped := strings.TrimSuffix(part, ".git")
		if !allowedRE.MatchString(stripped) {
			return nil, "", fmt.Errorf("invalid repository path component: %s", part)
		}
	}

	escapedParts := make([]string, len(uParts))
	for i, p := range uParts {
		escapedParts[i] = url.PathEscape(p)
	}
	rawURL := fmt.Sprintf("https://%s/%s", host, strings.Join(escapedParts, "/"))
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to build URL for %s: %w", repoURL, err)
	}
	return parsed, host, nil
}

// parseStandardURL handles non-SSH dependency strings.
func parseStandardURL(depStr string, isVirtual bool, virtualPath, validatedHost string) (
	host string, port int, repoURL, reference, alias string,
	effectiveIsVirtual bool, effectiveVirtualPath string, err error,
) {
	effectiveIsVirtual = isVirtual
	effectiveVirtualPath = virtualPath

	repoPart := depStr
	if idx := strings.LastIndex(depStr, "#"); idx >= 0 {
		repoPart = depStr[:idx]
		reference = strings.TrimSpace(depStr[idx+1:])
	}
	repoURL = strings.TrimSpace(repoPart)
	lower := strings.ToLower(repoURL)

	if isVirtual && !strings.HasPrefix(lower, "https://") && !strings.HasPrefix(lower, "http://") {
		host, repoURL = resolveVirtualShorthandRepo(repoURL, validatedHost, virtualPath)
	}

	lower = strings.ToLower(repoURL)
	var parsedURL *url.URL
	if strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://") {
		parsedURL, err = url.Parse(repoURL)
		if err != nil {
			return
		}
		host = parsedURL.Hostname()
		if p, e := parsePortInt(parsedURL.Port()); e == nil {
			port = p
		}
		scheme := strings.ToLower(parsedURL.Scheme)
		if port == defaultSchemePorts[scheme] {
			port = 0
		}
	} else {
		parsedURL, host, err = resolveShorthandToParsedURL(repoURL, host)
		if err != nil {
			return
		}
	}

	var urlVirtualPath string
	repoURL, urlVirtualPath, err = validateURLRepoPath(parsedURL)
	if err != nil {
		return
	}
	if urlVirtualPath != "" {
		effectiveIsVirtual = true
		effectiveVirtualPath = urlVirtualPath
	}
	if host == "" {
		host = githubhost.DefaultHost()
	}
	return
}

// validateFinalRepoFields checks the final repo_url and extracts ADO fields.
func validateFinalRepoFields(host, repoURL string) (adoOrg, adoProject, adoRepo string, err error) {
	isADO := host != "" && githubhost.IsAzureDevOpsHostname(host)
	if isADO {
		if !adoRepoRE.MatchString(repoURL) {
			err = fmt.Errorf("invalid Azure DevOps repository format: %s; expected 'org/project/repo'", repoURL)
			return
		}
		parts := strings.SplitN(repoURL, "/", 3)
		if err2 := pathsecurity.ValidatePathSegments(repoURL, "Azure DevOps repository path", false, false); err2 != nil {
			err = err2
			return
		}
		adoOrg, adoProject, adoRepo = parts[0], parts[1], parts[2]
		return
	}

	segments := strings.Split(repoURL, "/")
	if len(segments) < 2 {
		err = fmt.Errorf("invalid repository format: %s; expected 'user/repo'", repoURL)
		return
	}
	validRE := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for _, s := range segments {
		if !validRE.MatchString(s) {
			err = fmt.Errorf("invalid repository format: %s; contains invalid characters", repoURL)
			return
		}
		for _, ext := range virtualFileExtensions {
			if strings.HasSuffix(s, ext) {
				err = fmt.Errorf("invalid repository format: %q contains a virtual file extension", repoURL)
				return
			}
		}
	}
	if e := pathsecurity.ValidatePathSegments(repoURL, "repository path", false, false); e != nil {
		err = e
	}
	return
}

// extractArtifactoryPrefix extracts the Artifactory VCS prefix from the original dep string.
func extractArtifactoryPrefix(depStr, host string) string {
	s := depStr
	if idx := strings.Index(s, "#"); idx >= 0 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "@"); idx >= 0 {
		s = s[:idx]
	}
	if strings.Contains(s, "://") {
		s = strings.SplitN(s, "://", 2)[1]
	}
	s = strings.Replace(s, host+"/", "", 1)
	segs := filterEmpty(strings.Split(s, "/"))
	if githubhost.IsArtifactoryPath(segs) {
		prefix, _, _ := githubhost.ParseArtifactoryPath(segs)
		return prefix
	}
	return ""
}

// Parse parses a dependency string into a DependencyReference.
//
// Supports all forms: shorthand (user/repo), FQDN, HTTPS, SSH, SCP, local paths.
func Parse(depStr string) (*DependencyReference, error) {
	if strings.TrimSpace(depStr) == "" {
		return nil, fmt.Errorf("empty dependency string")
	}

	depStr, err := url.PathUnescape(depStr)
	if err != nil {
		depStr = depStr // keep original on error
	}

	for _, r := range depStr {
		if r < 32 && !unicode.IsSpace(r) {
			return nil, fmt.Errorf("dependency string contains invalid control characters")
		}
	}

	// Local path detection (must run before URL/host parsing)
	if IsLocalPath(depStr) {
		local := strings.TrimSpace(depStr)
		base := filepath.Base(local)
		if base == "" || base == "." || base == ".." {
			return nil, fmt.Errorf("local path %q does not resolve to a named directory", local)
		}
		return &DependencyReference{
			RepoURL:   "_local/" + base,
			IsLocal:   true,
			LocalPath: local,
		}, nil
	}

	if strings.HasPrefix(depStr, "//") {
		return nil, fmt.Errorf("protocol-relative URLs are not supported")
	}

	// Phase 1: detect virtual packages
	isVirtual, virtualPath, validatedHost, err := detectVirtualPackage(depStr)
	if err != nil {
		return nil, err
	}

	// Phase 2: SSH parsing
	var (
		host           string
		port           int
		repoURL        string
		reference      string
		alias          string
		explicitScheme string
	)

	if h, p, r, ref, al, ok := parseSSHProtocolURL(depStr); ok {
		host, port, repoURL, reference, alias = h, p, r, ref, al
		explicitScheme = "ssh"
	} else if h, p, r, ref, al, ok2 := parseSCPURL(depStr); ok2 {
		host, port, repoURL, reference, alias = h, p, r, ref, al
		explicitScheme = "ssh"
	} else {
		var effectiveIsVirtual bool
		var effectiveVirtualPath string
		host, port, repoURL, reference, alias, effectiveIsVirtual, effectiveVirtualPath, err =
			parseStandardURL(depStr, isVirtual, virtualPath, validatedHost)
		if err != nil {
			return nil, err
		}
		isVirtual = effectiveIsVirtual
		virtualPath = effectiveVirtualPath
		lower := strings.ToLower(strings.TrimSpace(depStr))
		if strings.HasPrefix(lower, "https://") {
			explicitScheme = "https"
		} else if strings.HasPrefix(lower, "http://") {
			explicitScheme = "http"
		}
	}

	// Phase 3: validate final fields
	adoOrg, adoProject, adoRepo, err := validateFinalRepoFields(host, repoURL)
	if err != nil {
		return nil, err
	}

	if alias != "" && !aliasRE.MatchString(alias) {
		return nil, fmt.Errorf("invalid alias: %s; aliases can only contain letters, numbers, dots, underscores, and hyphens", alias)
	}

	isADO := host != "" && githubhost.IsAzureDevOpsHostname(host)
	var artifactoryPrefix string
	if host != "" && !isADO {
		artifactoryPrefix = extractArtifactoryPrefix(depStr, host)
	}

	parsedScheme := ""
	if u, e := url.Parse(depStr); e == nil {
		parsedScheme = strings.ToLower(u.Scheme)
	}

	return &DependencyReference{
		RepoURL:                 repoURL,
		Host:                    host,
		Port:                    port,
		ExplicitScheme:          explicitScheme,
		Reference:               reference,
		Alias:                   alias,
		VirtualPath:             virtualPath,
		IsVirtual:               isVirtual,
		ADOOrganization:         adoOrg,
		ADOProject:              adoProject,
		ADORepo:                 adoRepo,
		ArtifactoryPrefix:       artifactoryPrefix,
		IsInsecure:              parsedScheme == "http",
		IsParentRepoInheritance: false,
	}, nil
}

// Canonicalize parses raw and returns its canonical form.
func Canonicalize(raw string) (string, error) {
	ref, err := Parse(raw)
	if err != nil {
		return "", err
	}
	return ref.ToCanonical(), nil
}

// ParseFromDict parses a dict-style dependency entry (as in apm.yml).
func ParseFromDict(entry map[string]interface{}) (*DependencyReference, error) {
	pathVal, hasPath := entry["path"]
	gitVal, hasGit := entry["git"]

	if hasPath && !hasGit {
		localStr, ok := pathVal.(string)
		if !ok || strings.TrimSpace(localStr) == "" {
			return nil, fmt.Errorf("'path' field must be a non-empty string")
		}
		localStr = strings.TrimSpace(localStr)
		if !IsLocalPath(localStr) {
			return nil, fmt.Errorf("object-style dependency must have a 'git' field, or 'path' must be a local filesystem path")
		}
		return Parse(localStr)
	}

	if !hasGit {
		return nil, fmt.Errorf("object-style dependency must have a 'git' or 'path' field")
	}

	gitURL, ok := gitVal.(string)
	if !ok || strings.TrimSpace(gitURL) == "" {
		return nil, fmt.Errorf("'git' field must be a non-empty string")
	}
	gitURL = strings.TrimSpace(gitURL)

	// Parent repo inheritance
	if gitURL == "parent" {
		pathRaw, _ := entry["path"].(string)
		if strings.TrimSpace(pathRaw) == "" {
			return nil, fmt.Errorf("object-style dependency with git: 'parent' requires a 'path' field")
		}
		normPath := normalizeParentRepoPath(pathRaw)
		if normPath == "" {
			return nil, fmt.Errorf("'path' field must be a non-empty string")
		}
		dep := &DependencyReference{
			RepoURL:                 "_parent",
			IsVirtual:               true,
			IsParentRepoInheritance: true,
			VirtualPath:             normPath,
		}
		if refRaw, ok2 := entry["ref"].(string); ok2 && strings.TrimSpace(refRaw) != "" {
			dep.Reference = strings.TrimSpace(refRaw)
		}
		if aliasRaw, ok2 := entry["alias"].(string); ok2 && strings.TrimSpace(aliasRaw) != "" {
			a := strings.TrimSpace(aliasRaw)
			if !aliasRE.MatchString(a) {
				return nil, fmt.Errorf("invalid alias: %s", a)
			}
			dep.Alias = a
		}
		return dep, nil
	}

	dep, err := Parse(gitURL)
	if err != nil {
		return nil, err
	}

	if allowInsecure, ok2 := entry["allow_insecure"].(bool); ok2 {
		dep.AllowInsecure = allowInsecure
	}

	if refRaw, ok2 := entry["ref"].(string); ok2 && strings.TrimSpace(refRaw) != "" {
		dep.Reference = strings.TrimSpace(refRaw)
	}

	if aliasRaw, ok2 := entry["alias"].(string); ok2 && strings.TrimSpace(aliasRaw) != "" {
		a := strings.TrimSpace(aliasRaw)
		if !aliasRE.MatchString(a) {
			return nil, fmt.Errorf("invalid alias: %s", a)
		}
		dep.Alias = a
	}

	if subPath, ok2 := entry["path"].(string); ok2 && strings.TrimSpace(subPath) != "" {
		sp := strings.TrimSpace(strings.ReplaceAll(subPath, `\`, "/"))
		sp = strings.Trim(sp, "/")
		if err2 := pathsecurity.ValidatePathSegments(sp, "path", false, false); err2 != nil {
			return nil, err2
		}
		dep.VirtualPath = sp
		dep.IsVirtual = true
	}

	if skillsRaw, ok2 := entry["skills"].([]interface{}); ok2 {
		if len(skillsRaw) == 0 {
			return nil, fmt.Errorf("skills: must contain at least one name")
		}
		seen := map[string]bool{}
		var validated []string
		for _, s := range skillsRaw {
			name, ok3 := s.(string)
			if !ok3 || strings.TrimSpace(name) == "" {
				return nil, fmt.Errorf("each entry in 'skills' must be a non-empty string")
			}
			name = strings.TrimSpace(name)
			if err2 := pathsecurity.ValidatePathSegments(name, "skills/<name>", false, false); err2 != nil {
				return nil, err2
			}
			if !seen[name] {
				seen[name] = true
				validated = append(validated, name)
			}
		}
		dep.SkillSubset = sortedStrings(validated)
	}

	return dep, nil
}

func normalizeParentRepoPath(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, `\`, "/")
	s = strings.Trim(s, "/")
	parts := filterEmpty(strings.Split(s, "/"))
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "/")
}

func sortedStrings(ss []string) []string {
	out := make([]string, len(ss))
	copy(out, ss)
	// simple insertion sort (skill lists are short)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// ToApmYMLEntry returns the value to store in apm.yml.
// Returns a string for simple deps, or a map for HTTP/skill-subset deps.
func (d *DependencyReference) ToApmYMLEntry() interface{} {
	if d.IsInsecure {
		host := d.effectiveHost()
		entry := map[string]interface{}{
			"git": "http://" + host + "/" + d.RepoURL,
		}
		if d.Reference != "" {
			entry["ref"] = d.Reference
		}
		if d.Alias != "" {
			entry["alias"] = d.Alias
		}
		entry["allow_insecure"] = d.AllowInsecure
		if len(d.SkillSubset) > 0 {
			entry["skills"] = sortedStrings(d.SkillSubset)
		}
		return entry
	}
	if len(d.SkillSubset) > 0 {
		entry := map[string]interface{}{
			"git": d.GetIdentity(),
		}
		if d.Reference != "" {
			entry["ref"] = d.Reference
		}
		if d.Alias != "" {
			entry["alias"] = d.Alias
		}
		entry["skills"] = sortedStrings(d.SkillSubset)
		return entry
	}
	return d.ToCanonical()
}

// VirtualSuffixIsInstallableShape returns true when virtualPath matches APM virtual package rules.
func VirtualSuffixIsInstallableShape(virtualPath string) bool {
	if strings.TrimSpace(virtualPath) == "" {
		return false
	}
	v := strings.Trim(strings.TrimSpace(virtualPath), "/")
	if err := pathsecurity.ValidatePathSegments(v, "virtual path", false, false); err != nil {
		return false
	}
	if strings.Contains(v, "/collections/") || strings.HasPrefix(v, "collections/") {
		return true
	}
	for _, ext := range virtualFileExtensions {
		if strings.HasSuffix(v, ext) {
			return true
		}
	}
	last := v
	if idx := strings.LastIndex(v, "/"); idx >= 0 {
		last = v[idx+1:]
	}
	return !strings.Contains(last, ".")
}
