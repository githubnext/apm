package client_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client"
)

// mockClient satisfies MCPClientAdapter for interface verification.
type mockClient struct{}

func (m *mockClient) GetCurrentConfig() (map[string]interface{}, error) {
	return map[string]interface{}{"mcpServers": map[string]interface{}{}}, nil
}
func (m *mockClient) UpdateConfig(config map[string]interface{}) error      { return nil }
func (m *mockClient) ConfigureMCPServer(name, pkg string, en bool) error    { return nil }
func (m *mockClient) RemoveMCPServer(name string) error                     { return nil }
func (m *mockClient) GetTargetName() string                                  { return "mock" }

// TestParityMCPClientAdapterInterface verifies the interface type exists.
func TestParityMCPClientAdapterInterface(t *testing.T) {
	var _ client.MCPClientAdapter = (*mockClient)(nil)
}

// TestParityMCPServerEntry verifies struct fields.
func TestParityMCPServerEntry(t *testing.T) {
	entry := client.MCPServerEntry{
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-filesystem"},
		Env:     map[string]string{"TOKEN": "${GITHUB_TOKEN}"},
	}
	if entry.Command != "npx" {
		t.Fatalf("unexpected command: %s", entry.Command)
	}
	if len(entry.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(entry.Args))
	}
}

// TestParityMCPClientErrors verifies sentinel errors are defined.
func TestParityMCPClientErrors(t *testing.T) {
	if client.ErrServerNotFound == nil {
		t.Fatal("ErrServerNotFound should not be nil")
	}
	if client.ErrConfigInvalid == nil {
		t.Fatal("ErrConfigInvalid should not be nil")
	}
}

// TestParityMCPServerEntryURL verifies URL-based (SSE) server config.
func TestParityMCPServerEntryURL(t *testing.T) {
	entry := client.MCPServerEntry{
		URL:  "http://localhost:3000/sse",
		Type: "sse",
	}
	if entry.URL == "" {
		t.Fatal("expected non-empty URL")
	}
}

// TestParityMCPClientAdapterMethodGetTargetName verifies method via mock.
func TestParityMCPClientAdapterMethodGetTargetName(t *testing.T) {
	var c client.MCPClientAdapter = &mockClient{}
	if c.GetTargetName() != "mock" {
		t.Fatal("unexpected target name")
	}
}
