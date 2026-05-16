package inittemplate_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/inittemplate"
)

func TestRenderMarketplaceYMLTemplate_Defaults(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("", "")
	if !strings.Contains(out, "my-marketplace") {
		t.Errorf("expected default name 'my-marketplace', got:\n%s", out)
	}
	if !strings.Contains(out, "acme-org") {
		t.Errorf("expected default owner 'acme-org', got:\n%s", out)
	}
}

func TestRenderMarketplaceYMLTemplate_CustomValues(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("my-mkt", "my-org")
	if !strings.Contains(out, "my-mkt") {
		t.Errorf("expected name 'my-mkt' in output")
	}
	if !strings.Contains(out, "my-org") {
		t.Errorf("expected owner 'my-org' in output")
	}
}

func TestRenderMarketplaceYMLTemplate_ContainsRequiredFields(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("mkt", "org")
	for _, field := range []string{"name:", "version:", "packages:", "description:"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in template output", field)
		}
	}
}

func TestRenderMarketplaceBlock_Defaults(t *testing.T) {
	out := inittemplate.RenderMarketplaceBlock("")
	if !strings.Contains(out, "acme-org") {
		t.Errorf("expected default owner 'acme-org', got:\n%s", out)
	}
	if !strings.Contains(out, "marketplace:") {
		t.Errorf("expected 'marketplace:' key in output")
	}
}

func TestRenderMarketplaceBlock_CustomOwner(t *testing.T) {
	out := inittemplate.RenderMarketplaceBlock("my-company")
	if !strings.Contains(out, "my-company") {
		t.Errorf("expected owner 'my-company' in block output")
	}
}
