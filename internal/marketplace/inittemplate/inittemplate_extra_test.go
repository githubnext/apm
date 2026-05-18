package inittemplate_test

import (
"strings"
"testing"

"github.com/githubnext/apm/internal/marketplace/inittemplate"
)

func TestRenderMarketplaceYMLTemplate_ContainsOwner(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("", "my-owner")
if !strings.Contains(out, "my-owner") {
t.Errorf("expected owner 'my-owner' in output:\n%s", out)
}
}

func TestRenderMarketplaceYMLTemplate_BothCustom(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("acme-mkt", "acme")
if !strings.Contains(out, "acme-mkt") {
t.Errorf("missing name 'acme-mkt'")
}
if !strings.Contains(out, "acme") {
t.Errorf("missing owner 'acme'")
}
}

func TestRenderMarketplaceYMLTemplate_IsValidYAMLLike(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("x", "y")
// Should contain colon-separated key: value pairs
if !strings.Contains(out, ": ") && !strings.Contains(out, ":\n") {
t.Error("output does not look like YAML")
}
}

func TestRenderMarketplaceYMLTemplate_NameOnlyCustom(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("my-pkg", "")
if !strings.Contains(out, "my-pkg") {
t.Errorf("expected custom name 'my-pkg' in output")
}
// Default owner should be present when empty string given.
if !strings.Contains(out, "acme-org") {
t.Errorf("expected default owner 'acme-org' when owner not provided")
}
}

func TestRenderMarketplaceBlock_IsNonEmpty(t *testing.T) {
for _, owner := range []string{"", "test-org", "github"} {
out := inittemplate.RenderMarketplaceBlock(owner)
if out == "" {
t.Errorf("RenderMarketplaceBlock(%q) returned empty string", owner)
}
}
}

func TestRenderMarketplaceBlock_ContainsMarketplaceKey(t *testing.T) {
out := inittemplate.RenderMarketplaceBlock("org")
if !strings.Contains(out, "marketplace:") {
t.Errorf("expected 'marketplace:' key in output:\n%s", out)
}
}

func TestRenderMarketplaceYMLTemplate_MetadataSection(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
if !strings.Contains(out, "metadata:") {
t.Error("expected 'metadata:' section in output")
}
}

func TestRenderMarketplaceYMLTemplate_TagPattern(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
if !strings.Contains(out, "tagPattern") {
t.Error("expected 'tagPattern' in output")
}
}

func TestRenderMarketplaceYMLTemplate_DefaultVersion(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("", "")
if !strings.Contains(out, "0.1.0") {
t.Error("expected default version '0.1.0' in output")
}
}

func TestRenderMarketplaceYMLTemplate_ExamplePackage(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
if !strings.Contains(out, "example-package") {
t.Error("expected example package stub in template output")
}
}

func TestRenderMarketplaceYMLTemplate_HasDescription(t *testing.T) {
out := inittemplate.RenderMarketplaceYMLTemplate("n", "o")
if !strings.Contains(out, "description:") {
t.Error("expected 'description:' field in template output")
}
}
