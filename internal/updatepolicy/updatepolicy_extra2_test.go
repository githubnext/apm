package updatepolicy

import "testing"

func TestIsSelfUpdateEnabled_defaultTrue(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = true
	if !IsSelfUpdateEnabled() {
		t.Error("expected self-update enabled")
	}
}

func TestIsSelfUpdateEnabled_disabled(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = false
	if IsSelfUpdateEnabled() {
		t.Error("expected self-update disabled")
	}
}

func TestGetSelfUpdateDisabledMessage_defaultFallback(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = DefaultSelfUpdateDisabledMessage
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("expected default message, got %q", got)
	}
}

func TestGetSelfUpdateDisabledMessage_empty(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = ""
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("empty message should fall back to default, got %q", got)
	}
}

func TestGetSelfUpdateDisabledMessage_nonASCIIFallback(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "Update via \u2601"
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("non-ASCII message should fall back to default, got %q", got)
	}
}

func TestGetUpdateHintMessage_disabledMessage(t *testing.T) {
	origEnabled := SelfUpdateEnabled
	origMsg := SelfUpdateDisabledMessage
	defer func() {
		SelfUpdateEnabled = origEnabled
		SelfUpdateDisabledMessage = origMsg
	}()
	SelfUpdateEnabled = false
	SelfUpdateDisabledMessage = DefaultSelfUpdateDisabledMessage
	got := GetUpdateHintMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("expected disabled message, got %q", got)
	}
}

func TestGetUpdateHintMessage_enabledContainsUpdate(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = true
	got := GetUpdateHintMessage()
	if got == "" {
		t.Error("enabled hint should not be empty")
	}
}

func TestDefaultSelfUpdateDisabledMessage_nonEmpty(t *testing.T) {
	if DefaultSelfUpdateDisabledMessage == "" {
		t.Error("default disabled message should not be empty")
	}
}

func TestGetSelfUpdateDisabledMessage_asciiSpecialChars(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "Use: apm-update --force (v2.0)"
	got := GetSelfUpdateDisabledMessage()
	if got != "Use: apm-update --force (v2.0)" {
		t.Errorf("printable ASCII with special chars should be accepted, got %q", got)
	}
}
