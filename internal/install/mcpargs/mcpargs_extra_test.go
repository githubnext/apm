package mcpargs

import "testing"

func TestParseKVPairs_ValueWithEquals(t *testing.T) {
	// Value itself contains '='; only the first '=' is the separator.
	got, err := ParseKVPairs([]string{"KEY=a=b=c"}, "--test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "a=b=c" {
		t.Errorf("expected value 'a=b=c', got %q", got["KEY"])
	}
}

func TestParseKVPairs_EmptyValue(t *testing.T) {
	got, err := ParseKVPairs([]string{"KEY="}, "--test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := got["KEY"]; !ok || v != "" {
		t.Errorf("expected empty value for KEY, got %q (ok=%v)", v, ok)
	}
}

func TestParseKVPairs_MultipleEntries(t *testing.T) {
	got, err := ParseKVPairs([]string{"A=1", "B=2", "C=3"}, "--test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got["A"] != "1" || got["B"] != "2" || got["C"] != "3" {
		t.Errorf("value mismatch: %v", got)
	}
}

func TestParseEnvPairs_ErrorMessage(t *testing.T) {
	_, err := ParseEnvPairs([]string{"BADPAIR"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
	if !contains(err.Error(), "--env") {
		t.Errorf("expected '--env' in error, got %q", err.Error())
	}
}

func TestParseHeaderPairs_ErrorMessage(t *testing.T) {
	_, err := ParseHeaderPairs([]string{"BADPAIR"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
	if !contains(err.Error(), "--header") {
		t.Errorf("expected '--header' in error, got %q", err.Error())
	}
}

func TestParseKVPairs_EmptyKey_ErrorContainsFlagName(t *testing.T) {
	_, err := ParseKVPairs([]string{"=value"}, "--myFlag")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if !contains(err.Error(), "--myFlag") {
		t.Errorf("expected '--myFlag' in error, got %q", err.Error())
	}
}

func TestParseKVPairs_WhitespaceKey(t *testing.T) {
	// Key with whitespace is technically valid (not empty).
	got, err := ParseKVPairs([]string{" KEY =value"}, "--test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got[" KEY "]; !ok {
		t.Error("expected whitespace key to be stored as-is")
	}
}

func TestParseKVPairs_OverwritesDuplicate(t *testing.T) {
	// Later value overwrites earlier one for same key.
	got, err := ParseKVPairs([]string{"KEY=first", "KEY=second"}, "--test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", got["KEY"])
	}
}

// contains is a helper to avoid importing strings in the test file.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
