package gitauthenv_test

import (
	"testing"

	"github.com/githubnext/apm/internal/deps/gitauthenv"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	b := gitauthenv.New(nil)
	if b == nil {
		t.Error("expected non-nil builder")
	}
}

func TestSetupEnvironment_ReturnsNonNilMap(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if env == nil {
		t.Error("expected non-nil environment map")
	}
}

func TestSetupEnvironment_PreservesBaseKey(t *testing.T) {
	base := map[string]string{"CUSTOM_KEY": "custom_val"}
	b := gitauthenv.New(base)
	env := b.SetupEnvironment()
	if env["CUSTOM_KEY"] != "custom_val" {
		t.Errorf("expected CUSTOM_KEY=custom_val, got %q", env["CUSTOM_KEY"])
	}
}

func TestSetupEnvironment_SSHEnvKeyPresent(t *testing.T) {
	b := gitauthenv.New(nil)
	env := b.SetupEnvironment()
	if _, ok := env["GIT_SSH_COMMAND"]; !ok {
		t.Error("expected GIT_SSH_COMMAND in environment")
	}
}

func TestNoninteractiveEnv_ContainsTerminalPrompt(t *testing.T) {
	env := gitauthenv.NoninteractiveEnv(nil, gitauthenv.NoninteractiveEnvOptions{})
	if v, ok := env["GIT_TERMINAL_PROMPT"]; !ok || v != "0" {
		t.Errorf("expected GIT_TERMINAL_PROMPT=0, got %q", env["GIT_TERMINAL_PROMPT"])
	}
}

func TestNoninteractiveEnv_SuppressCredentialsNonNil(t *testing.T) {
	env := gitauthenv.NoninteractiveEnv(nil, gitauthenv.NoninteractiveEnvOptions{
		SuppressCredentialHelpers: true,
	})
	if env == nil {
		t.Error("expected non-nil env with SuppressCredentialHelpers=true")
	}
}

func TestNoninteractiveEnv_PreserveConfigIsolationNonNil(t *testing.T) {
	env := gitauthenv.NoninteractiveEnv(nil, gitauthenv.NoninteractiveEnvOptions{
		PreserveConfigIsolation: true,
	})
	if env == nil {
		t.Error("expected non-nil env with PreserveConfigIsolation=true")
	}
}

func TestNoninteractiveEnv_BothOptionsWithBase(t *testing.T) {
	env := gitauthenv.NoninteractiveEnv(map[string]string{"BASE": "val"}, gitauthenv.NoninteractiveEnvOptions{
		PreserveConfigIsolation:   true,
		SuppressCredentialHelpers: true,
	})
	if env["BASE"] != "val" {
		t.Errorf("expected BASE=val, got %q", env["BASE"])
	}
}

func TestSubprocessEnvDict_ReturnsMap(t *testing.T) {
	env := gitauthenv.SubprocessEnvDict(nil)
	if env == nil {
		t.Error("expected non-nil subprocess env dict")
	}
}

func TestSubprocessEnvDict_WithBase(t *testing.T) {
	base := map[string]string{"MY_KEY": "my_val"}
	env := gitauthenv.SubprocessEnvDict(base)
	if env == nil {
		t.Error("expected non-nil env")
	}
}

func TestNoninteractiveEnv_NilBase(t *testing.T) {
	env := gitauthenv.NoninteractiveEnv(nil, gitauthenv.NoninteractiveEnvOptions{})
	if env == nil {
		t.Error("expected non-nil env for nil base")
	}
}

func TestNoninteractiveEnv_TerminalPromptOverride(t *testing.T) {
	base := map[string]string{"GIT_TERMINAL_PROMPT": "1"}
	env := gitauthenv.NoninteractiveEnv(base, gitauthenv.NoninteractiveEnvOptions{})
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("GIT_TERMINAL_PROMPT should always be 0 in NoninteractiveEnv, got %q", env["GIT_TERMINAL_PROMPT"])
	}
}
