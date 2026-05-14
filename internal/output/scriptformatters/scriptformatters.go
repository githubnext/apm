// Package scriptformatters provides ASCII-only CLI output formatters for
// APM script execution.
// Migrated from src/apm_cli/output/script_formatters.py.
// Rich/colour output is omitted -- all output is plain ASCII.
package scriptformatters

import (
	"fmt"
	"strings"
)

// ScriptExecutionFormatter formats script execution output as plain ASCII lines.
type ScriptExecutionFormatter struct{}

// NewScriptExecutionFormatter returns a new formatter.
func NewScriptExecutionFormatter() *ScriptExecutionFormatter {
	return &ScriptExecutionFormatter{}
}

// FormatScriptHeader formats the script execution header with parameters.
func (f *ScriptExecutionFormatter) FormatScriptHeader(scriptName string, params map[string]string) []string {
	lines := []string{fmt.Sprintf("[>] Running script: %s", scriptName)}
	for k, v := range params {
		lines = append(lines, fmt.Sprintf("  - %s: %s", k, v))
	}
	return lines
}

// FormatCompilationProgress formats prompt compilation progress.
func (f *ScriptExecutionFormatter) FormatCompilationProgress(promptFiles []string) []string {
	if len(promptFiles) == 0 {
		return nil
	}
	var lines []string
	if len(promptFiles) == 1 {
		lines = append(lines, "Compiling prompt...")
	} else {
		lines = append(lines, fmt.Sprintf("Compiling %d prompts...", len(promptFiles)))
	}
	for _, pf := range promptFiles {
		lines = append(lines, fmt.Sprintf("|- %s", pf))
	}
	if len(lines) > 1 {
		lines[len(lines)-1] = strings.Replace(lines[len(lines)-1], "|-", "+-", 1)
	}
	return lines
}

// FormatRuntimeExecution formats runtime command execution details.
func (f *ScriptExecutionFormatter) FormatRuntimeExecution(runtime, command string, contentLength int) []string {
	return []string{
		fmt.Sprintf("Executing %s runtime...", runtime),
		fmt.Sprintf("|- Command: %s", command),
		fmt.Sprintf("+- Prompt content: %d characters", contentLength),
	}
}

// FormatContentPreview formats a content preview (plain text, no rich boxes).
func (f *ScriptExecutionFormatter) FormatContentPreview(content string, maxPreview int) []string {
	if maxPreview <= 0 {
		maxPreview = 200
	}
	preview := content
	if len(content) > maxPreview {
		preview = content[:maxPreview] + "..."
	}
	return []string{
		"Prompt preview:",
		strings.Repeat("-", 50),
		preview,
		strings.Repeat("-", 50),
	}
}

// FormatEnvironmentSetup formats environment setup information.
func (f *ScriptExecutionFormatter) FormatEnvironmentSetup(runtime string, envVarsSet []string) []string {
	if len(envVarsSet) == 0 {
		return nil
	}
	lines := []string{"Environment setup:"}
	for _, v := range envVarsSet {
		lines = append(lines, fmt.Sprintf("|- %s: configured", v))
	}
	if len(lines) > 1 {
		lines[len(lines)-1] = strings.Replace(lines[len(lines)-1], "|-", "+-", 1)
	}
	return lines
}

// FormatExecutionSuccess formats a successful execution result.
// executionTime < 0 means not provided.
func (f *ScriptExecutionFormatter) FormatExecutionSuccess(runtime string, executionTime float64) []string {
	msg := fmt.Sprintf("[+] %s execution completed successfully", titleCase(runtime))
	if executionTime >= 0 {
		msg += fmt.Sprintf(" (%.2fs)", executionTime)
	}
	return []string{msg}
}

// FormatExecutionError formats an execution error result.
func (f *ScriptExecutionFormatter) FormatExecutionError(runtime string, errorCode int, errorMsg string) []string {
	lines := []string{
		fmt.Sprintf("x %s execution failed (exit code: %d)", titleCase(runtime), errorCode),
	}
	if errorMsg != "" {
		for _, line := range strings.Split(errorMsg, "\n") {
			if strings.TrimSpace(line) != "" {
				lines = append(lines, "  "+line)
			}
		}
	}
	return lines
}

// FormatSubprocessDetails formats subprocess execution details.
func (f *ScriptExecutionFormatter) FormatSubprocessDetails(args []string, contentLength int) []string {
	quoted := make([]string, len(args))
	for i, a := range args {
		if strings.Contains(a, " ") {
			quoted[i] = `"` + a + `"`
		} else {
			quoted[i] = a
		}
	}
	return []string{
		"Subprocess execution:",
		fmt.Sprintf("|- Args: %s", strings.Join(quoted, " ")),
		fmt.Sprintf("+- Content: +%d chars appended", contentLength),
	}
}

// FormatAutoDiscoveryMessage formats the message for auto-discovered prompts.
func (f *ScriptExecutionFormatter) FormatAutoDiscoveryMessage(scriptName, promptFile, runtime string) string {
	return fmt.Sprintf("[i] Auto-discovered: %s (runtime: %s)", promptFile, runtime)
}

// titleCase capitalises the first rune of s.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
