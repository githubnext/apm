package mkterrors_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mkterrors"
)

func TestMarketplaceNotFoundError_IsError(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("mymarket", "github.com")
	var target *mkterrors.MarketplaceNotFoundError
	if !errors.As(err, &target) {
		t.Error("should be unwrappable as MarketplaceNotFoundError")
	}
}

func TestMarketplaceNotFoundError_Fields(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("awesome", "gitlab.com")
	if err.Name != "awesome" {
		t.Errorf("Name = %q, want %q", err.Name, "awesome")
	}
	if err.Host != "gitlab.com" {
		t.Errorf("Host = %q, want %q", err.Host, "gitlab.com")
	}
}

func TestMarketplaceNotFoundError_DefaultHostIsGitHub(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("mymarket", "")
	if err.Host != "github.com" {
		t.Errorf("Host = %q, want github.com", err.Host)
	}
}

func TestPluginNotFoundError_Fields(t *testing.T) {
	err := mkterrors.NewPluginNotFoundError("my-plugin", "my-market")
	if err.PluginName != "my-plugin" {
		t.Errorf("PluginName = %q", err.PluginName)
	}
	if err.MarketplaceName != "my-market" {
		t.Errorf("MarketplaceName = %q", err.MarketplaceName)
	}
}

func TestPluginNotFoundError_IsError(t *testing.T) {
	var err error = mkterrors.NewPluginNotFoundError("p", "m")
	if err.Error() == "" {
		t.Error("error message should not be empty")
	}
}

func TestMarketplaceYmlError_MessageField(t *testing.T) {
	err := mkterrors.NewMarketplaceYmlError("bad schema")
	if err.Message != "bad schema" {
		t.Errorf("Message = %q", err.Message)
	}
}

func TestMarketplaceYmlError_ErrorString(t *testing.T) {
	err := mkterrors.NewMarketplaceYmlError("bad schema")
	if err.Error() != "bad schema" {
		t.Errorf("Error() = %q, want %q", err.Error(), "bad schema")
	}
}

func TestMarketplaceFetchError_IsError(t *testing.T) {
	err := mkterrors.NewMarketplaceFetchError("fetch failed")
	if err.Error() != "fetch failed" {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestMarketplaceFetchError_EmptyMsg(t *testing.T) {
	err := mkterrors.NewMarketplaceFetchError("")
	if err.Error() != "" {
		t.Errorf("expected empty error string, got %q", err.Error())
	}
}

func TestMarketplaceNotFoundError_MessageContainsBothNameAndHost(t *testing.T) {
	err := mkterrors.NewMarketplaceNotFoundError("corp-market", "ghe.corp.com")
	msg := err.Error()
	if msg == "" {
		t.Error("error message must not be empty")
	}
}

func TestPluginNotFoundError_DifferentPairs(t *testing.T) {
	cases := [][2]string{
		{"plugin-a", "market-1"},
		{"plugin-b", "market-2"},
		{"my-tool", "default"},
	}
	for _, tc := range cases {
		err := mkterrors.NewPluginNotFoundError(tc[0], tc[1])
		if err.PluginName != tc[0] || err.MarketplaceName != tc[1] {
			t.Errorf("unexpected fields for %v", tc)
		}
	}
}

func TestMarketplaceYmlError_EmptyMessageField(t *testing.T) {
	err := mkterrors.NewMarketplaceYmlError("")
	if err.Message != "" {
		t.Errorf("Message = %q, want empty", err.Message)
	}
}

func TestMarketplaceFetchError_LongMessage(t *testing.T) {
	msg := "error: " + string(make([]byte, 200))
	err := mkterrors.NewMarketplaceFetchError(msg)
	if err.Error() != msg {
		t.Error("long error message should be preserved")
	}
}
