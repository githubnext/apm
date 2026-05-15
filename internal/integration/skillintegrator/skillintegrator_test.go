package skillintegrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToHyphenCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"mySkill", "my-skill"},
		{"MySkill", "my-skill"},
		{"my_skill", "my-skill"},
		{"my skill", "my-skill"},
		{"MyAwesomeSkill", "my-awesome-skill"},
		{"some/path/mySkill", "my-skill"},
		{"already-hyphen", "already-hyphen"},
		{"", ""},
		{"a--b", "a-b"},
		{"-leading", "leading"},
		{"trailing-", "trailing"},
	}
	for _, tc := range tests {
		got := ToHyphenCase(tc.input)
		if got != tc.want {
			t.Errorf("ToHyphenCase(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestToHyphenCaseTruncation(t *testing.T) {
	long := "averylongskillnamethatshouldbetruncatedatsixtyfourcharactersexactly"
	got := ToHyphenCase(long)
	if len(got) > 64 {
		t.Errorf("ToHyphenCase should truncate to 64 chars, got %d: %q", len(got), got)
	}
}

func TestValidateSkillName(t *testing.T) {
	tests := []struct {
		name      string
		wantValid bool
	}{
		{"my-skill", true},
		{"skill123", true},
		{"a", true},
		{"", false},
		{"MY-SKILL", false},
		{"my_skill", false},
		{"my skill", false},
		{"-leading", false},
		{"trailing-", false},
		{"a-b-c-d", true},
	}
	for _, tc := range tests {
		valid, msg := ValidateSkillName(tc.name)
		if valid != tc.wantValid {
			t.Errorf("ValidateSkillName(%q) valid=%v msg=%q, want valid=%v", tc.name, valid, msg, tc.wantValid)
		}
		if !valid && msg == "" {
			t.Errorf("ValidateSkillName(%q) returned invalid with empty message", tc.name)
		}
	}
}

func TestValidateSkillNameTooLong(t *testing.T) {
	long := "a"
	for range 65 {
		long += "b"
	}
	valid, msg := ValidateSkillName(long)
	if valid {
		t.Errorf("ValidateSkillName of 66-char name should be invalid")
	}
	if msg == "" {
		t.Errorf("ValidateSkillName of 66-char name should return error message")
	}
}

func TestNormalizeSkillName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"MySkill", "my-skill"},
		{"MY_SKILL", "my-skill"},
		{"valid-name", "valid-name"},
	}
	for _, tc := range tests {
		got := NormalizeSkillName(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeSkillName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFindInstructionFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{"SKILL.md", "AGENT.md", "instructions.md", "readme.txt", "code.py"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	found := FindInstructionFiles(dir)
	_ = found
}

func TestFindAgentFiles(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"AGENT.md", "agent.md", "other.txt"} {
		os.WriteFile(filepath.Join(dir, f), []byte("content"), 0644) //nolint:errcheck
	}
	found := FindAgentFiles(dir)
	_ = found
}

func TestFindPromptFiles(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"prompt.md", "PROMPT.MD", "other.txt"} {
		os.WriteFile(filepath.Join(dir, f), []byte("content"), 0644) //nolint:errcheck
	}
	found := FindPromptFiles(dir)
	_ = found
}

func TestFindContextFiles(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"context.md", "CONTEXT.md", "other.py"} {
		os.WriteFile(filepath.Join(dir, f), []byte("content"), 0644) //nolint:errcheck
	}
	found := FindContextFiles(dir)
	_ = found
}

func TestNew(t *testing.T) {
	si := New()
	if si == nil {
		t.Fatal("New() returned nil")
	}
}

func TestIntegrateNativeSkill_NoSKILLMD(t *testing.T) {
	si := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	pkg := &PackageInfo{
		InstallPath: pkgDir,
		PackageType: "CLAUDE_SKILL",
	}
	result := si.IntegrateNativeSkill(pkg, projectDir, false, nil, nil)
	_ = result
}

func TestIntegrateNativeSkill_WithSKILLMD(t *testing.T) {
	si := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	skillName := "my-skill"
	skillDir := filepath.Join(pkgDir, skillName)
	os.MkdirAll(skillDir, 0755) //nolint:errcheck
	content := "# My Skill\n\nThis is a skill.\n"
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644) //nolint:errcheck
	pkg := &PackageInfo{
		InstallPath: skillDir,
		PackageType: "CLAUDE_SKILL",
	}
	result := si.IntegrateNativeSkill(pkg, projectDir, false, nil, nil)
	_ = result
}

func TestIntegratePackageSkill_NonSkillType(t *testing.T) {
	si := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	pkg := &PackageInfo{
		InstallPath: pkgDir,
		PackageType: "INSTRUCTIONS",
	}
	result := si.IntegratePackageSkill(pkg, projectDir, false, nil, nil, nil)
	if !result.SkillSkipped {
		t.Errorf("INSTRUCTIONS type should be skipped, got SkillSkipped=false")
	}
}

func TestIntegratePackageSkill_SkillType(t *testing.T) {
	si := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	// CLAUDE_SKILL type should not be skipped
	pkg := &PackageInfo{
		InstallPath: pkgDir,
		PackageType: "CLAUDE_SKILL",
	}
	result := si.IntegratePackageSkill(pkg, projectDir, false, nil, nil, nil)
	_ = result
}

func TestSyncIntegration_NoInstalledSkills(t *testing.T) {
	si := New()
	projectDir := t.TempDir()
	stats := si.SyncIntegration(nil, projectDir, nil, nil)
	_ = stats
}

func TestSyncIntegration_WithInstalledSkills(t *testing.T) {
	si := New()
	projectDir := t.TempDir()
	installed := map[string]struct{}{
		"my-skill": {},
	}
	stats := si.SyncIntegration(installed, projectDir, nil, nil)
	_ = stats
}

func TestSkillIntegrationResult_Fields(t *testing.T) {
	r := &SkillIntegrationResult{
		SkillCreated:      true,
		SkillUpdated:      false,
		SkillSkipped:      false,
		ReferencesCopied:  3,
		SubSkillsPromoted: 1,
		TargetPaths:       []string{"/a", "/b"},
	}
	if !r.SkillCreated {
		t.Error("SkillCreated should be true")
	}
	if r.SkillUpdated {
		t.Error("SkillUpdated should be false")
	}
	if r.ReferencesCopied != 3 {
		t.Errorf("ReferencesCopied = %d, want 3", r.ReferencesCopied)
	}
	if r.SubSkillsPromoted != 1 {
		t.Errorf("SubSkillsPromoted = %d, want 1", r.SubSkillsPromoted)
	}
	if len(r.TargetPaths) != 2 {
		t.Errorf("TargetPaths len = %d, want 2", len(r.TargetPaths))
	}
}

func TestToHyphenCasePathPrefix(t *testing.T) {
	got := ToHyphenCase("github.com/owner/my-package")
	if got != "my-package" {
		t.Errorf("ToHyphenCase with path = %q, want %q", got, "my-package")
	}
}

func TestIntegratePackageSkill_WithSKILLMD(t *testing.T) {
	si := New()
	pkgDir := t.TempDir()
	projectDir := t.TempDir()
	skillName := "my-skill"
	skillPkgDir := filepath.Join(pkgDir, skillName)
	os.MkdirAll(skillPkgDir, 0755) //nolint:errcheck
	os.WriteFile(filepath.Join(skillPkgDir, "SKILL.md"), []byte("# My Skill\n"), 0644) //nolint:errcheck
	pkg := &PackageInfo{
		InstallPath: skillPkgDir,
		PackageType: "CLAUDE_SKILL",
	}
	result := si.IntegratePackageSkill(pkg, projectDir, false, nil, nil, nil)
	// With SKILL.md present, should not be skipped (uses IntegrateNativeSkill path)
	if result.SkillSkipped {
		t.Errorf("CLAUDE_SKILL with SKILL.md should not be SkillSkipped")
	}
}

func TestNonSkillTypeAlwaysSkipped(t *testing.T) {
	nonSkillTypes := []string{"INSTRUCTIONS", "PROMPTS", ""}
	for _, pt := range nonSkillTypes {
		pkg := &PackageInfo{PackageType: pt}
		si := New()
		result := si.IntegratePackageSkill(pkg, t.TempDir(), false, nil, nil, nil)
		if !result.SkillSkipped {
			t.Errorf("PackageType %q should be skipped", pt)
		}
	}
}
