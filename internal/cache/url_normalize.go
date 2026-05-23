package cache

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// scpLikeRE matches SCP-style SSH URLs: user@host:path
var scpLikeRE = regexp.MustCompile(
	`^(?P<user>[a-zA-Z0-9_][a-zA-Z0-9_.+-]*)@` +
		`(?P<host>[^:/]+)` +
		`:(?P<path>.+)$`,
)

// defaultPorts maps schemes to their default TCP ports.
var defaultPorts = map[string]int{
	"https": 443,
	"ssh":   22,
	"http":  80,
	"git":   9418,
}

// caseInsensitiveHosts are hosts where the URL path is treated case-insensitively.
var caseInsensitiveHosts = map[string]bool{
	"github.com":    true,
	"gitlab.com":    true,
	"bitbucket.org": true,
}

// NormalizeRepoURL normalises a Git repository URL for cache key derivation.
// The result is a canonical string suitable for hashing. It is NOT necessarily
// a valid URL -- it is a deterministic representation.
func NormalizeRepoURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)

	// Convert SCP-like (git@host:path) to ssh:// form
	if m := scpLikeRE.FindStringSubmatch(rawURL); m != nil {
		user := scpLikeRE.SubexpIndex("user")
		host := scpLikeRE.SubexpIndex("host")
		path := scpLikeRE.SubexpIndex("path")
		p := m[path]
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		rawURL = fmt.Sprintf("ssh://%s@%s%s", m[user], m[host], p)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Lowercase hostname
	hostname := strings.ToLower(parsed.Hostname())

	// Keep username, drop password
	username := ""
	if parsed.User != nil {
		username = parsed.User.Username()
	}

	// Strip default ports
	scheme := strings.ToLower(parsed.Scheme)
	if scheme == "" {
		scheme = "https"
	}
	portInt := 0
	if p, err2 := url.ParseRequestURI(rawURL); err2 == nil {
		if p.Port() != "" {
			fmt.Sscanf(p.Port(), "%d", &portInt)
		}
	}
	if def, ok := defaultPorts[scheme]; ok && portInt == def {
		portInt = 0
	}

	// Reconstruct authority
	authority := hostname
	if username != "" {
		authority = username + "@" + hostname
	}
	if portInt != 0 {
		authority = fmt.Sprintf("%s:%d", authority, portInt)
	}

	// Strip trailing .git from path
	path := parsed.Path
	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}

	// Lowercase path for known case-insensitive hosts
	if caseInsensitiveHosts[hostname] {
		path = strings.ToLower(path)
	}

	// Strip trailing slash
	path = strings.TrimRight(path, "/")

	return fmt.Sprintf("%s://%s%s", scheme, authority, path)
}

// CacheShardKey derives a filesystem-safe shard key from a repository URL.
// Returns the first 16 hex characters of the SHA-256 of the normalised URL.
func CacheShardKey(rawURL string) string {
	normalized := NormalizeRepoURL(rawURL)
	h := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("%x", h)[:16]
}
