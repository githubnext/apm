package listcmd

import "testing"

func TestScript_CommandField_Extra4(t *testing.T) {
s := Script{Name: "test", Command: "go test ./..."}
if s.Command != "go test ./..." {
t.Errorf("unexpected command: %s", s.Command)
}
}

func TestScript_NameIsKey_Extra4(t *testing.T) {
s := Script{Name: "lint"}
if s.Name != "lint" {
t.Errorf("unexpected name: %s", s.Name)
}
}

func TestScript_EmptyName_Extra4(t *testing.T) {
s := Script{Name: "", Command: "echo hello"}
if s.Name != "" {
t.Errorf("expected empty name, got %s", s.Name)
}
}

func TestScript_SliceOfThree_Extra4(t *testing.T) {
scripts := []Script{
{Name: "a", Command: "cmd1"},
{Name: "b", Command: "cmd2"},
{Name: "c", Command: "cmd3"},
}
if len(scripts) != 3 {
t.Errorf("expected 3 scripts, got %d", len(scripts))
}
}

func TestScript_ZeroValue_Extra4(t *testing.T) {
s := Script{}
if s.Name != "" || s.Command != "" {
t.Errorf("expected zero value, got %+v", s)
}
}

func TestScript_LongCommand_Extra4(t *testing.T) {
cmd := "go build -v -ldflags '-X main.version=1.0' ./cmd/..."
s := Script{Name: "build", Command: cmd}
if s.Command != cmd {
t.Errorf("expected long command preserved, got %s", s.Command)
}
}

func TestScript_BothFields_Extra4(t *testing.T) {
s := Script{Name: "deploy", Command: "kubectl apply -f k8s/"}
if s.Name != "deploy" {
t.Errorf("unexpected name: %s", s.Name)
}
if s.Command != "kubectl apply -f k8s/" {
t.Errorf("unexpected command: %s", s.Command)
}
}

func TestScript_MultipleInstances_Extra4(t *testing.T) {
scripts := make([]Script, 5)
for i := range scripts {
scripts[i] = Script{Name: "script", Command: "echo"}
}
if len(scripts) != 5 {
t.Errorf("expected 5 scripts, got %d", len(scripts))
}
}
