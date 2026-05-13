// Package versionchecker provides version checking and update notification utilities.
package versionchecker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var versionRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(a\d+|b\d+|rc\d+)?$`)

// VersionComponents holds parsed version parts.
type VersionComponents struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

// GetLatestVersionFromGitHub fetches the latest release version from GitHub API.
// Returns empty string if unable to fetch.
func GetLatestVersionFromGitHub(repo string, timeoutSecs int) string {
	if repo == "" {
		repo = "microsoft/apm"
	}
	if timeoutSecs <= 0 {
		timeoutSecs = 2
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	client := &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var data struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return ""
	}
	tag := data.TagName
	if len(tag) > 0 && tag[0] == 'v' {
		tag = tag[1:]
	}
	if versionRe.MatchString(tag) {
		return tag
	}
	return ""
}

// ParseVersion parses a semantic version string into components.
// Returns nil if the string is not a valid version.
func ParseVersion(versionStr string) *VersionComponents {
	m := versionRe.FindStringSubmatch(versionStr)
	if m == nil {
		return nil
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	return &VersionComponents{Major: major, Minor: minor, Patch: patch, Prerelease: m[4]}
}

// IsNewerVersion returns true if latest is newer than current.
func IsNewerVersion(current, latest string) bool {
	c := ParseVersion(current)
	l := ParseVersion(latest)
	if c == nil || l == nil {
		return false
	}
	if l.Major != c.Major {
		return l.Major > c.Major
	}
	if l.Minor != c.Minor {
		return l.Minor > c.Minor
	}
	if l.Patch != c.Patch {
		return l.Patch > c.Patch
	}
	// Same major.minor.patch -- compare prerelease
	// Stable (no prerelease) is newer than prerelease
	if l.Prerelease == "" && c.Prerelease != "" {
		return true
	}
	if l.Prerelease != "" && c.Prerelease == "" {
		return false
	}
	return l.Prerelease > c.Prerelease
}

// GetUpdateCachePath returns the path to the version update cache file.
func GetUpdateCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var cacheDir string
	if runtime.GOOS == "windows" {
		cacheDir = filepath.Join(home, "AppData", "Local", "apm", "cache")
	} else {
		cacheDir = filepath.Join(home, ".cache", "apm")
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "last_version_check"), nil
}

// ShouldCheckForUpdates returns true if a version check is due (at most once per day).
func ShouldCheckForUpdates() bool {
	path, err := GetUpdateCachePath()
	if err != nil {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return true // file doesn't exist
	}
	return time.Since(info.ModTime()) > 24*time.Hour
}

// SaveVersionCheckTimestamp saves the timestamp of the last version check.
func SaveVersionCheckTimestamp() {
	path, err := GetUpdateCachePath()
	if err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return
	}
	f.Close()
}

// CheckForUpdates checks if a newer version is available. Returns the latest
// version string if an update is available, empty string otherwise.
func CheckForUpdates(currentVersion string) string {
	if !ShouldCheckForUpdates() {
		return ""
	}
	latest := GetLatestVersionFromGitHub("microsoft/apm", 2)
	SaveVersionCheckTimestamp()
	if latest == "" {
		return ""
	}
	if IsNewerVersion(currentVersion, latest) {
		return latest
	}
	return ""
}
