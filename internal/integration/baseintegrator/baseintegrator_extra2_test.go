package baseintegrator

import (
	"testing"
)

func TestIntegrationResult_ZeroValue(t *testing.T) {
	var r IntegrationResult
	if r.FilesIntegrated != 0 {
		t.Errorf("expected 0 FilesIntegrated, got %d", r.FilesIntegrated)
	}
	if r.FilesSkipped != 0 {
		t.Errorf("expected 0 FilesSkipped, got %d", r.FilesSkipped)
	}
	if r.SkillCreated {
		t.Error("expected SkillCreated false")
	}
}

func TestIntegrationResult_FieldAssignment(t *testing.T) {
	r := IntegrationResult{
		FilesIntegrated:  3,
		FilesSkipped:     1,
		LinksResolved:    2,
		ScriptsCopied:    4,
		SkillCreated:     true,
		SubSkillsPromoted: 5,
	}
	if r.FilesIntegrated != 3 {
		t.Errorf("unexpected FilesIntegrated %d", r.FilesIntegrated)
	}
	if !r.SkillCreated {
		t.Error("expected SkillCreated true")
	}
	if r.SubSkillsPromoted != 5 {
		t.Errorf("unexpected SubSkillsPromoted %d", r.SubSkillsPromoted)
	}
}

func TestPartitionBucketKey_WithTarget(t *testing.T) {
	key := PartitionBucketKey("myprim", "copilot")
	if key == "" {
		t.Error("expected non-empty partition key")
	}
}

func TestPartitionBucketKey_EmptyTarget(t *testing.T) {
	key := PartitionBucketKey("myprim", "")
	if key == "" {
		t.Error("expected non-empty key even with empty target")
	}
}

func TestNormalizeManagedFiles_EmptyMap(t *testing.T) {
	result := NormalizeManagedFiles(map[string]struct{}{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestNormalizeManagedFiles_NilMap(t *testing.T) {
	result := NormalizeManagedFiles(nil)
	if result != nil && len(result) != 0 {
		t.Errorf("expected nil or empty result for nil input")
	}
}

func TestNormalizeManagedFiles_ForwardSlashUnchanged(t *testing.T) {
	input := map[string]struct{}{".github/skills/foo.md": {}}
	result := NormalizeManagedFiles(input)
	if _, ok := result[".github/skills/foo.md"]; !ok {
		t.Error("expected forward-slash path to remain")
	}
}

func TestValidateDeployPath_AbsoluteRejected(t *testing.T) {
	ok := ValidateDeployPath("/absolute", "/root", []string{".github/"}, nil)
	if ok {
		t.Error("expected absolute path rejected")
	}
}

func TestValidateDeployPath_ValidPrefix(t *testing.T) {
	ok := ValidateDeployPath(".github/skills/foo.md", "/root", []string{".github/"}, nil)
	if !ok {
		t.Error("expected valid path to pass")
	}
}

func TestSyncRemoveResult_ZeroValue(t *testing.T) {
	var r SyncRemoveResult
	if r.FilesRemoved != 0 {
		t.Errorf("expected 0 FilesRemoved, got %d", r.FilesRemoved)
	}
	if r.Errors != 0 {
		t.Errorf("expected 0 Errors, got %d", r.Errors)
	}
}

func TestSyncRemoveResult_Fields(t *testing.T) {
	r := SyncRemoveResult{FilesRemoved: 2, Errors: 1}
	if r.FilesRemoved != 2 {
		t.Errorf("unexpected FilesRemoved %d", r.FilesRemoved)
	}
	if r.Errors != 1 {
		t.Errorf("unexpected Errors %d", r.Errors)
	}
}
