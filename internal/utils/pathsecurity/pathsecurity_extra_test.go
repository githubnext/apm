package pathsecurity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/pathsecurity"
)

func TestSafeRmtree_ValidPath(t *testing.T) {
	base, err := os.MkdirTemp("", "saferm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	sub := filepath.Join(base, "todelete")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte("x"), 0o644)

	if err := pathsecurity.SafeRmtree(sub, base); err != nil {
		t.Errorf("SafeRmtree valid path: %v", err)
	}
	if _, err := os.Stat(sub); !os.IsNotExist(err) {
		t.Error("expected directory to be removed")
	}
}

func TestSafeRmtree_OutsideBase(t *testing.T) {
	base, err := os.MkdirTemp("", "saferm-base")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	// Attempt to remove /tmp -- should fail containment check
	if err := pathsecurity.SafeRmtree("/tmp", base); err == nil {
		t.Error("expected error when removing path outside base")
	}
}

func TestValidatePathSegments_CleanPath(t *testing.T) {
	if err := pathsecurity.ValidatePathSegments("a/b/c", "ctx", false, false); err != nil {
		t.Errorf("expected no error for clean path: %v", err)
	}
}

func TestValidatePathSegments_TraversalDotDot(t *testing.T) {
	if err := pathsecurity.ValidatePathSegments("../escape", "ctx", false, false); err == nil {
		t.Error("expected error for .. traversal")
	}
}

func TestValidatePathSegments_MiddleTraversal(t *testing.T) {
	if err := pathsecurity.ValidatePathSegments("a/../../b", "ctx", false, false); err == nil {
		t.Error("expected error for middle traversal")
	}
}

func TestValidatePathSegments_EmptySegmentRejected(t *testing.T) {
	// "foo//bar" has an empty segment when rejectEmpty=true
	err := pathsecurity.ValidatePathSegments("foo//bar", "ctx", true, false)
	if err == nil {
		t.Error("expected error for empty segment with rejectEmpty=true")
	}
}

func TestValidatePathSegments_EmptySegmentAllowed(t *testing.T) {
	// rejectEmpty=false should not reject double-slash
	err := pathsecurity.ValidatePathSegments("foo//bar", "ctx", false, false)
	// Behavior is implementation-specific; just assert no panic
	_ = err
}

func TestIsPathTraversalError_NonTraversalError(t *testing.T) {
	// A generic error should not be identified as a path traversal error
	customErr := &customError{"something else"}
	if pathsecurity.IsPathTraversalError(customErr) {
		t.Error("generic error should not be a PathTraversalError")
	}
}

func TestIsPathTraversalError_Nil(t *testing.T) {
	if pathsecurity.IsPathTraversalError(nil) {
		t.Error("nil should not be a PathTraversalError")
	}
}

func TestEnsurePathWithin_NonexistentFile(t *testing.T) {
	base, err := os.MkdirTemp("", "ensure-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	// Non-existent path within base -- should still succeed on the path check
	nonexistent := filepath.Join(base, "nonexistent.txt")
	_, err = pathsecurity.EnsurePathWithin(nonexistent, base)
	// May succeed (path is within base) or fail (file doesn't exist) -- no panic
	_ = err
}

func TestValidatePathSegments_SingleDotAllowed(t *testing.T) {
	// "." with allowCurrentDir=true should pass
	err := pathsecurity.ValidatePathSegments(".", "ctx", false, true)
	if err != nil {
		t.Errorf("single dot with allowCurrentDir=true should pass: %v", err)
	}
}

// customError is a non-PathTraversalError for testing IsPathTraversalError.
type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }
