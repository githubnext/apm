package mcpwarnings_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpwarnings"
)

func TestWarnSSRFURL_SafeURL(t *testing.T) {
	w := mcpwarnings.WarnSSRFURL("https://example.com/mcp")
	if w != "" {
		t.Errorf("expected no warning for public URL, got %q", w)
	}
}

func TestWarnSSRFURL_Empty(t *testing.T) {
	if w := mcpwarnings.WarnSSRFURL(""); w != "" {
		t.Errorf("expected empty warning for empty URL, got %q", w)
	}
}

func TestWarnSSRFURL_Loopback(t *testing.T) {
	w := mcpwarnings.WarnSSRFURL("http://127.0.0.1:8080/mcp")
	if w == "" {
		t.Error("expected warning for loopback URL")
	}
}

func TestWarnSSRFURL_MetadataHost(t *testing.T) {
	w := mcpwarnings.WarnSSRFURL("http://169.254.169.254/latest/meta-data")
	if w == "" {
		t.Error("expected warning for metadata host")
	}
}

func TestWarnSSRFURL_PrivateNetwork(t *testing.T) {
	w := mcpwarnings.WarnSSRFURL("http://192.168.1.1/mcp")
	if w == "" {
		t.Error("expected warning for RFC1918 host")
	}
}

func TestIsInternalOrMetadataHost_EmptyFalse(t *testing.T) {
	if mcpwarnings.IsInternalOrMetadataHost("") {
		t.Error("empty host should return false")
	}
}

func TestIsInternalOrMetadataHost_PublicFalse(t *testing.T) {
	if mcpwarnings.IsInternalOrMetadataHost("8.8.8.8") {
		t.Error("public IP should return false")
	}
}

func TestIsInternalOrMetadataHost_LoopbackTrue(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("127.0.0.1") {
		t.Error("loopback should return true")
	}
}

func TestIsInternalOrMetadataHost_MetadataTrue(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("169.254.169.254") {
		t.Error("metadata host should return true")
	}
}

func TestWarnShellMetachars_NoWarning(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(map[string]string{"FOO": "bar"}, "npx")
	if len(ws) != 0 {
		t.Errorf("expected no warnings, got %v", ws)
	}
}

func TestWarnShellMetachars_EnvWarning(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(map[string]string{"CMD": "$(evil)"}, "")
	if len(ws) == 0 {
		t.Error("expected warning for $( in env value")
	}
}

func TestWarnShellMetachars_CommandWarning(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(nil, "cmd; rm -rf /")
	if len(ws) == 0 {
		t.Error("expected warning for ; in command")
	}
}
