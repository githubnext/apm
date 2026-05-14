// Package codexruntime provides the Codex CLI runtime adapter for APM.
// Migrated from src/apm_cli/runtime/codex_runtime.py
package codexruntime

import (
	"errors"
	"os/exec"
	"strings"
	"time"
)

// installCmd is the install instruction shown when codex is missing.
const installCmd = "npm i -g @openai/codex@native"

// CodexRuntime is the APM adapter for the Codex CLI.
type CodexRuntime struct {
	ModelName string
}

// IsAvailable returns true when the codex binary is on PATH.
func IsAvailable() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

// GetRuntimeName returns "codex".
func (r *CodexRuntime) GetRuntimeName() string { return "codex" }

// New creates a CodexRuntime.
// Returns an error when the codex binary is not available.
func New(modelName string) (*CodexRuntime, error) {
	if !IsAvailable() {
		return nil, errors.New("Codex CLI not available. Install with: " + installCmd)
	}
	if modelName == "" {
		modelName = "default"
	}
	return &CodexRuntime{ModelName: modelName}, nil
}

// NewDefault creates a CodexRuntime with the default model.
func NewDefault() (*CodexRuntime, error) { return New("") }

// ExecutePrompt runs the given prompt through codex exec with real-time streaming.
// Times out after 5 minutes.
func (r *CodexRuntime) ExecutePrompt(prompt string) (string, error) {
	cmd := exec.Command("codex", "exec", "--skip-git-repo-check", prompt)

	out, err := runWithTimeout(cmd, 5*time.Minute)
	if err != nil {
		if strings.Contains(out, "OPENAI_API_KEY") {
			return "", errors.New("Codex execution failed: Missing or invalid OPENAI_API_KEY. Please set your OpenAI API key.")
		}
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ListAvailableModels returns a static map of available Codex models.
// Codex does not expose model listing via CLI.
func (r *CodexRuntime) ListAvailableModels() map[string]interface{} {
	return map[string]interface{}{
		"codex-default": map[string]string{
			"id":          "codex-default",
			"provider":    "codex",
			"description": "Default Codex model (managed by Codex CLI)",
		},
	}
}

// GetRuntimeInfo returns metadata about this runtime adapter.
func (r *CodexRuntime) GetRuntimeInfo() map[string]interface{} {
	version := "unknown"
	if out, err := exec.Command("codex", "--version").Output(); err == nil {
		version = strings.TrimSpace(string(out))
	}
	return map[string]interface{}{
		"name":    "codex",
		"type":    "codex_cli",
		"version": version,
		"capabilities": map[string]interface{}{
			"model_execution": true,
			"mcp_servers":     "native_support",
			"configuration":   "config.toml",
			"sandboxing":      "built_in",
		},
		"description": "OpenAI Codex CLI runtime adapter",
	}
}

// String returns a human-readable representation.
func (r *CodexRuntime) String() string {
	return "CodexRuntime(model=" + r.ModelName + ")"
}

// runWithTimeout executes cmd, collecting all output, and returns it along with
// any error. The process is killed after timeout.
func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) (string, error) {
	var buf strings.Builder
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return "", errors.New("Codex CLI not found. Install with: " + installCmd)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		output := buf.String()
		if err != nil {
			return output, errors.New("Codex execution failed: " + err.Error())
		}
		return output, nil
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", errors.New("Codex execution timed out after 5 minutes")
	}
}
