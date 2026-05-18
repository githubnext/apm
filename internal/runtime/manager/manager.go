// Package manager handles AI runtime installation and configuration.
// Migrated from src/apm_cli/runtime/manager.py
package manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// RuntimeInfo describes a supported AI runtime.
type RuntimeInfo struct {
	Script      string
	Description string
	Binary      string
}

// RuntimeManager manages AI runtime installation and configuration via embedded scripts.
type RuntimeManager struct {
	RuntimeDir        string
	SupportedRuntimes map[string]RuntimeInfo
}

// New creates a RuntimeManager with the default runtime directory and supported runtimes.
func New() *RuntimeManager {
	home, _ := os.UserHomeDir()
	runtimeDir := filepath.Join(home, ".apm", "runtimes")

	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".ps1"
	}

	return &RuntimeManager{
		RuntimeDir: runtimeDir,
		SupportedRuntimes: map[string]RuntimeInfo{
			"copilot": {
				Script:      "setup-copilot" + ext,
				Description: "GitHub Copilot CLI with native MCP integration",
				Binary:      "copilot",
			},
			"codex": {
				Script:      "setup-codex" + ext,
				Description: "OpenAI Codex CLI with GitHub Models support",
				Binary:      "codex",
			},
			"llm": {
				Script:      "setup-llm" + ext,
				Description: "Simon Willison's LLM library with multiple providers",
				Binary:      "llm",
			},
			"gemini": {
				Script:      "setup-gemini" + ext,
				Description: "Google Gemini CLI with MCP integration",
				Binary:      "gemini",
			},
		},
	}
}

// IsInstalled reports whether the binary for a runtime is available on PATH.
func (m *RuntimeManager) IsInstalled(name string) bool {
	info, ok := m.SupportedRuntimes[name]
	if !ok {
		return false
	}
	_, err := exec.LookPath(info.Binary)
	return err == nil
}

// GetRuntimeDir returns the directory where a specific runtime is installed.
func (m *RuntimeManager) GetRuntimeDir(name string) string {
	return filepath.Join(m.RuntimeDir, name)
}

// GetInstalledRuntimes returns the names of all installed runtimes.
func (m *RuntimeManager) GetInstalledRuntimes() []string {
	var installed []string
	for name := range m.SupportedRuntimes {
		if m.IsInstalled(name) {
			installed = append(installed, name)
		}
	}
	return installed
}

// GetScriptPath returns the path to the setup script for a runtime.
func (m *RuntimeManager) GetScriptPath(name string) (string, error) {
	info, ok := m.SupportedRuntimes[name]
	if !ok {
		return "", fmt.Errorf("unknown runtime: %s", name)
	}
	return filepath.Join(m.RuntimeDir, info.Script), nil
}

// IsWindows reports whether the current OS is Windows.
func (m *RuntimeManager) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// ValidateRuntime checks that a runtime name is supported.
func (m *RuntimeManager) ValidateRuntime(name string) error {
	if _, ok := m.SupportedRuntimes[name]; !ok {
		return fmt.Errorf("unsupported runtime: %s; supported: copilot, codex, llm, gemini", name)
	}
	return nil
}

// GetCommonScriptPath returns the path to the shared common script.
func (m *RuntimeManager) GetCommonScriptPath() string {
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".ps1"
	}
	return filepath.Join(m.RuntimeDir, "common"+ext)
}

// SetupEnvironment configures the environment variables needed for a runtime.
func (m *RuntimeManager) SetupEnvironment(name string, token string) (map[string]string, error) {
	if err := m.ValidateRuntime(name); err != nil {
		return nil, err
	}
	env := map[string]string{}
	if token != "" {
		env["GITHUB_TOKEN"] = token
		env["GH_TOKEN"] = token
	}
	return env, nil
}

// Remove uninstalls a runtime by deleting its directory.
func (m *RuntimeManager) Remove(name string) error {
	if err := m.ValidateRuntime(name); err != nil {
		return err
	}
	dir := m.GetRuntimeDir(name)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.RemoveAll(dir)
}

// ListRuntimes returns all supported runtime names and their descriptions.
func (m *RuntimeManager) ListRuntimes() map[string]string {
	result := make(map[string]string, len(m.SupportedRuntimes))
	for name, info := range m.SupportedRuntimes {
		result[name] = info.Description
	}
	return result
}
