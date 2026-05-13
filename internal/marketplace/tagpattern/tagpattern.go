// Package tagpattern expands and builds regexes for marketplace version tag patterns.
package tagpattern

import (
"regexp"
"strings"
)

// RenderTag expands {name} and {version} placeholders in pattern.
func RenderTag(pattern, name, version string) string {
result := strings.ReplaceAll(pattern, "{version}", version)
result = strings.ReplaceAll(result, "{name}", name)
return result
}

// BuildTagRegex compiles a tag pattern into a regex that captures the {version} portion.
func BuildTagRegex(pattern string) (*regexp.Regexp, error) {
// Split on {version} to capture it, escape everything else, replace {name} with .+
withName := strings.ReplaceAll(pattern, "{name}", ".+")
parts := strings.SplitN(withName, "{version}", 2)
if len(parts) != 2 {
// No {version} placeholder -- exact match
return regexp.Compile("^" + regexp.QuoteMeta(withName) + "$")
}
re := "^" + regexp.QuoteMeta(parts[0]) + "(?P<version>.+)" + regexp.QuoteMeta(parts[1]) + "$"
return regexp.Compile(re)
}

// ExtractVersion extracts the version from a tag string given a compiled pattern regex.
func ExtractVersion(re *regexp.Regexp, tag string) (string, bool) {
m := re.FindStringSubmatch(tag)
if m == nil {
return "", false
}
for i, name := range re.SubexpNames() {
if name == "version" && i < len(m) {
return m[i], true
}
}
return "", false
}
