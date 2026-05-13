// Package urlnormalize provides URL normalization for cache key derivation.
package urlnormalize

import (
"crypto/sha256"
"fmt"
"regexp"
"strings"
)

var scpLikeRe = regexp.MustCompile(`^(?P<user>[a-zA-Z0-9_][a-zA-Z0-9_.+-]*)@(?P<host>[^:/]+):(?P<path>.+)$`)

var defaultPorts = map[string]string{
"https": "443",
"ssh":   "22",
"http":  "80",
"git":   "9418",
}

// NormalizeRepoURL normalizes a git repository URL for cache key derivation.
func NormalizeRepoURL(url string) string {
u := strings.TrimSpace(url)
// Strip trailing .git
u = strings.TrimSuffix(u, ".git")

// SCP -> SSH URL conversion
if m := scpLikeRe.FindStringSubmatch(u); m != nil {
user := m[scpLikeRe.SubexpIndex("user")]
host := m[scpLikeRe.SubexpIndex("host")]
path := m[scpLikeRe.SubexpIndex("path")]
u = fmt.Sprintf("ssh://%s@%s/%s", user, strings.ToLower(host), path)
}

// Parse scheme://[user@]host[:port]/path
scheme := ""
rest := u
if idx := strings.Index(u, "://"); idx >= 0 {
scheme = strings.ToLower(u[:idx])
rest = u[idx+3:]
}

// Separate userinfo@host:port from path
var userinfo, hostport, path string
if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
hostport = rest[:slashIdx]
path = rest[slashIdx:]
} else {
hostport = rest
}

// Split userinfo from host
if atIdx := strings.LastIndex(hostport, "@"); atIdx >= 0 {
userinfo = hostport[:atIdx]
hostport = hostport[atIdx+1:]
}

// Strip password from userinfo
if colonIdx := strings.Index(userinfo, ":"); colonIdx >= 0 {
userinfo = userinfo[:colonIdx]
}

// Lowercase host, strip default port
hostLower := strings.ToLower(hostport)
if colonIdx := strings.LastIndex(hostLower, ":"); colonIdx >= 0 {
host := hostLower[:colonIdx]
port := hostLower[colonIdx+1:]
if dp, ok := defaultPorts[scheme]; ok && port == dp {
hostLower = host
}
}

// Lowercase github/gitlab/bitbucket paths
pathNorm := path
if hostLower == "github.com" || hostLower == "gitlab.com" || hostLower == "bitbucket.org" {
pathNorm = strings.ToLower(path)
}

// Reassemble
result := ""
if scheme != "" {
result = scheme + "://"
}
if userinfo != "" {
result += userinfo + "@"
}
result += hostLower + pathNorm
return result
}

// CacheKey returns the first 16 hex chars of SHA256 of the normalized URL.
func CacheKey(url string) string {
normalized := NormalizeRepoURL(url)
sum := sha256.Sum256([]byte(normalized))
return fmt.Sprintf("%x", sum)[:16]
}
