package mcpwarnings

import (
	"strings"
	"testing"
)

func TestIsInternalOrMetadataHost_PublicIPFalse(t *testing.T) {
	if IsInternalOrMetadataHost("8.8.8.8") {
		t.Error("public IP should return false")
	}
}

func TestIsInternalOrMetadataHost_172Block(t *testing.T) {
	if !IsInternalOrMetadataHost("172.16.0.1") {
		t.Error("172.16.x.x should be internal")
	}
}

func TestIsInternalOrMetadataHost_192Block(t *testing.T) {
	if !IsInternalOrMetadataHost("192.168.1.100") {
		t.Error("192.168.x.x should be internal")
	}
}

func TestWarnSSRFURL_SafePublicURL(t *testing.T) {
	w := WarnSSRFURL("https://api.example.com/endpoint")
	if w != "" {
		t.Errorf("public URL should have no warning, got %q", w)
	}
}

func TestWarnSSRFURL_EmptyString(t *testing.T) {
	w := WarnSSRFURL("")
	if w != "" {
		t.Errorf("empty URL should have no warning, got %q", w)
	}
}

func TestWarnShellMetachars_NilEnvClean(t *testing.T) {
	warns := WarnShellMetachars(nil, "safe-command")
	if len(warns) != 0 {
		t.Errorf("expected no warnings, got %v", warns)
	}
}

func TestWarnShellMetachars_Semicolon(t *testing.T) {
	warns := WarnShellMetachars(nil, "cmd; rm -rf /")
	if len(warns) == 0 {
		t.Error("expected shell metachar warning for semicolon")
	}
}

func TestWarnShellMetachars_DollarParen(t *testing.T) {
	warns := WarnShellMetachars(nil, "echo $(id)")
	found := false
	for _, w := range warns {
		if strings.Contains(w, "$(") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected warning about $( in command, got %v", warns)
	}
}

func TestWarnShellMetachars_EnvWithMetachar(t *testing.T) {
	warns := WarnShellMetachars(map[string]string{"VAR": "val; evil"}, "safe")
	if len(warns) == 0 {
		t.Error("expected warning for env var with semicolon")
	}
}

func TestWarnShellMetachars_PipeSeparator(t *testing.T) {
	warns := WarnShellMetachars(nil, "cmd | cat")
	if len(warns) == 0 {
		t.Error("expected warning for pipe in command")
	}
}
