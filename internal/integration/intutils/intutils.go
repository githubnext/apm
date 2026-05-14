// Package intutils provides shared utility functions for integration modules.
package intutils

import "strings"

// NormalizeRepoURL normalizes a repo URL to owner/repo format.
func NormalizeRepoURL(packageRepoURL string) string {
url := packageRepoURL
if !strings.Contains(url, "://") {
url = strings.TrimSuffix(url, ".git")
return strings.TrimRight(url, "/")
}
parts := strings.SplitN(url, "://", 2)
if len(parts) < 2 {
return url
}
rest := parts[1]
slashIdx := strings.Index(rest, "/")
if slashIdx < 0 {
return url
}
path := rest[slashIdx+1:]
path = strings.TrimRight(path, "/")
path = strings.TrimSuffix(path, ".git")
return path
}
