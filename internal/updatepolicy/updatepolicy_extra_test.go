package updatepolicy

import "testing"

func TestGetSelfUpdateDisabledMessage_customASCII(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	msg := "Use brew upgrade apm to update."
	SelfUpdateDisabledMessage = msg
	got := GetSelfUpdateDisabledMessage()
	if got != msg {
		t.Errorf("expected custom message %q, got %q", msg, got)
	}
}

func TestGetUpdateHintMessage_enabledText(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = true
	got := GetUpdateHintMessage()
	if got == "" {
		t.Error("GetUpdateHintMessage returned empty string when enabled")
	}
}

func TestGetUpdateHintMessage_disabledText(t *testing.T) {
	origEnabled := SelfUpdateEnabled
	origMsg := SelfUpdateDisabledMessage
	defer func() {
		SelfUpdateEnabled = origEnabled
		SelfUpdateDisabledMessage = origMsg
	}()
	SelfUpdateEnabled = false
	SelfUpdateDisabledMessage = "Use apt upgrade apm."
	got := GetUpdateHintMessage()
	if got != "Use apt upgrade apm." {
		t.Errorf("expected custom disabled message, got %q", got)
	}
}

func TestDefaultSelfUpdateDisabledMessage_printableASCII(t *testing.T) {
	for i, c := range DefaultSelfUpdateDisabledMessage {
		if c < ' ' || c > '~' {
			t.Errorf("DefaultSelfUpdateDisabledMessage has non-printable ASCII at index %d: %U", i, c)
		}
	}
}

func TestGetSelfUpdateDisabledMessage_ControlChar(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "bad\x01message"
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("control char in message should return default, got %q", got)
	}
}

func TestGetSelfUpdateDisabledMessage_HighUnicode(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "update via \u2603"
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("unicode in message should return default, got %q", got)
	}
}

func TestIsSelfUpdateEnabled_FalseByAssignment(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = false
	if IsSelfUpdateEnabled() {
		t.Error("expected false after disabling")
	}
	SelfUpdateEnabled = true
	if !IsSelfUpdateEnabled() {
		t.Error("expected true after enabling")
	}
}

func TestGetSelfUpdateDisabledMessage_ReturnsString(t *testing.T) {
	got := GetSelfUpdateDisabledMessage()
	if got == "" {
		t.Error("GetSelfUpdateDisabledMessage must never return empty string")
	}
}
