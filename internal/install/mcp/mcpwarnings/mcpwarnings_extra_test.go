package mcpwarnings_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpwarnings"
)

func TestIsInternalOrMetadataHost_AlibabaCloud(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("100.100.100.200") {
		t.Error("Alibaba Cloud metadata host should return true")
	}
}

func TestIsInternalOrMetadataHost_RFC1918_10(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("10.0.0.1") {
		t.Error("10.x.x.x should be internal")
	}
}

func TestIsInternalOrMetadataHost_RFC1918_172(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("172.16.0.1") {
		t.Error("172.16.x.x should be internal")
	}
}

func TestIsInternalOrMetadataHost_LinkLocal(t *testing.T) {
	if !mcpwarnings.IsInternalOrMetadataHost("169.254.1.1") {
		t.Error("link-local address should return true")
	}
}

func TestWarnSSRFURL_MalformedURL(t *testing.T) {
	// Malformed URL should not crash; returns empty string
	w := mcpwarnings.WarnSSRFURL("://bad_url")
	_ = w // no assertion: just must not panic
}

func TestWarnSSRFURL_10Network(t *testing.T) {
	w := mcpwarnings.WarnSSRFURL("http://10.0.0.1/api")
	if w == "" {
		t.Error("expected warning for RFC1918 10.x.x.x URL")
	}
}

func TestWarnShellMetachars_BacktickCommand(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(nil, "`echo hello`")
	if len(ws) == 0 {
		t.Error("expected warning for backtick in command")
	}
}

func TestWarnShellMetachars_PipeCommand(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(nil, "cat /etc/passwd | grep root")
	if len(ws) == 0 {
		t.Error("expected warning for pipe in command")
	}
}

func TestWarnShellMetachars_MultipleEnvWarnings(t *testing.T) {
	env := map[string]string{
		"A": "$(cmd)",
		"B": "`other`",
	}
	ws := mcpwarnings.WarnShellMetachars(env, "")
	if len(ws) < 2 {
		t.Errorf("expected at least 2 warnings for 2 bad env vars, got %d", len(ws))
	}
}

func TestWarnShellMetachars_EmptyEnvAndCommand(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(map[string]string{}, "")
	if len(ws) != 0 {
		t.Errorf("expected no warnings for empty env and command, got %v", ws)
	}
}

func TestWarnShellMetachars_RedirectSymbol(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(nil, "echo hello > /tmp/out")
	if len(ws) == 0 {
		t.Error("expected warning for > in command")
	}
}

func TestWarnShellMetachars_AndAnd(t *testing.T) {
	ws := mcpwarnings.WarnShellMetachars(nil, "cmd1 && cmd2")
	if len(ws) == 0 {
		t.Error("expected warning for && in command")
	}
}
