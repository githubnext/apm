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

func TestRenderMarketplaceYMLTemplate_VersionField(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("pkg", "org")
if !strings.Contains(out, "version:") {
t.Error("expected 'version:' in template output")
}
if !strings.Contains(out, "0.1.0") {
t.Error("expected default version '0.1.0' in template output")
}
}

func TestRenderMarketplaceBlock_ContainsPackagesBlock(t *testing.T) {
out := inittemplate.RenderMarketplaceBlock("my-org")
if !strings.Contains(out, "packages:") {
t.Error("expected 'packages:' in block output")
}
if !strings.Contains(out, "tagPattern") {
t.Error("expected 'tagPattern' in block output")
}
}

func TestRenderMarketplaceYMLTemplate_OwnerURL(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("", "testowner")
if !strings.Contains(out, "https://github.com/testowner") {
t.Errorf("expected owner URL in template, got:\n%s", out)
}
}

func TestRenderMarketplaceBlock_NonEmpty(t *testing.T) {
out := inittemplate.RenderMarketplaceBlock("acme")
if len(out) == 0 {
t.Error("expected non-empty block output")
}
}

func TestRenderMarketplaceYMLTemplate_BuildSection(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
if !strings.Contains(out, "build:") {
t.Error("expected 'build:' section in template")
}
}
