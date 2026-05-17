package matcher_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy/matcher"
)

func TestMatchesPattern_Exact(t *testing.T) {
	if !matcher.MatchesPattern("github.com/owner/repo", "github.com/owner/repo") {
		t.Error("exact match should succeed")
	}
}

func TestMatchesPattern_Empty(t *testing.T) {
	if matcher.MatchesPattern("", "pattern") {
		t.Error("empty ref should not match")
	}
	if matcher.MatchesPattern("ref", "") {
		t.Error("empty pattern should not match")
	}
}

func TestMatchesPattern_SingleStar(t *testing.T) {
	if !matcher.MatchesPattern("github.com/owner/repo", "github.com/owner/*") {
		t.Error("single wildcard should match")
	}
	if matcher.MatchesPattern("github.com/owner/sub/nested", "github.com/owner/*") {
		t.Error("single wildcard should not cross /")
	}
}

func TestMatchesPattern_DoubleStar(t *testing.T) {
	if !matcher.MatchesPattern("github.com/owner/sub/nested", "github.com/**") {
		t.Error("double wildcard should match across /")
	}
	if !matcher.MatchesPattern("github.com/a/b/c/d", "github.com/**") {
		t.Error("double wildcard should match deep paths")
	}
}

func TestCheckAllowDeny_NilAllow(t *testing.T) {
	ok, reason := matcher.CheckAllowDeny("any/ref", nil, nil)
	if !ok {
		t.Errorf("nil allow list should allow everything, got reason: %s", reason)
	}
}

func TestCheckAllowDeny_EmptyAllow(t *testing.T) {
	ok, reason := matcher.CheckAllowDeny("any/ref", []string{}, nil)
	if ok {
		t.Error("empty allow list should block all")
	}
	if reason == "" {
		t.Error("should provide reason")
	}
}

func TestCheckAllowDeny_Denied(t *testing.T) {
	ok, reason := matcher.CheckAllowDeny("bad/ref", nil, []string{"bad/*"})
	if ok {
		t.Error("should be denied")
	}
	if reason == "" {
		t.Error("should provide reason")
	}
}

func TestCheckAllowDeny_AllowedByPattern(t *testing.T) {
	ok, _ := matcher.CheckAllowDeny("github.com/owner/repo", []string{"github.com/**"}, nil)
	if !ok {
		t.Error("should be allowed by pattern")
	}
}

func TestCheckAllowDeny_DenyTakesPrecedence(t *testing.T) {
	ok, _ := matcher.CheckAllowDeny("github.com/bad/repo", []string{"github.com/**"}, []string{"github.com/bad/*"})
	if ok {
		t.Error("deny list should take precedence over allow")
	}
}

func TestMatchesPattern_DoubleStarOnly(t *testing.T) {
if !matcher.MatchesPattern("anything/nested/deep", "**") {
t.Error("** should match any path")
}
}

func TestMatchesPattern_TrailingStar(t *testing.T) {
if !matcher.MatchesPattern("github.com/owner/repo", "github.com/**") {
t.Error("trailing ** should match")
}
if !matcher.MatchesPattern("github.com/owner/repo/sub", "github.com/**") {
t.Error("trailing ** should match deep path")
}
}

func TestMatchesPattern_ExactNoWildcard(t *testing.T) {
if matcher.MatchesPattern("github.com/owner/other", "github.com/owner/repo") {
t.Error("exact pattern should not match different path")
}
}

func TestCheckAllowDeny_DenyThenAllow(t *testing.T) {
// Deny overrides allow
ok, reason := matcher.CheckAllowDeny("bad/pkg", []string{"bad/**"}, []string{"bad/*"})
if ok {
t.Error("deny should override allow")
}
if reason == "" {
t.Error("expected non-empty denial reason")
}
}

func TestCheckAllowDeny_NilDeny(t *testing.T) {
ok, _ := matcher.CheckAllowDeny("github.com/ok/repo", []string{"github.com/**"}, nil)
if !ok {
t.Error("should be allowed when deny list is nil")
}
}

func TestCheckAllowDeny_NotInAllowList(t *testing.T) {
ok, reason := matcher.CheckAllowDeny("other.com/owner/repo", []string{"github.com/**"}, nil)
if ok {
t.Error("should be blocked when not in allow list")
}
if reason == "" {
t.Error("expected reason for rejection")
}
}

func TestMatchesPattern_SingleStarInMiddle(t *testing.T) {
if !matcher.MatchesPattern("github.com/owner/repo", "github.com/*/repo") {
t.Error("single wildcard in middle should match")
}
if matcher.MatchesPattern("github.com/a/b/repo", "github.com/*/repo") {
t.Error("single wildcard should not cross /")
}
}
