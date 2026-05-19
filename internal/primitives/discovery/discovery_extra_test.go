package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/primitives/primmodels"
)

func TestShouldReplace_NewLocalReplacesNonLocal(t *testing.T) {
	if !shouldReplace("dep-a", "local") {
		t.Error("local should replace non-local")
	}
}

func TestShouldReplace_LocalDoesNotReplaceLocal(t *testing.T) {
	if shouldReplace("local", "local") {
		t.Error("local should not replace local")
	}
}

func TestShouldReplace_NonLocalDoesNotReplaceLocal(t *testing.T) {
	if shouldReplace("local", "dep-b") {
		t.Error("non-local should not replace local")
	}
}

func TestShouldReplace_EmptyTreatedAsLocal(t *testing.T) {
	if !shouldReplace("dep-a", "") {
		t.Error("empty source should be treated as local and replace non-local")
	}
}

func TestShouldReplace_BothNonLocal(t *testing.T) {
	if shouldReplace("dep-a", "dep-b") {
		t.Error("non-local should not replace non-local")
	}
}

func TestGlobMatch_SimplePattern(t *testing.T) {
	if !globMatch("foo/bar.md", "foo/bar.md") {
		t.Error("exact match should succeed")
	}
}

func TestGlobMatch_StarPattern(t *testing.T) {
	if !globMatch("foo/bar.md", "foo/*.md") {
		t.Error("star pattern should match")
	}
}

func TestGlobMatch_DoubleStarPattern(t *testing.T) {
	if !globMatch("a/b/c/file.md", "**/*.md") {
		t.Error("double-star pattern should match")
	}
}

func TestGlobMatch_Mismatch(t *testing.T) {
	if globMatch("foo/bar.txt", "foo/*.md") {
		t.Error("should not match different extension")
	}
}

func TestGlobMatch_EmptyPath(t *testing.T) {
	if globMatch("", "foo/bar.md") {
		t.Error("empty path should not match non-empty pattern")
	}
}

func TestAddPrimitive_SkillAndContext(t *testing.T) {
	c := NewPrimitiveCollection()
	sk := &primmodels.Skill{Name: "my-skill", Source: "local"}
	ctx := &primmodels.Context{Name: "my-ctx", Source: "local"}
	if err := c.AddPrimitive(sk); err != nil {
		t.Fatalf("AddPrimitive skill: %v", err)
	}
	if err := c.AddPrimitive(ctx); err != nil {
		t.Fatalf("AddPrimitive context: %v", err)
	}
	if len(c.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(c.Skills))
	}
	if len(c.Contexts) != 1 {
		t.Errorf("expected 1 context, got %d", len(c.Contexts))
	}
}

func TestAddPrimitive_UnknownType(t *testing.T) {
	c := NewPrimitiveCollection()
	err := c.AddPrimitive(&fakeUnknownPrimitive{})
	if err == nil {
		t.Error("expected error for unknown primitive type")
	}
}

type fakeUnknownPrimitive struct{}

func (f *fakeUnknownPrimitive) Validate() []string { return nil }

func TestPrimitiveCollection_Conflict_SameNameDifferentSource(t *testing.T) {
	c := NewPrimitiveCollection()
	cm1 := &primmodels.Chatmode{Name: "agent", Source: "dep-a"}
	cm2 := &primmodels.Chatmode{Name: "agent", Source: "local"}
	c.AddPrimitive(cm1) //nolint
	c.AddPrimitive(cm2) //nolint
	if len(c.Conflicts) == 0 {
		t.Error("expected conflict to be recorded")
	}
}

func TestPrimitiveCollection_MultipleInstructions(t *testing.T) {
	c := NewPrimitiveCollection()
	for _, name := range []string{"ins-a", "ins-b", "ins-c"} {
		ins := &primmodels.Instruction{Name: name, Source: "local"}
		if err := c.AddPrimitive(ins); err != nil {
			t.Fatalf("AddPrimitive: %v", err)
		}
	}
	if len(c.Instructions) != 3 {
		t.Errorf("expected 3 instructions, got %d", len(c.Instructions))
	}
}

func TestPrimitiveConflict_Fields(t *testing.T) {
	pc := PrimitiveConflict{
		PrimitiveName: "my-agent",
		PrimitiveType: "chatmode",
		WinningSource: "local",
		LosingSource:  "dep-b",
		FilePath:      "/path/to/file.md",
	}
	if pc.PrimitiveName != "my-agent" || pc.WinningSource != "local" {
		t.Errorf("unexpected fields: %+v", pc)
	}
}

func TestFindPrimitiveFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files, err := FindPrimitiveFiles(dir, []string{"**/*.md"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestFindPrimitiveFiles_MatchesMD(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".apm", "skills")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "my.skill.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := FindPrimitiveFiles(dir, []string{"**/*.skill.md"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d: %v", len(files), files)
	}
}

func TestDiscoverPrimitives_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	col, err := DiscoverPrimitives(dir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col == nil {
		t.Fatal("expected non-nil collection")
	}
	total := len(col.Chatmodes) + len(col.Instructions) + len(col.Contexts) + len(col.Skills)
	if total != 0 {
		t.Errorf("expected 0 primitives in empty dir, got %d", total)
	}
}
