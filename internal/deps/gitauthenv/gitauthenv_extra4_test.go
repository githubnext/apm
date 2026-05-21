package gitauthenv_test

import (
	"testing"

	"github.com/githubnext/apm/internal/deps/gitauthenv"
)

func TestNew_NonNil(t *testing.T) {
	b := gitauthenv.New(map[string]string{"FOO": "bar"})
	if b == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestSetupEnvironment_TerminalPrompt(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("expected GIT_TERMINAL_PROMPT=0, got %s", env["GIT_TERMINAL_PROMPT"])
	}
}

func TestSetupEnvironment_AskPass(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("expected GIT_ASKPASS=echo, got %s", env["GIT_ASKPASS"])
	}
}

func TestSetupEnvironment_NoSystem(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if env["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("expected GIT_CONFIG_NOSYSTEM=1, got %s", env["GIT_CONFIG_NOSYSTEM"])
	}
}

func TestSetupEnvironment_SSHCommandPresent(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	ssh := env["GIT_SSH_COMMAND"]
	if ssh == "" {
		t.Error("expected non-empty GIT_SSH_COMMAND")
	}
}

func TestSetupEnvironment_PreservesBaseEnv(t *testing.T) {
	b := gitauthenv.New(map[string]string{"CUSTOM_KEY": "custom_val"})
	env := b.SetupEnvironment()
	if env["CUSTOM_KEY"] != "custom_val" {
		t.Errorf("expected CUSTOM_KEY=custom_val, got %s", env["CUSTOM_KEY"])
	}
}

func TestNoninteractiveEnv_NonNil(t *testing.T) {
	result := gitauthenv.NoninteractiveEnv(map[string]string{}, gitauthenv.NoninteractiveEnvOptions{})
	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestSubprocessEnvDict_NonNil(t *testing.T) {
	result := gitauthenv.SubprocessEnvDict(map[string]string{})
	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestSetupEnvironment_EmptyBaseEnv(t *testing.T) {
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if len(env) == 0 {
		t.Error("expected non-empty env map")
	}
}

func TestNew_NilBaseEnv(t *testing.T) {
	b := gitauthenv.New(nil)
	if b == nil {
		t.Fatal("expected non-nil builder even with nil base env")
	}
	env := b.SetupEnvironment()
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("expected GIT_TERMINAL_PROMPT=0, got %s", env["GIT_TERMINAL_PROMPT"])
	}
}
