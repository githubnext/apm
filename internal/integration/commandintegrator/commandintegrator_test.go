package commandintegrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidInputName(t *testing.T) {
	valid := []string{"foo", "Bar", "foo-bar", "foo123", "A1"}
	for _, v := range valid {
		if !isValidInputName(v) {
			t.Errorf("isValidInputName(%q) should be true", v)
		}
	}
	invalid := []string{"", "1foo", "-foo", "foo bar", "foo!bar", "a very long name that exceeds 64 chars in total length for sure yes"}
	for _, v := range invalid {
		if isValidInputName(v) {
			t.Errorf("isValidInputName(%q) should be false", v)
		}
	}
}

func TestExtractInputNamesString(t *testing.T) {
	valid, rejected := extractInputNames("myInput")
	if len(valid) != 1 || valid[0] != "myInput" {
		t.Errorf("expected [myInput], got %v", valid)
	}
	if len(rejected) != 0 {
		t.Errorf("expected no rejected, got %v", rejected)
	}
}

func TestExtractInputNamesInvalidString(t *testing.T) {
	valid, rejected := extractInputNames("123-bad")
	if len(valid) != 0 {
		t.Errorf("expected no valid, got %v", valid)
	}
	if len(rejected) != 1 {
		t.Errorf("expected 1 rejected, got %v", rejected)
	}
}

func TestExtractInputNamesSlice(t *testing.T) {
	input := []interface{}{"foo", "bar", "123bad"}
	valid, rejected := extractInputNames(input)
	if len(valid) != 2 {
		t.Errorf("expected 2 valid, got %v", valid)
	}
	if len(rejected) != 1 {
		t.Errorf("expected 1 rejected, got %v", rejected)
	}
}

func TestExtractInputNamesMapSlice(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{"name": "myArg", "description": "desc"},
		map[string]interface{}{"name": "123bad"},
	}
	valid, rejected := extractInputNames(input)
	if len(valid) != 1 || valid[0] != "myArg" {
		t.Errorf("expected [myArg], got %v", valid)
	}
	if len(rejected) != 1 {
		t.Errorf("expected 1 rejected, got %v", rejected)
	}
}

func TestExtractInputNamesNil(t *testing.T) {
	valid, rejected := extractInputNames(nil)
	if len(valid) != 0 || len(rejected) != 0 {
		t.Error("expected empty slices for nil input")
	}
}

func TestParseFrontmatter(t *testing.T) {
	content := "---\ndescription: my command\nmodel: claude-3\n---\n\nBody content here."
	meta, body := parseFrontmatter(content)
	if meta["description"] != "my command" {
		t.Errorf("description = %v", meta["description"])
	}
	if !strings.Contains(body, "Body content here.") {
		t.Errorf("body missing content: %s", body)
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	content := "Just a body without frontmatter."
	meta, body := parseFrontmatter(content)
	if len(meta) != 0 {
		t.Errorf("expected empty meta, got %v", meta)
	}
	if body != content {
		t.Errorf("body mismatch: got %q", body)
	}
}

func TestFindPromptFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.prompt.md"), []byte("content"), 0o644)
	os.WriteFile(filepath.Join(dir, "b.prompt.md"), []byte("content"), 0o644)
	os.WriteFile(filepath.Join(dir, "ignored.md"), []byte("content"), 0o644)

	files := FindPromptFiles(dir)
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(files), files)
	}
}

func TestBuildCommandContent(t *testing.T) {
	meta := map[string]interface{}{
		"description":   "Test command",
		"allowed-tools": "Bash",
	}
	body := "Do the thing."
	content := buildCommandContent(meta, body)
	if !strings.Contains(content, "description:") {
		t.Errorf("missing description in output: %s", content)
	}
	if !strings.Contains(content, "Do the thing.") {
		t.Errorf("missing body in output: %s", content)
	}
}
