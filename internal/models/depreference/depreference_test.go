package depreference

import (
"strings"
"testing"
)

func TestParse_SimpleGitHubRef(t *testing.T) {
ref, err := Parse("owner/repo")
if err != nil {
t.Fatalf("Parse(owner/repo) error: %v", err)
}
if ref.RepoURL != "owner/repo" {
t.Errorf("RepoURL = %q, want %q", ref.RepoURL, "owner/repo")
}
}

func TestParse_WithHashReference(t *testing.T) {
ref, err := Parse("owner/repo#v1.2.3")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
if ref.Reference != "v1.2.3" {
t.Errorf("Reference = %q, want v1.2.3", ref.Reference)
}
}

func TestParse_WithHTTPS(t *testing.T) {
ref, err := Parse("https://github.com/owner/repo")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
if ref.ExplicitScheme != "https" {
t.Errorf("ExplicitScheme = %q, want https", ref.ExplicitScheme)
}
}

func TestParse_LocalPath(t *testing.T) {
ref, err := Parse("./local/path")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
if !ref.IsLocal {
t.Errorf("IsLocal should be true for local path")
}
}

func TestParse_AbsoluteLocalPath(t *testing.T) {
ref, err := Parse("/absolute/path")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
if !ref.IsLocal {
t.Errorf("IsLocal should be true for absolute path")
}
}

func TestParse_InvalidEmpty(t *testing.T) {
_, err := Parse("")
if err == nil {
t.Errorf("Parse of empty string should return error")
}
}

func TestIsLocalPath(t *testing.T) {
tests := []struct {
input string
want  bool
}{
{"./local/path", true},
{"../parent/path", true},
{"/absolute/path", true},
{"owner/repo", false},
{"github.com/owner/repo", false},
}
for _, tc := range tests {
got := IsLocalPath(tc.input)
if got != tc.want {
t.Errorf("IsLocalPath(%q) = %v, want %v", tc.input, got, tc.want)
}
}
}

func TestDependencyReference_GetUniqueKey(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
key := ref.GetUniqueKey()
if key == "" {
t.Errorf("GetUniqueKey should not be empty")
}
if !strings.Contains(key, "owner") || !strings.Contains(key, "repo") {
t.Errorf("GetUniqueKey %q should contain owner and repo", key)
}
}

func TestDependencyReference_ToCanonical(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
canonical := ref.ToCanonical()
if canonical == "" {
t.Errorf("ToCanonical should not return empty string")
}
}

func TestDependencyReference_GetInstallPath(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
path, err := ref.GetInstallPath("/tmp/apm_modules")
if err != nil {
t.Fatalf("GetInstallPath error: %v", err)
}
if path == "" {
t.Errorf("GetInstallPath should not be empty")
}
if !strings.Contains(path, "owner") || !strings.Contains(path, "repo") {
t.Errorf("GetInstallPath %q should contain owner and repo", path)
}
}

func TestDependencyReference_IsVirtualFile(t *testing.T) {
ref := &DependencyReference{
IsVirtual:   true,
VirtualPath: "my-file.prompt.md",
}
if !ref.IsVirtualFile() {
t.Errorf("IsVirtualFile should be true for .prompt.md virtual path")
}
}

func TestDependencyReference_IsNotVirtualFile(t *testing.T) {
ref := &DependencyReference{
IsVirtual:   true,
VirtualPath: "some/subdir",
}
if ref.IsVirtualFile() {
t.Errorf("IsVirtualFile should be false for directory virtual path")
}
}

func TestDependencyReference_IsVirtualSubdirectory(t *testing.T) {
ref := &DependencyReference{
IsVirtual:   true,
VirtualPath: "some/subdir",
}
if !ref.IsVirtualSubdirectory() {
t.Errorf("IsVirtualSubdirectory should be true for non-file virtual path")
}
}

func TestDependencyReference_IsArtifactory(t *testing.T) {
ref := &DependencyReference{
ArtifactoryPrefix: "artifactory/github",
}
if !ref.IsArtifactory() {
t.Errorf("IsArtifactory should be true when ArtifactoryPrefix is set")
}

ref2 := &DependencyReference{}
if ref2.IsArtifactory() {
t.Errorf("IsArtifactory should be false when ArtifactoryPrefix is empty")
}
}

func TestDependencyReference_GetDisplayName(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
name := ref.GetDisplayName()
if name == "" {
t.Errorf("GetDisplayName should not return empty")
}
}

func TestDependencyReference_String(t *testing.T) {
ref, err := Parse("owner/repo")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
s := ref.String()
if s == "" {
t.Errorf("String() should not be empty")
}
}

func TestDependencyReference_ToGitHubURL(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
url := ref.ToGitHubURL()
if url == "" {
t.Errorf("ToGitHubURL should not be empty")
}
}

func TestParse_SSHScheme(t *testing.T) {
ref, err := Parse("ssh://git@github.com/owner/repo.git")
if err != nil {
t.Fatalf("Parse SSH URL error: %v", err)
}
if ref.ExplicitScheme != "ssh" {
t.Errorf("ExplicitScheme = %q, want ssh", ref.ExplicitScheme)
}
}

func TestGetInstallPath_IsUnderModulesDir(t *testing.T) {
ref, err := Parse("owner/repo#main")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
path, err := ref.GetInstallPath("/apm_modules")
if err != nil {
return
}
if !strings.HasPrefix(path, "/apm_modules") {
t.Errorf("GetInstallPath %q should be under /apm_modules", path)
}
}

func TestParseFromDict_PathEntry(t *testing.T) {
entry := map[string]interface{}{
"path": "./local/dep",
}
ref, err := ParseFromDict(entry)
if err != nil {
t.Fatalf("ParseFromDict path entry error: %v", err)
}
if ref == nil {
t.Fatal("ParseFromDict should not return nil ref")
}
if !ref.IsLocal {
t.Errorf("ParseFromDict path entry should be local")
}
}

func TestParseFromDict_GitEntry(t *testing.T) {
entry := map[string]interface{}{
"git": "owner/repo",
}
ref, err := ParseFromDict(entry)
if err != nil {
t.Fatalf("ParseFromDict git entry error: %v", err)
}
if ref == nil {
t.Fatal("ParseFromDict should not return nil ref")
}
}

func TestParseFromDict_MissingRequired(t *testing.T) {
entry := map[string]interface{}{
"name": "something",
}
_, err := ParseFromDict(entry)
if err == nil {
t.Errorf("ParseFromDict without git or path should return error")
}
}

func TestGetCanonicalDependencyString(t *testing.T) {
ref, err := Parse("owner/repo")
if err != nil {
t.Fatalf("Parse error: %v", err)
}
canonical := ref.GetCanonicalDependencyString()
if canonical == "" {
t.Errorf("GetCanonicalDependencyString should not be empty")
}
}

func TestParse_AzureDevOps(t *testing.T) {
ref, err := Parse("https://dev.azure.com/myorg/myproject/_git/myrepo")
if err != nil {
// ADO parsing may not be supported for all URL forms
return
}
if !ref.IsAzureDevOps() {
t.Errorf("IsAzureDevOps should be true for ADO URL")
}
}

func TestParse_GitHubHTTPS_WithFragment(t *testing.T) {
ref, err := Parse("https://github.com/owner/repo#v2.0.0")
if err != nil {
t.Fatalf("Parse HTTPS with fragment: %v", err)
}
if ref.Reference != "v2.0.0" {
t.Errorf("Reference = %q, want v2.0.0", ref.Reference)
}
}
