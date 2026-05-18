// Package inittemplate provides template renderers for marketplace authoring scaffolds.
//
// Mirrors src/apm_cli/marketplace/init_template.py.
package inittemplate

import (
	"fmt"
	"strings"
)

// RenderMarketplaceYMLTemplate returns the scaffold content for a new marketplace.yml.
// name defaults to "my-marketplace" and owner defaults to "acme-org".
func RenderMarketplaceYMLTemplate(name, owner string) string {
	if name == "" {
		name = "my-marketplace"
	}
	if owner == "" {
		owner = "acme-org"
	}

	// The template uses {version} placeholders for literal braces in YAML.
	// In Go we use fmt.Sprintf with %% for literal braces.
	template := `# APM marketplace descriptor
#
# This file (marketplace.yml) is the SOURCE for your marketplace.
# Run 'apm pack' to compile it to marketplace.json.
# Both files must be committed to the repository.
#
# For the full schema, see:
#   https://microsoft.github.io/apm/guides/marketplace-authoring/

name: %s
description: A short description of what your marketplace offers

# Semantic version of this marketplace (bump on release)
version: 0.1.0

owner:
  name: %s
  url: https://github.com/%s
  # email: maintainers@%s.example       # optional

# APM-only build options (stripped from compiled marketplace.json)
build:
  # Default tag pattern used to resolve {version} for each package.
  # Supports {name} and {version} placeholders. Override per-package below.
  tagPattern: "v{version}"

# Opaque pass-through metadata (copied verbatim to marketplace.json).
# Use this for Anthropic-recognised or marketplace-specific fields.
metadata:
  # Example: maintained by %s
  homepage: https://example.com

packages:
  - name: example-package
    description: Human-readable description of the package
    source: %s/example-package
    version: "^1.0.0"
    # Optional overrides:
    # subdir: path/inside/repo
    # tagPattern: "example-package-v{version}"
    # include_prerelease: false
    # ref: abcdef1234  # pin to explicit SHA/tag/branch (overrides version range)

  # Alternative: pin a package to an explicit branch or SHA instead of a
  # version range.  Uncomment the entry below and remove the 'version' line.
  #
  # - name: pinned-package
  #   description: Pinned to a specific commit
  #   source: %s/pinned-package
  #   ref: main
`
	return fmt.Sprintf(template, name, owner, owner, owner, owner, owner, owner)
}

// RenderMarketplaceBlock returns a YAML snippet for the marketplace: block of apm.yml.
// Used by 'apm init --marketplace'. owner defaults to "acme-org".
func RenderMarketplaceBlock(owner string) string {
	if owner == "" {
		owner = "acme-org"
	}
	// Replace {version} placeholders with literal strings in the YAML comment.
	template := `# Marketplace authoring config (APM-only).
# Run 'apm pack' to compile this block to .claude-plugin/marketplace.json.
#
# Top-level 'name', 'description', and 'version' are inherited from
# the project (above) by default.  Override them inside this block when
# the marketplace is published independently of the project's release
# cadence.
#
# For the full schema, see:
#   https://microsoft.github.io/apm/guides/marketplace-authoring/
marketplace:
  owner:
    name: %[1]s
    url: https://github.com/%[1]s

  # Default tag pattern used to resolve version ranges for each package.
  build:
    tagPattern: "v{version}"

  packages:
    - name: example-package
      description: Human-readable description of the package
      source: %[1]s/example-package
      version: "^1.0.0"
      # Optional overrides:
      # subdir: path/inside/repo
      # tagPattern: "example-package-v{version}"
      # include_prerelease: false
      # ref: main  # pin to an explicit ref instead of a version range

    # Local-path entry: ship a package shipped alongside this repo.
    # - name: local-tool
    #   source: ./packages/local-tool
    #   description: A locally vendored tool
    #   version: 0.1.0
`
	return fmt.Sprintf(template, owner)
}

// stripBraces converts Python-style {{...}} doubled braces to single {}.
// Used when the caller passes a template string with doubled braces.
func stripBraces(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "{{", "{"), "}}", "}")
}

var _ = stripBraces // exported for potential use by callers
