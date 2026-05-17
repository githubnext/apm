package policy

import "testing"

func TestStripSourcePrefix(t *testing.T) {
	tests := []struct{ in, want string }{
		{"org:myorg", "myorg"},
		{"url:https://example.com", "https://example.com"},
		{"file:/tmp/policy.yml", "/tmp/policy.yml"},
		{"plain", "plain"},
		{"", ""},
	}
	for _, tc := range tests {
		got := stripSourcePrefix(tc.in)
		if got != tc.want {
			t.Errorf("stripSourcePrefix(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		secs int
		want string
	}{
		{-1, "n/a"},
		{0, "0s ago"},
		{30, "30s ago"},
		{59, "59s ago"},
		{60, "1m ago"},
		{3599, "59m ago"},
		{3600, "1h ago"},
		{86399, "23h ago"},
		{86400, "1d ago"},
	}
	for _, tc := range tests {
		got := formatAge(tc.secs)
		if got != tc.want {
			t.Errorf("formatAge(%d) = %q, want %q", tc.secs, got, tc.want)
		}
	}
}
