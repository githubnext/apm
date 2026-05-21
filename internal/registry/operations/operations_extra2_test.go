package operations

import (
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// extractServerIDs additional JSON shapes
// ---------------------------------------------------------------------------

func TestExtractServerIDs_NestedMCPServers(t *testing.T) {
	data := []byte(`{"mcpServers":{"server1":{"id":"s1"},"server2":{"id":"s2"}}}`)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := map[string]bool{}
	for _, id := range ids {
		found[id] = true
	}
	if !found["s1"] || !found["s2"] {
		t.Errorf("expected s1 and s2, got %v", ids)
	}
}

func TestExtractServerIDs_ServersKeyVariant(t *testing.T) {
	data := []byte(`{"servers":{"alpha":{"id":"alpha"}}}`)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "alpha" {
		t.Errorf("expected [alpha], got %v", ids)
	}
}

func TestExtractServerIDs_MixedValidAndMissingID(t *testing.T) {
	data := []byte(`{"mcpServers":{"ok":{"id":"yes"},"bad":{}}}`)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, id := range ids {
		if id == "yes" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'yes' in ids, got %v", ids)
	}
}

func TestExtractServerIDs_InvalidJSONExtra2(t *testing.T) {
	_, err := extractServerIDs([]byte("{not valid"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestExtractServerIDs_EmptyServers(t *testing.T) {
	data := []byte(`{"mcpServers":{}}`)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %v", ids)
	}
}

// ---------------------------------------------------------------------------
// mcpConfigPaths additional runtime variants
// ---------------------------------------------------------------------------

func TestMCPConfigPaths_Gemini_UserScope(t *testing.T) {
	paths := mcpConfigPaths("gemini", "/project", true)
	_ = paths // should not panic; may return empty
}

func TestMCPConfigPaths_Codex_ProjectScope(t *testing.T) {
	paths := mcpConfigPaths("codex", "/project", false)
	_ = paths
}

func TestMCPConfigPaths_EmptyRuntime(t *testing.T) {
	paths := mcpConfigPaths("", "/project", false)
	_ = paths
}

// ---------------------------------------------------------------------------
// ServerNeed struct roundtrip
// ---------------------------------------------------------------------------

func TestServerNeed_JSONRoundtrip(t *testing.T) {
	sn := ServerNeed{Reference: "my-server", NeedsInstall: true, Reason: "not installed"}
	b, err := json.Marshal(sn)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var got ServerNeed
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.Reference != sn.Reference || got.NeedsInstall != sn.NeedsInstall || got.Reason != sn.Reason {
		t.Errorf("roundtrip mismatch: got %+v", got)
	}
}

// ---------------------------------------------------------------------------
// InstallStatus struct roundtrip
// ---------------------------------------------------------------------------

func TestInstallStatus_JSONRoundtrip(t *testing.T) {
	is := InstallStatus{ServerID: "srv", Installed: true, Runtime: "copilot"}
	b, err := json.Marshal(is)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var got InstallStatus
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.ServerID != is.ServerID || got.Installed != is.Installed || got.Runtime != is.Runtime {
		t.Errorf("roundtrip mismatch: got %+v", got)
	}
}
