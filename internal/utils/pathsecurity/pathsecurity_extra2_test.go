package pathsecurity_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/pathsecurity"
)

func TestValidatePathSegments_Clean(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner/repo/file.go", "path", false, false)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidatePathSegments_TraversalDotDotVariant(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner/../etc/passwd", "path", false, false)
	if err == nil {
		t.Error("expected error for .. segment")
	}
	if !pathsecurity.IsPathTraversalError(err) {
		t.Errorf("expected PathTraversalError, got %T", err)
	}
}

func TestValidatePathSegments_TraversalCurrentDir(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner/./repo", "path", false, false)
	if err == nil {
		t.Error("expected error for . segment when allowCurrentDir=false")
	}
}

func TestValidatePathSegments_CurrentDirAllowed(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner/./repo", "path", false, true)
	if err != nil {
		t.Errorf("expected nil when allowCurrentDir=true, got %v", err)
	}
}

func TestValidatePathSegments_EmptySegmentReject(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner//repo", "path", true, false)
	if err == nil {
		t.Error("expected error for empty segment when rejectEmpty=true")
	}
}

func TestValidatePathSegments_EmptySegmentAllow(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner//repo", "path", false, false)
	if err != nil {
		t.Errorf("expected nil when rejectEmpty=false, got %v", err)
	}
}

func TestValidatePathSegments_PercentEncodedTraversal(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("owner/%2e%2e/repo", "path", false, false)
	if err == nil {
		t.Error("expected error for percent-encoded ..")
	}
}

func TestValidatePathSegments_EmptyPath(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("", "ctx", false, false)
	if err != nil {
		t.Errorf("expected nil for empty path, got %v", err)
	}
}

func TestIsPathTraversalError_NonTraversal(t *testing.T) {
	if pathsecurity.IsPathTraversalError(nil) {
		t.Error("nil should not be a PathTraversalError")
	}
}

func TestIsPathTraversalError_ErrorInterface(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("a/../b", "x", false, false)
	if err == nil {
		t.Fatal("expected error")
	}
	if !pathsecurity.IsPathTraversalError(err) {
		t.Errorf("expected PathTraversalError, got %T: %v", err, err)
	}
}

func TestValidatePathSegments_BackslashNormalized(t *testing.T) {
	err := pathsecurity.ValidatePathSegments(`owner\..\repo`, "path", false, false)
	if err == nil {
		t.Error("expected error for backslash-encoded traversal")
	}
}

func TestValidatePathSegments_LongCleanPath(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("a/b/c/d/e/f/g/h", "path", false, false)
	if err != nil {
		t.Errorf("expected nil for long clean path, got %v", err)
	}
}
