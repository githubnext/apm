// cmdmeta.go holds per-command full descriptions and option help text.
// These mirror the Python Click CLI output for golden-file parity testing.
package main

// commandFullDesc provides the full first-paragraph description for each command,
// matching the Python CLI's Click help text exactly.
var commandFullDesc = map[string]string{
	"audit":    "Scan installed packages for hidden Unicode characters",
	"cache":    "Manage the local package cache",
	"compile":  "Compile APM context into distributed AGENTS.md files",
	"config":   "Configure APM CLI",
	"deps":     "Manage APM package dependencies",
	"experimental": "Manage experimental feature flags",
	"init":     "Initialize a new APM project",
	"install":  "Install APM and MCP dependencies (supports APM packages, Claude skills\n  (SKILL.md), and plugin collections (plugin.json); auto-creates apm.yml; use\n  --allow-insecure for http:// packages)",
	"list":     "List available scripts in the current project",
	"marketplace": "Manage marketplaces for discovery and governance",
	"mcp":      "Discover, inspect, and install MCP servers",
	"outdated": "Show outdated locked dependencies",
	"pack":     "Pack distributable artifacts from your APM project.",
	"plugin":   "Scaffold and manage plugins (plugin-author workflows)",
	"policy":   "Inspect and diagnose APM policy",
	"preview":  "Preview a script's compiled prompt files",
	"prune":    "Remove APM packages not listed in apm.yml",
	"run":      "Run a script with parameters (experimental)",
	"runtime":  "Manage AI runtimes (experimental)",
	"search":   "Search plugins in a marketplace (QUERY@MARKETPLACE)",
	"self-update": "Update the APM CLI binary itself to the latest version",
	"targets":  "Show resolved targets for the current project.",
	"uninstall": "Remove APM packages, their integrated files, and apm.yml entries",
	"unpack":   "[Deprecated] Extract an APM bundle into the current project.",
	"update":   "Refresh APM dependencies to the latest matching refs",
	"view":     "View package metadata or list remote versions.",
}

// commandOptions provides key options for each command, matching Python CLI help.
// Only the most commonly referenced options are listed; --help is added by printCmdHelp.
var commandOptions = map[string][]string{
	"compile": {
		"  -o, --output TEXT   Output file path (for single-file mode)",
		"  -t, --target TARGET Target platform (comma-separated)",
		"  --dry-run           Preview compilation without writing files",
		"  --no-links          Skip markdown link resolution",
		"  --watch             Auto-regenerate on changes",
		"  --validate          Validate primitives without compiling",
		"  --clean             Remove orphaned AGENTS.md files",
		"  --all               Compile for all canonical targets",
		"  -v, --verbose       Show detailed source attribution",
	},
	"install": {
		"  --runtime TEXT      Target specific runtime only",
		"  --exclude TEXT      Exclude specific runtime from installation",
		"  --only [apm|mcp]    Install only specific dependency type",
		"  --update            Update dependencies to latest Git references (deprecated)",
		"  --dry-run           Show what would be installed without installing",
		"  --force             Overwrite locally-authored files on collision",
		"  --frozen            Refuse to install when apm.lock.yaml is missing",
		"  -v, --verbose       Show detailed installation information",
		"  -t, --target TARGET Target harness(es) to deploy to",
		"  -g, --global        Install to user scope (~/.apm/)",
		"  --ssh               Prefer SSH transport for shorthand dependencies",
		"  --https             Prefer HTTPS transport for shorthand dependencies",
		"  --mcp NAME          Add an MCP server entry to apm.yml",
		"  --skill NAME        Install only named skill(s) from a SKILL_BUNDLE",
		"  --no-policy         Skip org policy enforcement for this invocation",
		"  --refresh           Bypass the persistent cache and re-fetch all dependencies",
		"  --dev               Install as development dependency",
		"  --allow-insecure    Allow HTTP (insecure) dependencies",
	},
	"init": {
		"  -y, --yes   Skip interactive prompts and use auto-detected defaults",
	},
	"update": {
		"  --yes          Apply updates without interactive confirmation",
		"  --dry-run      Show what would be updated without applying",
		"  -v, --verbose  Show detailed update information",
		"  -t, --target TARGET  Target harness(es) to deploy to",
	},
	"audit": {
		"  --ci       Exit non-zero if any issues are found (CI mode)",
		"  --verbose  Show detailed audit information",
	},
	"view": {
		"  --versions  List all available versions",
		"  --json      Output as JSON",
	},
	"uninstall": {
		"  --dry-run  Show what would be removed without removing",
		"  -g, --global  Uninstall from user scope (~/.apm/)",
	},
	"list": {
		"  --json     Output as JSON",
	},
}
