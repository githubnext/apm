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

func TestTreeNode_Fields(t *testing.T) {
node := TreeNode{
Name:    "mypkg",
Version: "v2.0.0",
Children: []TreeNode{
{Name: "child", Version: "v1.0.0"},
},
}
if node.Name != "mypkg" {
t.Errorf("unexpected Name %q", node.Name)
}
if len(node.Children) != 1 {
t.Errorf("expected 1 child, got %d", len(node.Children))
}
}

func TestDepEntry_InsecureFlag(t *testing.T) {
e := DepEntry{
Name:       "insecure-pkg",
IsInsecure: true,
}
if !e.IsInsecure {
t.Error("expected IsInsecure true")
}
}

func TestDepEntry_Primitives(t *testing.T) {
e := DepEntry{
Name:       "mypkg",
Primitives: []string{"skills", "instructions"},
}
if len(e.Primitives) != 2 {
t.Errorf("expected 2 primitives, got %d", len(e.Primitives))
}
if e.Primitives[0] != "skills" {
t.Errorf("unexpected primitive[0]: %q", e.Primitives[0])
}
}

func TestSanitizeMermaid_SpecialChars(t *testing.T) {
cases := []struct{ in, want string }{
{"my/pkg", "my_pkg"},
{"v1.2.3", "v1_2_3"},
{"@scope/name", "_scope_name"},
{"simple", "simple"},
}
for _, tc := range cases {
got := sanitizeMermaid(tc.in)
if got != tc.want {
t.Errorf("sanitizeMermaid(%q) = %q, want %q", tc.in, got, tc.want)
}
}
}

func TestSourceLabel_Extended(t *testing.T) {
cases := []struct {
dm   map[string]any
want string
}{
{map[string]any{"host": "gitlab.example.com"}, "gitlab"},
{map[string]any{"host": "dev.azure.com/org"}, "azure-devops"},
{map[string]any{"local": true, "host": "github.com"}, "local"},
}
for _, tc := range cases {
got := sourceLabel(tc.dm)
if got != tc.want {
t.Errorf("sourceLabel(%v) = %q, want %q", tc.dm, got, tc.want)
}
}
}
