// Package runtime provides runtime adapter interfaces and types for APM.
package runtime

import "errors"

// RuntimeAdapter is the base interface for LLM runtime adapters.
type RuntimeAdapter interface {
	// ExecutePrompt executes a single prompt and returns the response.
	ExecutePrompt(promptContent string, opts map[string]interface{}) (string, error)
	// ListAvailableModels lists all available models in the runtime.
	ListAvailableModels() (map[string]interface{}, error)
	// GetRuntimeInfo returns information about this runtime.
	GetRuntimeInfo() map[string]interface{}
	// IsAvailable checks if this runtime is available on the system.
	IsAvailable() bool
	// GetRuntimeName returns the name of this runtime.
	GetRuntimeName() string
}

// ErrRuntimeNotFound is returned when a requested runtime cannot be found.
var ErrRuntimeNotFound = errors.New("runtime not found")

// ErrRuntimeUnavailable is returned when a runtime is found but not available.
var ErrRuntimeUnavailable = errors.New("runtime not available")

// RuntimeInfo holds metadata about a runtime adapter.
type RuntimeInfo struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Available    bool                   `json:"available"`
	CurrentModel string                 `json:"current_model,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Error        string                 `json:"error,omitempty"`
}
