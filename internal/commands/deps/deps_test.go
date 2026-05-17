package deps

import "testing"

func TestSanitizeMermaid(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"foo/bar", "foo_bar"},
		{"a-b.c@d", "a_b_c_d"},
		{"plain", "plain"},
		{"a/b-c.d@e", "a_b_c_d_e"},
	}
	for _, tc := range tests {
		got := sanitizeMermaid(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeMermaid(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSourceLabel(t *testing.T) {
	tests := []struct {
		dm   map[string]any
		want string
	}{
		{map[string]any{"local": true}, "local"},
		{map[string]any{"host": "dev.azure.com"}, "azure-devops"},
		{map[string]any{"host": "mycompany.visualstudio.com"}, "azure-devops"},
		{map[string]any{"host": "gitlab.com"}, "gitlab"},
		{map[string]any{"host": "github.com"}, "github"},
		{map[string]any{"host": "bitbucket.org"}, "github"},
		{map[string]any{}, "github"},
	}
	for _, tc := range tests {
		got := sourceLabel(tc.dm)
		if got != tc.want {
			t.Errorf("sourceLabel(%v) = %q, want %q", tc.dm, got, tc.want)
		}
	}
}

func TestDepEntryStruct(t *testing.T) {
	e := DepEntry{
		Name:    "mypkg",
		Version: "v1.0.0",
		Source:  "github",
	}
	if e.Name != "mypkg" {
		t.Errorf("unexpected name %q", e.Name)
	}
	if e.IsOrphaned {
		t.Error("expected IsOrphaned false")
	}
}
