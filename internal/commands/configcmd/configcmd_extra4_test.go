package configcmd

import "testing"

func TestParseBoolValue_True_Extra4(t *testing.T) {
	cases := []string{"true", "1", "yes", "TRUE", "YES"}
	for _, c := range cases {
		v, err := ParseBoolValue(c)
		if err != nil {
			t.Errorf("ParseBoolValue(%q) error: %v", c, err)
		}
		if !v {
			t.Errorf("ParseBoolValue(%q) = false, want true", c)
		}
	}
}

func TestParseBoolValue_False_Extra4(t *testing.T) {
	cases := []string{"false", "0", "no", "FALSE", "NO"}
	for _, c := range cases {
		v, err := ParseBoolValue(c)
		if err != nil {
			t.Errorf("ParseBoolValue(%q) error: %v", c, err)
		}
		if v {
			t.Errorf("ParseBoolValue(%q) = true, want false", c)
		}
	}
}

func TestParseBoolValue_Invalid_Extra4(t *testing.T) {
	cases := []string{"maybe", "on", "off", "2", ""}
	for _, c := range cases {
		_, err := ParseBoolValue(c)
		if err == nil {
			t.Errorf("ParseBoolValue(%q) should return error", c)
		}
	}
}

func TestAPMConfig_ZeroValue_Extra4(t *testing.T) {
	var c APMConfig
	if c.Name != "" {
		t.Errorf("zero Name = %q", c.Name)
	}
	if c.Version != "" {
		t.Errorf("zero Version = %q", c.Version)
	}
	if c.MCPDepCount != 0 {
		t.Errorf("zero MCPDepCount = %d", c.MCPDepCount)
	}
}

func TestAPMConfig_Fields_Extra4(t *testing.T) {
	c := APMConfig{
		Name:        "my-tool",
		Version:     "1.2.3",
		Entrypoint:  "main.py",
		MCPDepCount: 5,
	}
	if c.Name != "my-tool" {
		t.Errorf("Name = %q", c.Name)
	}
	if c.Version != "1.2.3" {
		t.Errorf("Version = %q", c.Version)
	}
	if c.Entrypoint != "main.py" {
		t.Errorf("Entrypoint = %q", c.Entrypoint)
	}
	if c.MCPDepCount != 5 {
		t.Errorf("MCPDepCount = %d", c.MCPDepCount)
	}
}

func TestParseBoolValue_Whitespace_Extra4(t *testing.T) {
	v, err := ParseBoolValue("  true  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("should return true for '  true  '")
	}
}

func TestParseBoolValue_MixedCase_Extra4(t *testing.T) {
	v, err := ParseBoolValue("True")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("should return true for 'True'")
	}
}
