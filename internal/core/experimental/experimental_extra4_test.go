package experimental_test

import (
"testing"

"github.com/githubnext/apm/internal/core/experimental"
)

func TestFlags_NonEmptyExtra4(t *testing.T) {
flags := experimental.Flags()
if len(flags) == 0 {
t.Error("expected at least one experimental flag")
}
}

func TestDisplayName_MapHasEntriesExtra4(t *testing.T) {
flags := experimental.Flags()
for name := range flags {
dn := experimental.DisplayName(name)
if dn == "" {
t.Errorf("expected non-empty DisplayName for %q", name)
}
break
}
}

func TestValidateFlagName_ValidFlagExtra4(t *testing.T) {
flags := experimental.Flags()
for name := range flags {
normalized, err := experimental.ValidateFlagName(name)
if err != nil {
t.Errorf("expected valid flag name %q, got error: %v", name, err)
}
if normalized == "" {
t.Errorf("expected non-empty normalized name for %q", name)
}
break
}
}

func TestValidateFlagName_InvalidFlagExtra4(t *testing.T) {
_, err := experimental.ValidateFlagName("totally-unknown-flag-xyz-abc")
if err == nil {
t.Error("expected error for unknown flag name")
}
}

func TestValidateFlagName_ConsistentResultsExtra4(t *testing.T) {
flags := experimental.Flags()
for name := range flags {
n1, _ := experimental.ValidateFlagName(name)
n2, _ := experimental.ValidateFlagName(name)
if n1 != n2 {
t.Errorf("inconsistent ValidateFlagName results: %q vs %q", n1, n2)
}
break
}
}
