// Package lockfile provides APM lock file structures for reproducible installs.
//
// Migrated from src/apm_cli/deps/lockfile.py
package lockfile

import (
"bufio"
"fmt"
"os"
"path/filepath"
"sort"
"strconv"
"strings"
"time"
)

const (
LockfileName       = "apm.lock.yaml"
LegacyLockfileName = "apm.lock"
selfKey            = "."
)

// LockedDependency represents a resolved dependency with exact version info.
type LockedDependency struct {
RepoURL              string
Host                 string
Port                 int // 0 = unset
RegistryPrefix       string
ResolvedCommit       string
ResolvedRef          string
Version              string
VirtualPath          string
IsVirtual            bool
Depth                int
ResolvedBy           string
PackageType          string
DeployedFiles        []string
DeployedFileHashes   map[string]string
Source               string // "local" for local deps
LocalPath            string
ContentHash          string
IsDev                bool
DiscoveredVia        string
MarketplacePluginName string
IsInsecure           bool
AllowInsecure        bool
SkillSubset          []string
}

// GetUniqueKey returns the unique key for this dependency.
func (d *LockedDependency) GetUniqueKey() string {
if d.Source == "local" && d.LocalPath != "" {
return d.LocalPath
}
if d.IsVirtual && d.VirtualPath != "" {
return d.RepoURL + "/" + d.VirtualPath
}
return d.RepoURL
}

// ToDict serializes the dependency to a string map for YAML output.
func (d *LockedDependency) ToDict() map[string]interface{} {
result := map[string]interface{}{"repo_url": d.RepoURL}
if d.Host != "" {
result["host"] = d.Host
}
if d.Port != 0 {
result["port"] = d.Port
}
if d.RegistryPrefix != "" {
result["registry_prefix"] = d.RegistryPrefix
}
if d.ResolvedCommit != "" {
result["resolved_commit"] = d.ResolvedCommit
}
if d.ResolvedRef != "" {
result["resolved_ref"] = d.ResolvedRef
}
if d.Version != "" {
result["version"] = d.Version
}
if d.VirtualPath != "" {
result["virtual_path"] = d.VirtualPath
}
if d.IsVirtual {
result["is_virtual"] = true
}
if d.Depth != 1 {
result["depth"] = d.Depth
}
if d.ResolvedBy != "" {
result["resolved_by"] = d.ResolvedBy
}
if d.PackageType != "" {
result["package_type"] = d.PackageType
}
if len(d.DeployedFiles) > 0 {
sorted := append([]string{}, d.DeployedFiles...)
sort.Strings(sorted)
result["deployed_files"] = sorted
}
if len(d.DeployedFileHashes) > 0 {
result["deployed_file_hashes"] = sortedMapCopy(d.DeployedFileHashes)
}
if d.Source != "" {
result["source"] = d.Source
}
if d.LocalPath != "" {
result["local_path"] = d.LocalPath
}
if d.ContentHash != "" {
result["content_hash"] = d.ContentHash
}
if d.IsDev {
result["is_dev"] = true
}
if d.DiscoveredVia != "" {
result["discovered_via"] = d.DiscoveredVia
}
if d.MarketplacePluginName != "" {
result["marketplace_plugin_name"] = d.MarketplacePluginName
}
if d.IsInsecure {
result["is_insecure"] = true
}
if d.AllowInsecure {
result["allow_insecure"] = true
}
if len(d.SkillSubset) > 0 {
sorted := append([]string{}, d.SkillSubset...)
sort.Strings(sorted)
result["skill_subset"] = sorted
}
return result
}

// LockedDepFromMap deserializes a LockedDependency from a parsed YAML map.
func LockedDepFromMap(data map[string]interface{}) (*LockedDependency, error) {
repoURL, ok := data["repo_url"].(string)
if !ok || repoURL == "" {
return nil, fmt.Errorf("missing repo_url")
}

deployedFiles := strSlice(data["deployed_files"])
// Migrate legacy deployed_skills -> deployed_files
if oldSkills := strSlice(data["deployed_skills"]); len(oldSkills) > 0 && len(deployedFiles) == 0 {
for _, sk := range oldSkills {
deployedFiles = append(deployedFiles, ".github/skills/"+sk+"/")
deployedFiles = append(deployedFiles, ".claude/skills/"+sk+"/")
}
}

var port int
if pRaw, ok := data["port"]; ok && pRaw != nil {
switch v := pRaw.(type) {
case int:
if v >= 1 && v <= 65535 {
port = v
}
case float64:
p := int(v)
if p >= 1 && p <= 65535 {
port = p
}
case string:
if p, err := strconv.Atoi(v); err == nil && p >= 1 && p <= 65535 {
port = p
}
}
}

dep := &LockedDependency{
RepoURL:              repoURL,
Host:                 strVal(data["host"]),
Port:                 port,
RegistryPrefix:       strVal(data["registry_prefix"]),
ResolvedCommit:       strVal(data["resolved_commit"]),
ResolvedRef:          strVal(data["resolved_ref"]),
Version:              strVal(data["version"]),
VirtualPath:          strVal(data["virtual_path"]),
IsVirtual:            boolVal(data["is_virtual"]),
Depth:                intVal(data["depth"], 1),
ResolvedBy:           strVal(data["resolved_by"]),
PackageType:          strVal(data["package_type"]),
DeployedFiles:        deployedFiles,
DeployedFileHashes:   strMap(data["deployed_file_hashes"]),
Source:               strVal(data["source"]),
LocalPath:            strVal(data["local_path"]),
ContentHash:          strVal(data["content_hash"]),
IsDev:                boolVal(data["is_dev"]),
DiscoveredVia:        strVal(data["discovered_via"]),
MarketplacePluginName: strVal(data["marketplace_plugin_name"]),
IsInsecure:           boolVal(data["is_insecure"]),
AllowInsecure:        boolVal(data["allow_insecure"]),
SkillSubset:          strSlice(data["skill_subset"]),
}
return dep, nil
}

// LockFile represents an APM lock file.
type LockFile struct {
LockfileVersion        string
GeneratedAt            string
APMVersion             string
Dependencies           map[string]*LockedDependency
MCPServers             []string
MCPConfigs             map[string]map[string]interface{}
LocalDeployedFiles     []string
LocalDeployedFileHashes map[string]string
}

// NewLockFile creates a new empty LockFile.
func NewLockFile() *LockFile {
return &LockFile{
LockfileVersion:        "1",
GeneratedAt:            time.Now().UTC().Format(time.RFC3339),
Dependencies:           make(map[string]*LockedDependency),
MCPConfigs:             make(map[string]map[string]interface{}),
LocalDeployedFileHashes: make(map[string]string),
}
}

// AddDependency adds a dependency to the lock file.
func (lf *LockFile) AddDependency(dep *LockedDependency) {
lf.Dependencies[dep.GetUniqueKey()] = dep
}

// GetDependency returns a dependency by key.
func (lf *LockFile) GetDependency(key string) *LockedDependency {
return lf.Dependencies[key]
}

// HasDependency checks if a dependency exists.
func (lf *LockFile) HasDependency(key string) bool {
_, ok := lf.Dependencies[key]
return ok
}

// GetAllDependencies returns all dependencies sorted by depth then repo_url.
func (lf *LockFile) GetAllDependencies() []*LockedDependency {
deps := make([]*LockedDependency, 0, len(lf.Dependencies))
for _, d := range lf.Dependencies {
deps = append(deps, d)
}
sort.Slice(deps, func(i, j int) bool {
if deps[i].Depth != deps[j].Depth {
return deps[i].Depth < deps[j].Depth
}
return deps[i].RepoURL < deps[j].RepoURL
})
return deps
}

// GetPackageDependencies returns all dependencies excluding the virtual self-entry.
func (lf *LockFile) GetPackageDependencies() []*LockedDependency {
var result []*LockedDependency
for _, d := range lf.GetAllDependencies() {
if d.LocalPath != "." {
result = append(result, d)
}
}
return result
}

// IsSemanticalllyEquivalent returns true if other has the same deps/MCP/configs.
func (lf *LockFile) IsSemanticalllyEquivalent(other *LockFile) bool {
if lf.LockfileVersion != other.LockfileVersion {
return false
}
if len(lf.Dependencies) != len(other.Dependencies) {
return false
}
for key, dep := range lf.Dependencies {
od, ok := other.Dependencies[key]
if !ok {
return false
}
if fmt.Sprint(dep.ToDict()) != fmt.Sprint(od.ToDict()) {
return false
}
}
// MCP servers
as := append([]string{}, lf.MCPServers...)
bs := append([]string{}, other.MCPServers...)
sort.Strings(as)
sort.Strings(bs)
if strings.Join(as, ",") != strings.Join(bs, ",") {
return false
}
if fmt.Sprint(lf.MCPConfigs) != fmt.Sprint(other.MCPConfigs) {
return false
}
af := append([]string{}, lf.LocalDeployedFiles...)
bf := append([]string{}, other.LocalDeployedFiles...)
sort.Strings(af)
sort.Strings(bf)
if strings.Join(af, ",") != strings.Join(bf, ",") {
return false
}
return fmt.Sprint(sortedMapCopy(lf.LocalDeployedFileHashes)) == fmt.Sprint(sortedMapCopy(other.LocalDeployedFileHashes))
}

// FromYAML parses a LockFile from a simple line-by-line YAML reader.
// This is a minimal parser for the known lockfile schema.
func FromYAML(content string) (*LockFile, error) {
lf := NewLockFile()
scanner := bufio.NewScanner(strings.NewReader(content))
var lines []string
for scanner.Scan() {
lines = append(lines, scanner.Text())
}

// Simple state machine parser
i := 0
for i < len(lines) {
line := lines[i]
trimmed := strings.TrimSpace(line)

if strings.HasPrefix(trimmed, "lockfile_version:") {
lf.LockfileVersion = yamlValue(trimmed)
i++
} else if strings.HasPrefix(trimmed, "generated_at:") {
lf.GeneratedAt = yamlValue(trimmed)
i++
} else if strings.HasPrefix(trimmed, "apm_version:") {
lf.APMVersion = yamlValue(trimmed)
i++
} else if trimmed == "dependencies:" {
i++
// Parse list of dependency maps
for i < len(lines) {
dl := lines[i]
dtrimmed := strings.TrimSpace(dl)
if strings.HasPrefix(dtrimmed, "- repo_url:") || dtrimmed == "-" {
depMap, n := parseYAMLMap(lines, i)
i += n
dep, err := LockedDepFromMap(depMap)
if err == nil {
lf.AddDependency(dep)
}
} else if !strings.HasPrefix(dl, " ") && !strings.HasPrefix(dl, "\t") && dl != "" {
break
} else {
i++
}
}
} else if trimmed == "mcp_servers:" {
i++
for i < len(lines) {
sl := strings.TrimSpace(lines[i])
if strings.HasPrefix(sl, "- ") {
lf.MCPServers = append(lf.MCPServers, strings.TrimPrefix(sl, "- "))
i++
} else if sl == "" || !strings.HasPrefix(lines[i], " ") {
break
} else {
i++
}
}
} else if trimmed == "local_deployed_files:" {
i++
for i < len(lines) {
sl := strings.TrimSpace(lines[i])
if strings.HasPrefix(sl, "- ") {
lf.LocalDeployedFiles = append(lf.LocalDeployedFiles, strings.TrimPrefix(sl, "- "))
i++
} else if sl == "" || !strings.HasPrefix(lines[i], " ") {
break
} else {
i++
}
}
} else if trimmed == "local_deployed_file_hashes:" {
i++
for i < len(lines) {
kl := lines[i]
ktrimmed := strings.TrimSpace(kl)
if strings.HasPrefix(lines[i], "  ") && strings.Contains(ktrimmed, ":") {
parts := strings.SplitN(ktrimmed, ":", 2)
if len(parts) == 2 {
k := strings.Trim(strings.TrimSpace(parts[0]), `"'`)
v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
lf.LocalDeployedFileHashes[k] = v
}
i++
} else {
break
}
}
} else {
i++
}
}

// Synthesize self-entry
if len(lf.LocalDeployedFiles) > 0 {
lf.Dependencies[selfKey] = &LockedDependency{
RepoURL:            "<self>",
Source:             "local",
LocalPath:          ".",
IsDev:              true,
Depth:              0,
DeployedFiles:      append([]string{}, lf.LocalDeployedFiles...),
DeployedFileHashes: copyStrMap(lf.LocalDeployedFileHashes),
}
}

return lf, nil
}

// GetLockfilePath returns the path to the lock file for a project.
func GetLockfilePath(projectRoot string) string {
return filepath.Join(projectRoot, LockfileName)
}

// MigrateLockfileIfNeeded renames legacy apm.lock to apm.lock.yaml.
func MigrateLockfileIfNeeded(projectRoot string) bool {
newPath := GetLockfilePath(projectRoot)
legacyPath := filepath.Join(projectRoot, LegacyLockfileName)
if _, err := os.Stat(newPath); os.IsNotExist(err) {
if _, err2 := os.Stat(legacyPath); err2 == nil {
if err3 := os.Rename(legacyPath, newPath); err3 == nil {
return true
}
}
}
return false
}

// ReadLockfile reads a lock file from disk.
func ReadLockfile(path string) (*LockFile, error) {
data, err := os.ReadFile(path)
if err != nil {
return nil, err
}
return FromYAML(string(data))
}

// LoadOrCreate loads a lock file or creates a new one.
func LoadOrCreate(path string) *LockFile {
lf, err := ReadLockfile(path)
if err != nil || lf == nil {
return NewLockFile()
}
return lf
}

// --- YAML parsing helpers ---

// parseYAMLMap parses a YAML list item (map) starting at lines[start].
// Returns the map and the number of lines consumed.
func parseYAMLMap(lines []string, start int) (map[string]interface{}, int) {
result := make(map[string]interface{})
i := start

// Consume leading "- " prefix on first line
firstLine := strings.TrimSpace(lines[i])
if strings.HasPrefix(firstLine, "- ") {
kv := strings.TrimPrefix(firstLine, "- ")
if strings.Contains(kv, ":") {
parts := strings.SplitN(kv, ":", 2)
k := strings.TrimSpace(parts[0])
v := strings.TrimSpace(parts[1])
result[k] = unquote(v)
}
i++
} else if firstLine == "-" {
i++
}

// indent of the block items
blockIndent := ""
for i < len(lines) {
line := lines[i]
if strings.TrimSpace(line) == "" {
i++
continue
}
// Detect indentation
for _, c := range line {
if c == ' ' {
blockIndent += " "
} else {
break
}
}
break
}
if blockIndent == "" {
blockIndent = "  "
}

for i < len(lines) {
line := lines[i]
trimmed := strings.TrimSpace(line)

if trimmed == "" {
i++
continue
}
// End of this map item
if strings.HasPrefix(trimmed, "- ") || (!strings.HasPrefix(line, blockIndent) && !strings.HasPrefix(line, " ")) {
break
}
// Nested list
if strings.Contains(trimmed, ":") {
parts := strings.SplitN(trimmed, ":", 2)
key := strings.TrimSpace(parts[0])
val := strings.TrimSpace(parts[1])
if val == "" {
// collect sub-list or sub-map
i++
var subList []string
subMap := make(map[string]interface{})
isList := false
for i < len(lines) {
sl := lines[i]
strimmed := strings.TrimSpace(sl)
if strimmed == "" {
i++
continue
}
if !strings.HasPrefix(sl, blockIndent) {
break
}
if strings.HasPrefix(strimmed, "- ") {
isList = true
subList = append(subList, strings.TrimPrefix(strimmed, "- "))
i++
} else if strings.Contains(strimmed, ":") {
kp := strings.SplitN(strimmed, ":", 2)
sk := strings.TrimSpace(kp[0])
sv := strings.Trim(strings.TrimSpace(kp[1]), `"'`)
subMap[sk] = sv
i++
} else {
break
}
}
if isList {
result[key] = subList
} else {
result[key] = subMap
}
continue
}
result[key] = parseScalar(val)
i++
} else {
i++
}
}
return result, i - start
}

func yamlValue(line string) string {
idx := strings.Index(line, ":")
if idx < 0 {
return ""
}
return strings.Trim(strings.TrimSpace(line[idx+1:]), `"'`)
}

func unquote(s string) interface{} {
s = strings.TrimSpace(s)
if s == "" {
return nil
}
return parseScalar(s)
}

func parseScalar(s string) interface{} {
s = strings.Trim(s, `"'`)
if s == "true" {
return true
}
if s == "false" {
return false
}
if s == "null" || s == "~" {
return nil
}
if n, err := strconv.Atoi(s); err == nil {
return n
}
if f, err := strconv.ParseFloat(s, 64); err == nil {
return f
}
return s
}

// --- type coercion helpers ---

func strVal(v interface{}) string {
if v == nil {
return ""
}
if s, ok := v.(string); ok {
return s
}
return fmt.Sprint(v)
}

func boolVal(v interface{}) bool {
if v == nil {
return false
}
b, ok := v.(bool)
return ok && b
}

func intVal(v interface{}, def int) int {
if v == nil {
return def
}
switch n := v.(type) {
case int:
return n
case float64:
return int(n)
}
return def
}

func strSlice(v interface{}) []string {
if v == nil {
return nil
}
switch s := v.(type) {
case []string:
return s
case []interface{}:
result := make([]string, 0, len(s))
for _, item := range s {
result = append(result, strVal(item))
}
return result
}
return nil
}

func strMap(v interface{}) map[string]string {
if v == nil {
return make(map[string]string)
}
switch m := v.(type) {
case map[string]string:
return m
case map[string]interface{}:
result := make(map[string]string, len(m))
for k, val := range m {
result[k] = strVal(val)
}
return result
case map[interface{}]interface{}:
result := make(map[string]string, len(m))
for k, val := range m {
result[strVal(k)] = strVal(val)
}
return result
}
return make(map[string]string)
}

func sortedMapCopy(m map[string]string) map[string]string {
result := make(map[string]string, len(m))
for k, v := range m {
result[k] = v
}
return result
}

func copyStrMap(m map[string]string) map[string]string {
result := make(map[string]string, len(m))
for k, v := range m {
result[k] = v
}
return result
}
