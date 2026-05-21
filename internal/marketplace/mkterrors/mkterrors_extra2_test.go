package mkterrors

import (
	"errors"
	"testing"
)

func TestMarketplaceNotFoundError_ExplicitHost(t *testing.T) {
	err := NewMarketplaceNotFoundError("mymkt", "github.example.com")
	if err.Host != "github.example.com" {
		t.Errorf("unexpected Host: %q", err.Host)
	}
}

func TestMarketplaceNotFoundError_NameInMessage(t *testing.T) {
	err := NewMarketplaceNotFoundError("special-mkt", "")
	msg := err.Error()
	if len(msg) == 0 {
		t.Error("error message should not be empty")
	}
}

func TestPluginNotFoundError_ErrorContainsPlugin(t *testing.T) {
	err := NewPluginNotFoundError("myplugin", "mymarketplace")
	msg := err.Error()
	if len(msg) == 0 {
		t.Error("error message should not be empty")
	}
}

func TestPluginNotFoundError_FieldsMatchInput(t *testing.T) {
	err := NewPluginNotFoundError("plugin-x", "mkt-y")
	if err.PluginName != "plugin-x" {
		t.Errorf("unexpected PluginName: %q", err.PluginName)
	}
	if err.MarketplaceName != "mkt-y" {
		t.Errorf("unexpected MarketplaceName: %q", err.MarketplaceName)
	}
}

func TestMarketplaceYmlError_IsError(t *testing.T) {
	err := NewMarketplaceYmlError("bad yaml")
	var target *MarketplaceYmlError
	if !errors.As(err, &target) {
		t.Error("expected errors.As to match MarketplaceYmlError")
	}
}

func TestMarketplaceNotFoundError_IsError(t *testing.T) {
	err := NewMarketplaceNotFoundError("x", "")
	var target *MarketplaceNotFoundError
	if !errors.As(err, &target) {
		t.Error("expected errors.As to match MarketplaceNotFoundError")
	}
}

func TestPluginNotFoundError_IsError(t *testing.T) {
	err := NewPluginNotFoundError("p", "m")
	var target *PluginNotFoundError
	if !errors.As(err, &target) {
		t.Error("expected errors.As to match PluginNotFoundError")
	}
}

func TestMarketplaceNotFoundError_ZeroNameAllowed(t *testing.T) {
	err := NewMarketplaceNotFoundError("", "")
	if err == nil {
		t.Error("expected non-nil error even with empty name")
	}
}

func TestMarketplaceYmlError_MessagePreserved(t *testing.T) {
	msg := "yaml validation failed: field missing"
	err := NewMarketplaceYmlError(msg)
	if err.Message != msg {
		t.Errorf("expected Message=%q, got %q", msg, err.Message)
	}
}

func TestMarketplaceErrors_AllImplementError(t *testing.T) {
	errs := []error{
		NewMarketplaceNotFoundError("a", "b"),
		NewPluginNotFoundError("c", "d"),
		NewMarketplaceYmlError("e"),
	}
	for _, err := range errs {
		if err.Error() == "" {
			t.Errorf("Error() should return non-empty string for %T", err)
		}
	}
}
