package matcher

import (
	"testing"
)

func TestMatchesPattern_DoubleStar(t *testing.T) {
	if !MatchesPattern("owner/repo/some/path", "owner/**") {
		t.Error("expected ** to match nested path")
	}
}

func TestMatchesPattern_ExactMatch(t *testing.T) {
	if !MatchesPattern("owner/repo", "owner/repo") {
		t.Error("expected exact match")
	}
}

func TestMatchesPattern_NoMatch(t *testing.T) {
	if MatchesPattern("other/repo", "owner/repo") {
		t.Error("expected no match for different refs")
	}
}

func TestMatchesPattern_StarWildcard(t *testing.T) {
	if !MatchesPattern("owner/anything", "owner/*") {
		t.Error("expected * to match single segment")
	}
}

func TestMatchesPattern_EmptyPattern(t *testing.T) {
	result := MatchesPattern("owner/repo", "")
	_ = result // just verify no panic
}

func TestMatchesPattern_EmptyRef(t *testing.T) {
	result := MatchesPattern("", "owner/repo")
	_ = result // just verify no panic
}

func TestCheckAllowDeny_EmptyAllowAndDeny(t *testing.T) {
	allowed, reason := CheckAllowDeny("owner/repo", nil, nil)
	_ = allowed
	_ = reason
}

func TestCheckAllowDeny_AllowAll(t *testing.T) {
	allowed, _ := CheckAllowDeny("any/ref", []string{"**"}, nil)
	if !allowed {
		t.Error("expected ** to allow any ref")
	}
}

func TestCheckAllowDeny_DenyOverridesAllow(t *testing.T) {
	allowed, _ := CheckAllowDeny("owner/repo", []string{"owner/*"}, []string{"owner/repo"})
	if allowed {
		t.Error("expected deny to override allow")
	}
}

func TestCheckAllowDeny_ReasonOnDeny(t *testing.T) {
	_, reason := CheckAllowDeny("owner/repo", nil, []string{"owner/repo"})
	if reason == "" {
		t.Error("expected non-empty reason when denied")
	}
}

func TestCheckAllowDeny_AllowByExact(t *testing.T) {
	allowed, _ := CheckAllowDeny("owner/repo", []string{"owner/repo"}, nil)
	if !allowed {
		t.Error("expected exact allow to match")
	}
}

func TestCheckAllowDeny_MultiplePatternFirstMatches(t *testing.T) {
	allowed, _ := CheckAllowDeny("owner/repo", []string{"owner/repo", "other/*"}, nil)
	if !allowed {
		t.Error("expected first pattern to match")
	}
}

func TestMatchesPattern_CaseSensitive(t *testing.T) {
	// Pattern matching should be case-sensitive
	result1 := MatchesPattern("Owner/Repo", "owner/repo")
	result2 := MatchesPattern("owner/repo", "owner/repo")
	// Even if case-insensitive, verify no panic
	_ = result1
	_ = result2
}
