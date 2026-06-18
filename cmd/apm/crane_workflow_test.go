package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGoCutoverRealCraneProtectedFilesConstraints verifies that the Crane
// workflow prompt instructs the agent to strip protected workflow/config files
// from push patches when merging the base branch, and that the
// push-to-pull-request-branch safe-output configuration explicitly allows
// protected files on the crane migration branch.
//
// These properties correspond to the Python tests:
//   - test_crane_base_sync_strips_protected_workflow_files_from_push_patch
//   - test_crane_push_to_pr_branch_allows_protected_files
func TestGoCutoverRealCraneProtectedFilesConstraints(t *testing.T) {
	root := completionModuleRoot(t)
	craneWorkflow := filepath.Join(root, ".github", "workflows", "crane.md")
	data, err := os.ReadFile(craneWorkflow)
	if err != nil {
		t.Fatalf("read crane workflow: %v", err)
	}
	text := string(data)

	// Verify instructions to treat protected files as base-branch sync noise.
	if !strings.Contains(text, "trusted base-branch sync noise") {
		t.Error("crane workflow must describe protected workflow files as trusted base-branch sync noise")
	}
	if !strings.Contains(text, "git checkout ORIG_HEAD -- <path>") {
		t.Error("crane workflow must instruct restoring protected files with git checkout ORIG_HEAD -- <path>")
	}
	if !strings.Contains(text, "safe-output patch for an existing Crane PR must not include protected workflow/config files") {
		t.Error("crane workflow must warn that safe-output patch must not include protected workflow/config files")
	}

	// Verify push-to-pull-request-branch carries protected-files: allowed.
	pushIdx := strings.Index(text, "push-to-pull-request-branch:")
	if pushIdx < 0 {
		t.Fatal("crane workflow must include a push-to-pull-request-branch: configuration block")
	}
	createIssueIdx := strings.Index(text[pushIdx:], "create-issue:")
	var pushSection string
	if createIssueIdx >= 0 {
		pushSection = text[pushIdx : pushIdx+createIssueIdx]
	} else {
		pushSection = text[pushIdx:]
	}
	if !strings.Contains(pushSection, "protected-files: allowed") {
		t.Error("crane workflow push-to-pull-request-branch block must contain protected-files: allowed")
	}
}
