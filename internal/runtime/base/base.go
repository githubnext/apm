// Package base defines the RuntimeAdapter interface for LLM runtimes.
package base

// RuntimeAdapter is the base interface for LLM runtime adapters.
type RuntimeAdapter interface {
ExecutePrompt(promptContent string, args map[string]any) (string, error)
ListAvailableModels() map[string]any
GetRuntimeInfo() map[string]any
IsAvailable() bool
GetRuntimeName() string
}
