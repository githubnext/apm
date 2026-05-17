package update

import "testing"

func TestShortSHA(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"abc1234def", "abc1234"},
		{"abc1234", "abc1234"},
		{"abc12", "abc12"},
		{"", ""},
	}
	for _, tc := range tests {
		got := shortSHA(tc.in)
		if got != tc.want {
			t.Errorf("shortSHA(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRenderPlanEntry(t *testing.T) {
	tests := []struct {
		e    PlanEntry
		want string
	}{
		{
			PlanEntry{Package: "mypkg", NewRef: "v1.0.0", ChangeType: "added"},
			"[+] mypkg  (new: v1.0.0)",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "v1.0.0", ChangeType: "removed"},
			"[-] mypkg  (was: v1.0.0)",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "v1.0.0", NewRef: "v2.0.0", ChangeType: "updated"},
			"[~] mypkg  v1.0.0  ->  v2.0.0",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "main", NewRef: "main", OldSHA: "abc1234def", NewSHA: "xyz5678abc", ChangeType: "updated"},
			"[~] mypkg  abc1234  ->  xyz5678",
		},
	}
	for _, tc := range tests {
		got := renderPlanEntry(tc.e)
		if got != tc.want {
			t.Errorf("renderPlanEntry(%+v) = %q, want %q", tc.e, got, tc.want)
		}
	}
}
