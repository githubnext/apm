package mkterrors_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mkterrors"
)

func TestMarketplaceNotFoundError_DefaultHost(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("my-market", "")
	if err.Name != "my-market" {
		t.Fatalf("Name mismatch: %q", err.Name)
	}
	if err.Host != "github.com" {
		t.Fatalf("Host should default to github.com, got %q", err.Host)
	}
	if !strings.Contains(err.Error(), "my-market") {
		t.Fatalf("error message should mention name")
	}
}

func TestMarketplaceNotFoundError_CustomHost(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("my-market", "ghe.example.com")
	if err.Host != "ghe.example.com" {
		t.Fatalf("Host mismatch: %q", err.Host)
	}
}

func TestPluginNotFoundError(t *testing.T) {
	err := mkterrors.NewPluginNotFoundError("my-plugin", "my-market")
	if err.PluginName != "my-plugin" {
		t.Fatalf("PluginName mismatch")
	}
	if err.MarketplaceName != "my-market" {
		t.Fatalf("MarketplaceName mismatch")
	}
	if !strings.Contains(err.Error(), "my-plugin") {
		t.Fatalf("error message should mention plugin name")
	}
}

func TestMarketplaceYmlError(t *testing.T) {
	err := mkterrors.NewMarketplaceYmlError("invalid yml")
	if err.Message != "invalid yml" {
		t.Fatalf("Message mismatch")
	}
	if err.Error() != "invalid yml" {
		t.Fatalf("Error() mismatch")
	}
}

func TestMarketplaceFetchError(t *testing.T) {
	err := mkterrors.NewMarketplaceFetchError("fetch failed")
	if err.Error() != "fetch failed" {
		t.Fatalf("Error() mismatch")
	}
}

func TestMarketplaceNotFoundError_MessageContainsHost(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("corp-market", "ghe.corp.io")
	msg := err.Error()
	if !strings.Contains(msg, "ghe.corp.io") {
		t.Errorf("error message should contain custom host, got %q", msg)
	}
}

func TestMarketplaceNotFoundError_MessageContainsRunInstruction(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("x", "")
	msg := err.Error()
	if !strings.Contains(msg, "apm marketplace add") {
		t.Errorf("error message should contain apm marketplace add, got %q", msg)
	}
}

func TestPluginNotFoundError_MessageContainsMarketplace(t *testing.T) {
	err := mkterrors.NewPluginNotFoundError("cool-plugin", "central")
	msg := err.Error()
	if !strings.Contains(msg, "central") {
		t.Errorf("error message should mention marketplace name, got %q", msg)
	}
}

func TestMarketplaceYmlError_EmptyMessage(t *testing.T) {
	err := mkterrors.NewMarketplaceYmlError("")
	if err.Message != "" {
		t.Errorf("expected empty Message, got %q", err.Message)
	}
	if err.Error() != "" {
		t.Errorf("expected empty Error(), got %q", err.Error())
	}
}

func TestMarketplaceFetchError_EmptyMessage(t *testing.T) {
	err := mkterrors.NewMarketplaceFetchError("")
	if err.Error() != "" {
		t.Errorf("expected empty Error(), got %q", err.Error())
	}
}

func TestMarketplaceNotFoundError_DifferentNames(t *testing.T) {
	names := []string{"alpha", "beta-market", "gamma_market"}
	for _, name := range names {
		err := mkterrors.NewMarketplaceNotFoundError(name, "")
		if err.Name != name {
			t.Errorf("Name=%q, want %q", err.Name, name)
		}
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error should mention %q, got %q", name, err.Error())
		}
	}
}

func TestPluginNotFoundError_DifferentPlugins(t *testing.T) {
	cases := []struct{ plugin, market string }{
		{"plugin-a", "mkt-1"},
		{"plugin-b", "mkt-2"},
	}
	for _, c := range cases {
		err := mkterrors.NewPluginNotFoundError(c.plugin, c.market)
		if err.PluginName != c.plugin {
			t.Errorf("PluginName=%q, want %q", err.PluginName, c.plugin)
		}
		if err.MarketplaceName != c.market {
			t.Errorf("MarketplaceName=%q, want %q", err.MarketplaceName, c.market)
		}
	}
}
