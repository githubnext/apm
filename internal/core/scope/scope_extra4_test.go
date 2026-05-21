package scope_test

import (
"testing"

"github.com/githubnext/apm/internal/core/scope"
)

func TestParseScope_CaseInsensitiveExtra4(t *testing.T) {
for _, input := range []string{"USER", "User", "user"} {
s, ok := scope.ParseScope(input)
if !ok {
t.Errorf("expected ok for %q", input)
}
if s != scope.ScopeUser {
t.Errorf("expected ScopeUser for %q", input)
}
}
}

func TestParseScope_ProjectVariantsExtra4(t *testing.T) {
for _, input := range []string{"PROJECT", "Project", "project"} {
s, ok := scope.ParseScope(input)
if !ok {
t.Errorf("expected ok for %q", input)
}
if s != scope.ScopeProject {
t.Errorf("expected ScopeProject for %q", input)
}
}
}

func TestParseScope_InvalidReturnsProjectExtra4(t *testing.T) {
s, ok := scope.ParseScope("global")
if ok {
t.Error("expected ok=false for 'global'")
}
if s != scope.ScopeProject {
t.Errorf("expected ScopeProject default, got %v", s)
}
}

func TestGetDeployRoot_User_NonEmptyExtra4(t *testing.T) {
root, err := scope.GetDeployRoot(scope.ScopeUser)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if root == "" {
t.Error("expected non-empty deploy root for user scope")
}
}
