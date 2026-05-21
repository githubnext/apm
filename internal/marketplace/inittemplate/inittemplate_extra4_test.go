package inittemplate_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/inittemplate"
)

func TestRenderMarketplaceYMLTemplate_ContainsName(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("mymarket", "myorg")
	if !strings.Contains(got, "mymarket") {
		t.Errorf("expected 'mymarket' in output, got: %q", got[:100])
	}
}

func TestRenderMarketplaceYMLTemplate_NoLeadingNewline(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if len(got) > 0 && got[0] == '\n' {
		t.Error("template should not start with newline")
	}
}

func TestRenderMarketplaceYMLTemplate_DefaultOwnerPresent(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("", "")
	if !strings.Contains(got, "acme-org") {
		t.Errorf("expected default owner 'acme-org' in output")
	}
}

func TestRenderMarketplaceYMLTemplate_DefaultNamePresent(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("", "")
	if !strings.Contains(got, "my-marketplace") {
		t.Errorf("expected default name 'my-marketplace' in output")
	}
}

func TestRenderMarketplaceYMLTemplate_CustomNameOverridesDefault(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("custom-name", "some-org")
	if strings.Contains(got, "my-marketplace") {
		t.Errorf("default name 'my-marketplace' should not appear when custom name given")
	}
}

func TestRenderMarketplaceYMLTemplate_LengthReasonable(t *testing.T) {
	got := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if len(got) < 20 {
		t.Errorf("template too short: %d chars", len(got))
	}
}
