package policy

import "testing"

func TestPolicySource_AllFields_Extra4(t *testing.T) {
	s := PolicySource{
		Label:    "remote",
		FilePath: "/etc/policy.yml",
		Stale:    false,
		CacheAge: 120,
	}
	if s.Label != "remote" {
		t.Errorf("Label = %q", s.Label)
	}
	if s.FilePath != "/etc/policy.yml" {
		t.Errorf("FilePath = %q", s.FilePath)
	}
	if s.Stale {
		t.Error("Stale should be false")
	}
	if s.CacheAge != 120 {
		t.Errorf("CacheAge = %d", s.CacheAge)
	}
}

func TestPolicySource_ZeroValue_Extra4(t *testing.T) {
	var s PolicySource
	if s.Label != "" {
		t.Errorf("zero Label = %q", s.Label)
	}
	if s.Stale {
		t.Error("zero Stale should be false")
	}
	if s.CacheAge != 0 {
		t.Errorf("zero CacheAge = %d", s.CacheAge)
	}
}

func TestPolicyStatus_ZeroValue_Extra4(t *testing.T) {
	var st PolicyStatus
	if st.Discovered {
		t.Error("zero Discovered should be false")
	}
	if len(st.InheritanceChain) != 0 {
		t.Errorf("zero InheritanceChain len = %d", len(st.InheritanceChain))
	}
}

func TestPolicyStatus_DiscoveredTrue_Extra4(t *testing.T) {
	st := PolicyStatus{Discovered: true}
	if !st.Discovered {
		t.Error("Discovered should be true")
	}
}

func TestPolicyStatus_MultipleChain_Extra4(t *testing.T) {
	st := PolicyStatus{
		InheritanceChain: []PolicySource{
			{Label: "root"},
			{Label: "child"},
			{Label: "leaf"},
		},
	}
	if len(st.InheritanceChain) != 3 {
		t.Errorf("InheritanceChain len = %d, want 3", len(st.InheritanceChain))
	}
	if st.InheritanceChain[2].Label != "leaf" {
		t.Errorf("chain[2] = %q", st.InheritanceChain[2].Label)
	}
}

func TestStatusOptions_ZeroValue_Extra4(t *testing.T) {
	var o StatusOptions
	if o.ProjectRoot != "" {
		t.Errorf("zero ProjectRoot = %q", o.ProjectRoot)
	}
	if o.Verbose {
		t.Error("zero Verbose should be false")
	}
}

func TestStatusOptions_Fields_Extra4(t *testing.T) {
	o := StatusOptions{
		ProjectRoot: "/repo",
		Verbose:     true,
		Format:      "json",
	}
	if o.ProjectRoot != "/repo" {
		t.Errorf("ProjectRoot = %q", o.ProjectRoot)
	}
	if !o.Verbose {
		t.Error("Verbose should be true")
	}
	if o.Format != "json" {
		t.Errorf("Format = %q", o.Format)
	}
}

func TestDebugOptions_Fields_Extra4(t *testing.T) {
	o := DebugOptions{
		ProjectRoot: "/debug",
		Format:      "json",
		Source:      "local",
	}
	if o.ProjectRoot != "/debug" {
		t.Errorf("ProjectRoot = %q", o.ProjectRoot)
	}
	if o.Format != "json" {
		t.Errorf("Format = %q", o.Format)
	}
	if o.Source != "local" {
		t.Errorf("Source = %q", o.Source)
	}
}

func TestDebugOptions_ZeroValue_Extra4(t *testing.T) {
	var o DebugOptions
	if o.ProjectRoot != "" {
		t.Errorf("zero ProjectRoot = %q", o.ProjectRoot)
	}
}
