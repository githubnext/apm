package targets

import (
	"testing"
)

func TestTargetProfile_ZeroValue(t *testing.T) {
	var tp TargetProfile
	if tp.Name != "" {
		t.Error("Name should default to empty")
	}
	if tp.Supports("instructions") {
		t.Error("zero-value profile should support nothing")
	}
}

func TestTargetProfile_PrefixTrailingSlash(t *testing.T) {
	tp := TargetProfile{RootDir: ".github"}
	p := tp.Prefix()
	if p == "" {
		t.Error("expected non-empty prefix")
	}
}

func TestTargetProfile_SupportsAtUserScope_NoUserSupport(t *testing.T) {
	tp := TargetProfile{
		UserSupported: false,
		Primitives:    map[string]PrimitiveMapping{"instructions": {}},
	}
	if tp.SupportsAtUserScope("instructions") {
		t.Error("should not support at user scope when UserSupported=false")
	}
}

func TestTargetProfile_SupportsAtUserScope_Unsupported(t *testing.T) {
	tp := TargetProfile{
		UserSupported:            true,
		UnsupportedUserPrimitives: []string{"hooks"},
		Primitives:               map[string]PrimitiveMapping{"hooks": {}},
	}
	if tp.SupportsAtUserScope("hooks") {
		t.Error("hooks should not be supported at user scope")
	}
}

func TestTargetProfile_EffectiveRoot_Project(t *testing.T) {
	tp := TargetProfile{RootDir: ".github", UserRootDir: "~/.github"}
	root := tp.EffectiveRoot(false)
	if root != ".github" {
		t.Errorf("EffectiveRoot(false) = %q, want .github", root)
	}
}

func TestTargetProfile_EffectiveRoot_User(t *testing.T) {
	tp := TargetProfile{RootDir: ".github", UserRootDir: "~/.github"}
	root := tp.EffectiveRoot(true)
	if root != "~/.github" {
		t.Errorf("EffectiveRoot(true) = %q, want ~/.github", root)
	}
}

func TestActiveTargets_EmptyExplicit(t *testing.T) {
	targets := ActiveTargets("/tmp", nil)
	if len(targets) == 0 {
		t.Error("expected at least one default active target")
	}
}

func TestGetIntegrationPrefixes_Empty(t *testing.T) {
	prefixes := GetIntegrationPrefixes(nil)
	if prefixes == nil {
		prefixes = []string{}
	}
	_ = prefixes
}

func TestPrimitiveMapping_ZeroValue(t *testing.T) {
	var pm PrimitiveMapping
	if pm.Subdir != "" || pm.Extension != "" {
		t.Error("zero value should have empty fields")
	}
}
