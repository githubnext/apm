package mcpwarnings_test

import (
"testing"

"github.com/githubnext/apm/internal/install/mcp/mcpwarnings"
)

func TestIsInternalOrMetadataHost_empty(t *testing.T) {
if mcpwarnings.IsInternalOrMetadataHost("") {
t.Error("empty host should return false")
}
}

func TestIsInternalOrMetadataHost_loopback(t *testing.T) {
if !mcpwarnings.IsInternalOrMetadataHost("127.0.0.1") {
t.Error("loopback should be internal")
}
}

func TestIsInternalOrMetadataHost_loopback6(t *testing.T) {
if !mcpwarnings.IsInternalOrMetadataHost("::1") {
t.Error("IPv6 loopback should be internal")
}
}

func TestIsInternalOrMetadataHost_RFC1918_192(t *testing.T) {
if !mcpwarnings.IsInternalOrMetadataHost("192.168.1.100") {
t.Error("192.168.x.x should be internal")
}
}

func TestIsInternalOrMetadataHost_AWS_IMDS(t *testing.T) {
if !mcpwarnings.IsInternalOrMetadataHost("169.254.169.254") {
t.Error("AWS IMDS should be a metadata host")
}
}

func TestIsInternalOrMetadataHost_publicIP_false(t *testing.T) {
if mcpwarnings.IsInternalOrMetadataHost("8.8.8.8") {
t.Error("public IP should not be internal")
}
}

func TestIsInternalOrMetadataHost_publicIP2_false(t *testing.T) {
if mcpwarnings.IsInternalOrMetadataHost("1.1.1.1") {
t.Error("1.1.1.1 should not be internal")
}
}

func TestWarnSSRFURL_empty(t *testing.T) {
w := mcpwarnings.WarnSSRFURL("")
if w != "" {
t.Errorf("expected empty warning for empty URL, got %q", w)
}
}

func TestWarnSSRFURL_safePublicURL(t *testing.T) {
w := mcpwarnings.WarnSSRFURL("https://api.example.com/v1")
if w != "" {
t.Errorf("expected no warning for public URL, got %q", w)
}
}

func TestWarnSSRFURL_localhost(t *testing.T) {
w := mcpwarnings.WarnSSRFURL("http://127.0.0.1:8080/api")
if w == "" {
t.Error("expected warning for localhost URL")
}
}

func TestWarnSSRFURL_192_range(t *testing.T) {
w := mcpwarnings.WarnSSRFURL("http://192.168.0.1/resource")
if w == "" {
t.Error("expected warning for 192.168.x.x URL")
}
}

func TestWarnShellMetachars_Semicolon(t *testing.T) {
ws := mcpwarnings.WarnShellMetachars(nil, "echo a; rm -rf /")
if len(ws) == 0 {
t.Error("expected warning for semicolon in command")
}
}

func TestWarnShellMetachars_OrOr(t *testing.T) {
ws := mcpwarnings.WarnShellMetachars(nil, "cmd1 || cmd2")
if len(ws) == 0 {
t.Error("expected warning for || in command")
}
}

func TestWarnShellMetachars_AppendRedirect(t *testing.T) {
ws := mcpwarnings.WarnShellMetachars(nil, "echo test >> /tmp/log")
if len(ws) == 0 {
t.Error("expected warning for >> in command")
}
}

func TestWarnShellMetachars_EnvDollarParen(t *testing.T) {
env := map[string]string{"MY_VAR": "$(whoami)"}
ws := mcpwarnings.WarnShellMetachars(env, "")
if len(ws) == 0 {
t.Error("expected warning for $() in env value")
}
}

func TestWarnShellMetachars_CleanEnvAndCommand(t *testing.T) {
env := map[string]string{
"HOME":  "/home/user",
"TOKEN": "abc123",
}
ws := mcpwarnings.WarnShellMetachars(env, "node server.js")
if len(ws) != 0 {
t.Errorf("expected no warnings for clean env and command, got %v", ws)
}
}

func TestWarnShellMetachars_NilEnvCleanCommand(t *testing.T) {
ws := mcpwarnings.WarnShellMetachars(nil, "python3 app.py")
if len(ws) != 0 {
t.Errorf("expected no warnings for clean command, got %v", ws)
}
}

func TestWarnShellMetachars_InputRedirect(t *testing.T) {
ws := mcpwarnings.WarnShellMetachars(nil, "wc -l < /etc/passwd")
if len(ws) == 0 {
t.Error("expected warning for < in command")
}
}
