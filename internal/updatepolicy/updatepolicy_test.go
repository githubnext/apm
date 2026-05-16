package updatepolicy

import "testing"

func TestIsSelfUpdateEnabled_default(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = true
	if !IsSelfUpdateEnabled() {
		t.Error("expected true")
	}
	SelfUpdateEnabled = false
	if IsSelfUpdateEnabled() {
		t.Error("expected false")
	}
}

func TestGetSelfUpdateDisabledMessage_default(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = ""
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("expected default, got %q", got)
	}
}

func TestGetSelfUpdateDisabledMessage_custom(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "Use brew upgrade apm"
	got := GetSelfUpdateDisabledMessage()
	if got != "Use brew upgrade apm" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestGetSelfUpdateDisabledMessage_nonASCII(t *testing.T) {
	orig := SelfUpdateDisabledMessage
	defer func() { SelfUpdateDisabledMessage = orig }()
	SelfUpdateDisabledMessage = "Use \u2014 to update"
	got := GetSelfUpdateDisabledMessage()
	if got != DefaultSelfUpdateDisabledMessage {
		t.Errorf("expected fallback for non-ASCII, got %q", got)
	}
}

func TestGetUpdateHintMessage_enabled(t *testing.T) {
	orig := SelfUpdateEnabled
	defer func() { SelfUpdateEnabled = orig }()
	SelfUpdateEnabled = true
	got := GetUpdateHintMessage()
	if got != "Run apm update to upgrade" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestGetUpdateHintMessage_disabled(t *testing.T) {
	origEnabled := SelfUpdateEnabled
	origMsg := SelfUpdateDisabledMessage
	defer func() {
		SelfUpdateEnabled = origEnabled
		SelfUpdateDisabledMessage = origMsg
	}()
	SelfUpdateEnabled = false
	SelfUpdateDisabledMessage = "Use snap install apm"
	got := GetUpdateHintMessage()
	if got != "Use snap install apm" {
		t.Errorf("unexpected: %q", got)
	}
}
