// Package llmruntime provides the LLM CLI runtime adapter for APM.
// Migrated from src/apm_cli/runtime/llm_runtime.py
package llmruntime

import (
	"errors"
	"os/exec"
	"strings"
)

// LLMRuntime is the APM adapter for the llm CLI tool.
type LLMRuntime struct {
	ModelName string
}

// IsAvailable returns true when the llm binary is on PATH and responds to --version.
func IsAvailable() bool {
	cmd := exec.Command("llm", "--version")
	return cmd.Run() == nil
}

// GetRuntimeName returns "llm".
func (r *LLMRuntime) GetRuntimeName() string { return "llm" }

// New creates an LLMRuntime for the given model.
// Returns an error when the llm binary is not available.
func New(modelName string) (*LLMRuntime, error) {
	if !IsAvailable() {
		return nil, errors.New("llm CLI not found. Please install: pip install llm")
	}
	return &LLMRuntime{ModelName: modelName}, nil
}

// NewDefault creates an LLMRuntime using the llm CLI default model.
func NewDefault() (*LLMRuntime, error) { return New("") }

// ExecutePrompt runs the given prompt through the llm CLI and returns the response.
func (r *LLMRuntime) ExecutePrompt(prompt string) (string, error) {
	args := []string{}
	if r.ModelName != "" {
		args = append(args, "-m", r.ModelName)
	}
	args = append(args, prompt)

	cmd := exec.Command("llm", args...)
	var buf strings.Builder
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return "", errors.New("LLM execution failed: " + buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

// ListAvailableModels returns a map of available models by querying `llm models list`.
func (r *LLMRuntime) ListAvailableModels() map[string]interface{} {
	out, err := exec.Command("llm", "models", "list").Output()
	if err != nil {
		return map[string]interface{}{"error": "failed to list models: " + err.Error()}
	}
	models := map[string]interface{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			models[line] = map[string]string{"id": line, "provider": "llm"}
		}
	}
	return models
}

// GetRuntimeInfo returns metadata about this runtime adapter.
func (r *LLMRuntime) GetRuntimeInfo() map[string]interface{} {
	model := r.ModelName
	if model == "" {
		model = "default"
	}
	return map[string]interface{}{
		"name":          "llm",
		"type":          "llm_library",
		"current_model": model,
		"capabilities": map[string]interface{}{
			"model_execution": true,
			"mcp_servers":     "runtime_dependent",
			"configuration":   "llm_commands",
			"sandboxing":      "runtime_dependent",
		},
		"description": "LLM CLI runtime adapter",
	}
}

// String returns a human-readable representation.
func (r *LLMRuntime) String() string {
	return "LLMRuntime(model=" + r.ModelName + ")"
}
