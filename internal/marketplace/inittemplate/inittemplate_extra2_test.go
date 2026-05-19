package inittemplate_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/inittemplate"
)

func TestRenderMarketplaceYMLTemplate_EmptyInputsDefaults(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("", "")
	if !strings.Contains(out, "my-marketplace") {
		t.Error("expected default name 'my-marketplace'")
	}
	if !strings.Contains(out, "acme-org") {
		t.Error("expected default owner 'acme-org'")
	}
}

func TestRenderMarketplaceYMLTemplate_OwnerInURL(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("my-mkt", "myowner")
	if !strings.Contains(out, "github.com/myowner") {
		t.Errorf("expected github.com/myowner in output:\n%s", out)
	}
}

func TestRenderMarketplaceYMLTemplate_NonEmpty(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("x", "y")
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderMarketplaceYMLTemplate_VersionPresent(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if !strings.Contains(out, "version:") {
		t.Error("expected 'version:' field in output")
	}
}

func TestRenderMarketplaceYMLTemplate_PackagesPresent(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if !strings.Contains(out, "packages:") {
		t.Error("expected 'packages:' section in output")
	}
}

func TestRenderMarketplaceBlock_OwnerCustom(t *testing.T) {
	out := inittemplate.RenderMarketplaceBlock("custom-owner")
	if !strings.Contains(out, "custom-owner") {
		t.Errorf("expected 'custom-owner' in block output: %s", out)
	}
}

func TestRenderMarketplaceBlock_DefaultOwner(t *testing.T) {
	out := inittemplate.RenderMarketplaceBlock("")
	if len(out) == 0 {
		t.Error("expected non-empty block output for empty owner")
	}
}

func TestRenderMarketplaceYMLTemplate_BuildSectionV2(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if !strings.Contains(out, "build:") {
		t.Error("expected 'build:' section in output")
	}
}

func TestRenderMarketplaceYMLTemplate_MetadataSectionV2(t *testing.T) {
	out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
	if !strings.Contains(out, "metadata:") {
		t.Error("expected 'metadata:' section in output")
	}
}
