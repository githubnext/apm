package discovery

import (
	"testing"

	"github.com/githubnext/apm/internal/primitives/primmodels"
)

func TestNewPrimitiveCollection(t *testing.T) {
	c := NewPrimitiveCollection()
	if c == nil {
		t.Fatal("nil collection")
	}
	if len(c.Chatmodes) != 0 || len(c.Instructions) != 0 || len(c.Contexts) != 0 || len(c.Skills) != 0 {
		t.Error("expected empty slices")
	}
}

func TestAddPrimitive_Chatmode(t *testing.T) {
	c := NewPrimitiveCollection()
	cm := &primmodels.Chatmode{Name: "test", Source: "local"}
	if err := c.AddPrimitive(cm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Chatmodes) != 1 {
		t.Errorf("want 1 chatmode, got %d", len(c.Chatmodes))
	}
}

func TestAddPrimitive_Instruction(t *testing.T) {
	c := NewPrimitiveCollection()
	i := &primmodels.Instruction{Name: "ins", Source: "local"}
	if err := c.AddPrimitive(i); err != nil {
		t.Fatal(err)
	}
	if len(c.Instructions) != 1 {
		t.Errorf("want 1 instruction, got %d", len(c.Instructions))
	}
}

func TestAddPrimitive_Context(t *testing.T) {
	c := NewPrimitiveCollection()
	ctx := &primmodels.Context{Name: "ctx", Source: "local"}
	if err := c.AddPrimitive(ctx); err != nil {
		t.Fatal(err)
	}
	if len(c.Contexts) != 1 {
		t.Errorf("want 1 context, got %d", len(c.Contexts))
	}
}

func TestAddPrimitive_Skill(t *testing.T) {
	c := NewPrimitiveCollection()
	s := &primmodels.Skill{Name: "sk", Source: "local"}
	if err := c.AddPrimitive(s); err != nil {
		t.Fatal(err)
	}
	if len(c.Skills) != 1 {
		t.Errorf("want 1 skill, got %d", len(c.Skills))
	}
}

type unknownPrimitive struct{}

func (u *unknownPrimitive) Validate() []string { return nil }

func TestAddPrimitive_Unknown(t *testing.T) {
	c := NewPrimitiveCollection()
	err := c.AddPrimitive(&unknownPrimitive{})
	if err == nil {
		t.Error("expected error for unknown primitive type")
	}
}

func TestAddPrimitive_ConflictLocalWins(t *testing.T) {
	c := NewPrimitiveCollection()
	dep := &primmodels.Chatmode{Name: "chat", Source: "dependency:org/repo"}
	local := &primmodels.Chatmode{Name: "chat", Source: "local"}

	c.AddPrimitive(dep)   //nolint:errcheck
	c.AddPrimitive(local) //nolint:errcheck

	if len(c.Chatmodes) != 1 {
		t.Fatalf("want 1 chatmode, got %d", len(c.Chatmodes))
	}
	if c.Chatmodes[0].Source != "local" {
		t.Errorf("expected local to win, got %s", c.Chatmodes[0].Source)
	}
	if len(c.Conflicts) != 1 {
		t.Errorf("want 1 conflict, got %d", len(c.Conflicts))
	}
}

func TestAddPrimitive_ConflictDepDoesNotReplaceLocal(t *testing.T) {
	c := NewPrimitiveCollection()
	local := &primmodels.Chatmode{Name: "chat", Source: "local"}
	dep := &primmodels.Chatmode{Name: "chat", Source: "dependency:org/repo"}

	c.AddPrimitive(local) //nolint:errcheck
	c.AddPrimitive(dep)   //nolint:errcheck

	if c.Chatmodes[0].Source != "local" {
		t.Errorf("expected local to remain, got %s", c.Chatmodes[0].Source)
	}
	if len(c.Conflicts) != 1 {
		t.Errorf("want 1 conflict, got %d", len(c.Conflicts))
	}
}

func TestGlobMatch(t *testing.T) {
	tests := []struct {
		path    string
		pattern string
		want    bool
	}{
		{"foo/bar.chatmode.md", "**/*.chatmode.md", true},
		{"bar.chatmode.md", "**/*.chatmode.md", true},
		{"foo/bar.txt", "**/*.chatmode.md", false},
		{".apm/chatmodes/x.chatmode.md", "**/.apm/chatmodes/*.chatmode.md", true},
		{"a/b/c.instructions.md", "**/*.instructions.md", true},
	}
	for _, tc := range tests {
		got := globMatch(tc.path, tc.pattern)
		if got != tc.want {
			t.Errorf("globMatch(%q, %q) = %v, want %v", tc.path, tc.pattern, got, tc.want)
		}
	}
}
