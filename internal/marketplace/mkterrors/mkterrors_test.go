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
