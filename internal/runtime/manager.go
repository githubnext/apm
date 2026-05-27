package runtime

import "strings"

// SupportedRuntimeEntry describes a supported AI runtime.
type SupportedRuntimeEntry struct {
	Script      string
	Description string
	Binary      string
}

// SupportedRuntimes is the registry of known runtimes.
var SupportedRuntimes = map[string]SupportedRuntimeEntry{
	"copilot": {
		Script:      "setup-copilot",
		Description: "GitHub Copilot CLI with native MCP integration",
		Binary:      "copilot",
	},
	"codex": {
		Script:      "setup-codex",
		Description: "OpenAI Codex CLI with GitHub Models support",
		Binary:      "codex",
	},
	"llm": {
		Script:      "setup-llm",
		Description: "Simon Willison's LLM library with multiple providers",
		Binary:      "llm",
	},
	"gemini": {
		Script:      "setup-gemini",
		Description: "Google Gemini CLI with MCP integration",
		Binary:      "gemini",
	},
}

// IsKnownRuntime returns true if name is a known runtime.
func IsKnownRuntime(name string) bool {
	_, ok := SupportedRuntimes[strings.ToLower(name)]
	return ok
}

// GetSupportedRuntimeNames returns a sorted list of known runtime names.
func GetSupportedRuntimeNames() []string {
	names := make([]string, 0, len(SupportedRuntimes))
	for k := range SupportedRuntimes {
		names = append(names, k)
	}
	return names
}
