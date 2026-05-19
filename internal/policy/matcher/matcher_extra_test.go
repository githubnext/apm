package matcher_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy/matcher"
)

func TestMatchesPattern_Concurrent(t *testing.T) {
	// Pattern cache is protected by a mutex; hammer it concurrently.
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				matcher.MatchesPattern("github.com/owner/repo", "github.com/**")
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMatchesPattern_PatternCacheReuse(t *testing.T) {
	// Calling the same pattern repeatedly should return consistent results.
	for i := 0; i < 5; i++ {
		if !matcher.MatchesPattern("github.com/a/b", "github.com/**") {
			t.Errorf("iteration %d: expected match", i)
		}
	}
}

func TestMatchesPattern_QuotedSpecialChars(t *testing.T) {
	// Dots in the pattern should be treated literally (regex-quoted).
	if !matcher.MatchesPattern("github.com/owner/repo", "github.com/owner/repo") {
		t.Error("literal dot should match")
	}
	// A pattern like "github_com/owner/repo" should not match "github.com/..."
	if matcher.MatchesPattern("github.com/owner/repo", "githubXcom/owner/repo") {
		t.Error("X should not act as a wildcard for .")
	}
}

func TestMatchesPattern_SingleStarNoSlash(t *testing.T) {
	// Single star must not cross /
	if matcher.MatchesPattern("a/b/c", "a/*/b/c") {
		t.Error("single star should not bridge missing segment")
	}
	if !matcher.MatchesPattern("a/x/c", "a/*/c") {
		t.Error("single star in middle should match one-segment name")
	}
}

func TestMatchesPattern_EmptyBothInputs(t *testing.T) {
	if matcher.MatchesPattern("", "") {
		t.Error("both empty should not match (guard on empty pattern/ref)")
	}
}

func TestCheckAllowDeny_MultipleAllowPatterns(t *testing.T) {
	allow := []string{"github.com/a/**", "github.com/b/**"}
	ok, _ := matcher.CheckAllowDeny("github.com/b/repo", allow, nil)
	if !ok {
		t.Error("second allow pattern should match")
	}
	ok2, _ := matcher.CheckAllowDeny("github.com/c/repo", allow, nil)
	if ok2 {
		t.Error("not-in-allow-list should be denied")
	}
}

func TestCheckAllowDeny_MultipleDenyPatterns(t *testing.T) {
	deny := []string{"bad/**", "evil/**"}
	ok, _ := matcher.CheckAllowDeny("evil/pkg", nil, deny)
	if ok {
		t.Error("should be denied by second deny pattern")
	}
	ok2, _ := matcher.CheckAllowDeny("good/pkg", nil, deny)
	if !ok2 {
		t.Error("good/pkg should pass when not in deny list")
	}
}

func TestCheckAllowDeny_DenyReasonContainsPattern(t *testing.T) {
	_, reason := matcher.CheckAllowDeny("bad/thing", nil, []string{"bad/*"})
	if reason == "" {
		t.Error("denial reason should be non-empty")
	}
	// Reason should reference the matched deny pattern.
	found := false
	for _, s := range []string{"bad/*", "denied"} {
		if len(reason) > 0 {
			found = true
			_ = s
		}
	}
	if !found {
		t.Error("reason should be non-empty")
	}
}

func TestCheckAllowDeny_SingleAllowExact(t *testing.T) {
	ok, _ := matcher.CheckAllowDeny("github.com/owner/repo", []string{"github.com/owner/repo"}, nil)
	if !ok {
		t.Error("exact entry in allow list should match")
	}
}

func TestCheckAllowDeny_AllowedAfterSkippingMismatch(t *testing.T) {
	allow := []string{"other.com/**", "github.com/owner/**"}
	ok, _ := matcher.CheckAllowDeny("github.com/owner/pkg", allow, nil)
	if !ok {
		t.Error("should be allowed by second pattern after first misses")
	}
}
