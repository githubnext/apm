package deps

import "testing"

func TestDepEntry_NameField_Extra4(t *testing.T) {
e := DepEntry{Name: "org/myrepo"}
if e.Name != "org/myrepo" {
t.Errorf("unexpected Name: %s", e.Name)
}
}

func TestDepEntry_VersionField_Extra4(t *testing.T) {
e := DepEntry{Version: "v1.2.3"}
if e.Version != "v1.2.3" {
t.Errorf("unexpected Version: %s", e.Version)
}
}

func TestDepEntry_CommitField_Extra4(t *testing.T) {
e := DepEntry{Commit: "abc1234567"}
if e.Commit != "abc1234567" {
t.Errorf("unexpected Commit: %s", e.Commit)
}
}

func TestDepEntry_RefField_Extra4(t *testing.T) {
e := DepEntry{Ref: "main"}
if e.Ref != "main" {
t.Errorf("unexpected Ref: %s", e.Ref)
}
}

func TestDepEntry_SourceField_Extra4(t *testing.T) {
e := DepEntry{Source: "github"}
if e.Source != "github" {
t.Errorf("unexpected Source: %s", e.Source)
}
}

func TestDepEntry_IsInsecureField_Extra4b(t *testing.T) {
e := DepEntry{IsInsecure: true}
if !e.IsInsecure {
t.Error("expected IsInsecure true")
}
}

func TestCheckIssue_NameField_Extra4(t *testing.T) {
ci := CheckIssue{Name: "org/repo"}
if ci.Name != "org/repo" {
t.Errorf("unexpected Name: %s", ci.Name)
}
}

func TestCheckIssue_ProblemField_Extra4(t *testing.T) {
ci := CheckIssue{Problem: "insecure dependency"}
if ci.Problem != "insecure dependency" {
t.Errorf("unexpected Problem: %s", ci.Problem)
}
}

func TestCheckResult_OKField_Extra4(t *testing.T) {
cr := CheckResult{OK: true}
if !cr.OK {
t.Error("expected OK true")
}
}

func TestListOptions_ScopeField_Extra4(t *testing.T) {
opts := ListOptions{Scope: "dev"}
if opts.Scope != "dev" {
t.Errorf("unexpected Scope: %s", opts.Scope)
}
}

func TestSyncOptions_DryRunField_Extra4(t *testing.T) {
opts := SyncOptions{DryRun: true}
if !opts.DryRun {
t.Error("expected DryRun true")
}
}

func TestOrphanOptions_RemoveField_Extra4(t *testing.T) {
opts := OrphanOptions{Remove: true}
if !opts.Remove {
t.Error("expected Remove true")
}
}
