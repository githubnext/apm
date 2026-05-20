package gitauthenv

import (
	"testing"
)

func TestNew_ReturnsBuilderNotNil(t *testing.T) {
	b := New(map[string]string{"K": "V"})
	if b == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestSetupEnvironment_TerminalPromptZero(t *testing.T) {
	b := New(nil)
	env := b.SetupEnvironment()
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("expected GIT_TERMINAL_PROMPT=0, got %q", env["GIT_TERMINAL_PROMPT"])
	}
}

func TestSetupEnvironment_AskPassEcho(t *testing.T) {
	b := New(nil)
	env := b.SetupEnvironment()
	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("expected GIT_ASKPASS=echo, got %q", env["GIT_ASKPASS"])
	}
}

func TestSetupEnvironment_ConfigNoSystem(t *testing.T) {
	b := New(nil)
	env := b.SetupEnvironment()
	if env["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("expected GIT_CONFIG_NOSYSTEM=1, got %q", env["GIT_CONFIG_NOSYSTEM"])
	}
}

func TestSetupEnvironment_SSHCommandPresent(t *testing.T) {
	b := New(nil)
	env := b.SetupEnvironment()
	if env["GIT_SSH_COMMAND"] == "" {
		t.Error("expected GIT_SSH_COMMAND to be set")
	}
}

func TestSetupEnvironment_BaseKeysPreserved(t *testing.T) {
	b := New(map[string]string{"CUSTOM_KEY": "custom_val"})
	env := b.SetupEnvironment()
	if env["CUSTOM_KEY"] != "custom_val" {
		t.Errorf("expected CUSTOM_KEY=custom_val, got %q", env["CUSTOM_KEY"])
	}
}

func TestNoninteractiveEnv_TerminalPromptAlwaysZeroE3(t *testing.T) {
	env := NoninteractiveEnv(nil, NoninteractiveEnvOptions{})
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("expected 0, got %q", env["GIT_TERMINAL_PROMPT"])
	}
}

func TestNoninteractiveEnv_AskPassRemovedByDefault(t *testing.T) {
	env := NoninteractiveEnv(map[string]string{"GIT_ASKPASS": "echo"}, NoninteractiveEnvOptions{})
	if _, ok := env["GIT_ASKPASS"]; ok {
		t.Error("GIT_ASKPASS should be removed when not suppressing credentials")
	}
}

func TestNoninteractiveEnv_SuppressSetsAskPass(t *testing.T) {
	env := NoninteractiveEnv(nil, NoninteractiveEnvOptions{SuppressCredentialHelpers: true})
	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("expected GIT_ASKPASS=echo, got %q", env["GIT_ASKPASS"])
	}
}

func TestSubprocessEnvDict_ReturnsNonNil(t *testing.T) {
	env := SubprocessEnvDict(nil)
	if env == nil {
		t.Error("expected non-nil map")
	}
}

func TestSubprocessEnvDict_WithNilBase(t *testing.T) {
	env := SubprocessEnvDict(nil)
	_ = env
}

func TestSubprocessEnvDict_WithBaseKeys(t *testing.T) {
	env := SubprocessEnvDict(map[string]string{"GIT_TOKEN": "tok"})
	_ = env
}

func TestNoninteractiveEnv_PreserveConfigKeepsConfigGlobal(t *testing.T) {
	base := map[string]string{"GIT_CONFIG_GLOBAL": "/custom/path"}
	env := NoninteractiveEnv(base, NoninteractiveEnvOptions{PreserveConfigIsolation: true})
	if env["GIT_CONFIG_GLOBAL"] != "/custom/path" {
		t.Errorf("expected /custom/path, got %q", env["GIT_CONFIG_GLOBAL"])
	}
}

func TestNoninteractiveEnv_DefaultDropsConfigGlobal(t *testing.T) {
	base := map[string]string{"GIT_CONFIG_GLOBAL": "/dev/null"}
	env := NoninteractiveEnv(base, NoninteractiveEnvOptions{})
	if _, ok := env["GIT_CONFIG_GLOBAL"]; ok {
		t.Error("GIT_CONFIG_GLOBAL should be removed without PreserveConfigIsolation")
	}
}

func TestSetupEnvironment_ConfigGlobalSet(t *testing.T) {
	b := New(nil)
	env := b.SetupEnvironment()
	if env["GIT_CONFIG_GLOBAL"] == "" {
		t.Error("expected GIT_CONFIG_GLOBAL to be set")
	}
}
