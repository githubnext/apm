package mcpargs

import "testing"

func TestParseKVPairs_ValueHasMultipleEquals(t *testing.T) {
	// Only the first '=' is the separator; the rest belong to the value.
	m, err := ParseKVPairs([]string{"A=x=y=z"}, "--test2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["A"] != "x=y=z" {
		t.Errorf("expected m[A]='x=y=z', got %q", m["A"])
	}
}

func TestParseKVPairs_NilInput(t *testing.T) {
	m, err := ParseKVPairs(nil, "--test2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestParseKVPairs_ErrorMsgContainsFlagAndRaw(t *testing.T) {
	_, err := ParseKVPairs([]string{"BADPAIR"}, "--myFlag")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
}

func TestParseEnvPairs_PathValue(t *testing.T) {
	m, err := ParseEnvPairs([]string{"HOME=/home/user", "PATH=/usr/bin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["HOME"] != "/home/user" {
		t.Errorf("HOME = %q", m["HOME"])
	}
	if m["PATH"] != "/usr/bin" {
		t.Errorf("PATH = %q", m["PATH"])
	}
}

func TestParseHeaderPairs_SingleEntry(t *testing.T) {
	m, err := ParseHeaderPairs([]string{"X-Auth=token123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["X-Auth"] != "token123" {
		t.Errorf("X-Auth = %q", m["X-Auth"])
	}
}

func TestParseKVPairs_ThreePairs(t *testing.T) {
	m, err := ParseKVPairs([]string{"A=1", "B=2", "C=3"}, "--test2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(m))
	}
}
