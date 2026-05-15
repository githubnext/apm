// Package runner executes APM workflow files via configured runtimes.
// It mirrors src/apm_cli/workflow/runner.py.
package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/workflow/discovery"
	"github.com/githubnext/apm/internal/workflow/wfparser"
)

// RunResult holds the outcome of a workflow execution.
type RunResult struct {
	Success  bool
	Output   string
	ErrorMsg string
}

// SubstituteParameters replaces ${input:key} placeholders in content.
func SubstituteParameters(content string, params map[string]string) string {
	result := content
	for key, value := range params {
		placeholder := fmt.Sprintf("${input:%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// CollectParameters fills in any missing parameters from defaults.
// Interactive prompting is not supported in the Go implementation;
// missing parameters are returned as empty strings.
func CollectParameters(wf *wfparser.WorkflowDefinition, provided map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range provided {
		result[k] = v
	}
	for _, param := range wf.InputParameters {
		if _, ok := result[param]; !ok {
			result[param] = "" // default empty; callers can override
		}
	}
	return result
}

// FindWorkflowByName searches for a workflow by name or file path.
// baseDir defaults to the current working directory if empty.
func FindWorkflowByName(name, baseDir string) (*wfparser.WorkflowDefinition, error) {
	if baseDir == "" {
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	// Direct file path
	if strings.HasSuffix(name, ".prompt.md") || strings.HasSuffix(name, ".workflow.md") {
		p := name
		if !filepath.IsAbs(p) {
			p = filepath.Join(baseDir, p)
		}
		if _, err := os.Stat(p); err == nil {
			return wfparser.ParseWorkflowFile(p)
		}
	}

	// Search by name
	workflows, errs := discovery.DiscoverWorkflows(baseDir)
	if len(errs) > 0 && len(workflows) == 0 {
		return nil, errs[0]
	}
	for _, wf := range workflows {
		if wf.Name == name {
			return wf, nil
		}
	}
	return nil, fmt.Errorf("workflow %q not found", name)
}

// PreviewWorkflow finds the named workflow, substitutes parameters,
// and returns the processed content without executing it.
func PreviewWorkflow(workflowName string, params map[string]string, baseDir string) RunResult {
	if params == nil {
		params = make(map[string]string)
	}
	wf, err := FindWorkflowByName(workflowName, baseDir)
	if err != nil {
		return RunResult{ErrorMsg: err.Error()}
	}
	if errs := wf.Validate(); len(errs) > 0 {
		return RunResult{ErrorMsg: fmt.Sprintf("invalid workflow: %s", strings.Join(errs, ", "))}
	}
	allParams := CollectParameters(wf, params)
	content := SubstituteParameters(wf.Content, allParams)
	return RunResult{Success: true, Output: content}
}

// RunWorkflow finds, parameterises, and dispatches a workflow to a runtime.
// The runtime lookup is deferred to the caller via RuntimeExecutor to avoid
// a hard dependency on the runtime package from this low-level module.
func RunWorkflow(
	workflowName string,
	params map[string]string,
	baseDir string,
	executor func(content, model string) (string, error),
) RunResult {
	if params == nil {
		params = make(map[string]string)
	}
	wf, err := FindWorkflowByName(workflowName, baseDir)
	if err != nil {
		return RunResult{ErrorMsg: err.Error()}
	}
	if errs := wf.Validate(); len(errs) > 0 {
		return RunResult{ErrorMsg: fmt.Sprintf("invalid workflow: %s", strings.Join(errs, ", "))}
	}
	allParams := CollectParameters(wf, params)
	content := SubstituteParameters(wf.Content, allParams)

	if executor == nil {
		return RunResult{ErrorMsg: "no runtime executor configured"}
	}
	output, err := executor(content, wf.LLMModel)
	if err != nil {
		return RunResult{ErrorMsg: fmt.Sprintf("runtime execution failed: %v", err)}
	}
	return RunResult{Success: true, Output: output}
}
