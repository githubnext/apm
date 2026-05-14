// Package factory provides a factory for creating runtime adapters with auto-detection.
package factory

import "fmt"

// RuntimeInfo holds metadata about an available runtime.
type RuntimeInfo struct {
Name      string
Available bool
Error     string
}

// RuntimeAdapter is the interface that all runtime adapters must implement.
type RuntimeAdapter interface {
GetRuntimeName() string
IsAvailable() bool
GetRuntimeInfo() RuntimeInfo
}

// ConstructableAdapter extends RuntimeAdapter with constructors.
type ConstructableAdapter interface {
RuntimeAdapter
New(modelName string) (RuntimeAdapter, error)
NewDefault() (RuntimeAdapter, error)
}

// Registry holds the ordered list of runtime adapter constructors.
type Registry struct {
adapters []ConstructableAdapter
}

// NewRegistry creates a Registry with the given adapter constructors in preference order.
func NewRegistry(adapters ...ConstructableAdapter) *Registry {
return &Registry{adapters: adapters}
}

// GetAvailableRuntimes returns metadata for all available runtimes.
func (r *Registry) GetAvailableRuntimes() []RuntimeInfo {
var out []RuntimeInfo
for _, a := range r.adapters {
if !a.IsAvailable() {
continue
}
info := a.GetRuntimeInfo()
info.Available = true
if info.Error != "" {
out = append(out, info)
continue
}
instance, err := a.NewDefault()
if err != nil {
out = append(out, RuntimeInfo{
Name:      a.GetRuntimeName(),
Available: true,
Error:     fmt.Sprintf("Available but failed to initialize: %v", err),
})
continue
}
info = instance.GetRuntimeInfo()
info.Available = true
out = append(out, info)
}
return out
}

// GetRuntimeByName returns a runtime adapter by name.
// Returns an error if the runtime is not found or not available.
func (r *Registry) GetRuntimeByName(runtimeName, modelName string) (RuntimeAdapter, error) {
for _, a := range r.adapters {
if a.GetRuntimeName() != runtimeName {
continue
}
if !a.IsAvailable() {
return nil, fmt.Errorf("runtime %q is not available on this system", runtimeName)
}
if modelName != "" {
return a.New(modelName)
}
return a.NewDefault()
}
return nil, fmt.Errorf("unknown runtime: %s", runtimeName)
}

// GetBestAvailableRuntime returns the first available runtime in preference order.
func (r *Registry) GetBestAvailableRuntime(modelName string) (RuntimeAdapter, error) {
for _, a := range r.adapters {
if !a.IsAvailable() {
continue
}
var (
instance RuntimeAdapter
err      error
)
if modelName != "" {
instance, err = a.New(modelName)
} else {
instance, err = a.NewDefault()
}
if err == nil {
return instance, nil
}
}
return nil, fmt.Errorf("no runtimes available; install at least one of: " +
"Copilot CLI (npm i -g @github/copilot), Codex CLI (npm i -g @openai/codex@native), " +
"or LLM library (pip install llm)")
}

// CreateRuntime creates a runtime adapter with optional name and model.
// If runtimeName is empty, returns the best available runtime.
func (r *Registry) CreateRuntime(runtimeName, modelName string) (RuntimeAdapter, error) {
if runtimeName != "" {
return r.GetRuntimeByName(runtimeName, modelName)
}
return r.GetBestAvailableRuntime(modelName)
}

// RuntimeExists checks if a runtime exists and is available.
func (r *Registry) RuntimeExists(runtimeName string) bool {
_, err := r.GetRuntimeByName(runtimeName, "")
return err == nil
}
