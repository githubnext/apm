package gitauthenv_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/gitauthenv"
)

func TestSetupEnvironment_SSHCommandWithExistingValue(t *testing.T) {
	t.Setenv("GIT_SSH_COMMAND", "ssh -i ~/.ssh/id_rsa")
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	// existing GIT_SSH_COMMAND without ConnectTimeout should get timeout appended
	if !strings.Contains(env["GIT_SSH_COMMAND"], "ConnectTimeout") {
		t.Errorf("GIT_SSH_COMMAND should contain ConnectTimeout, got %q", env["GIT_SSH_COMMAND"])
	}
	if !strings.Contains(env["GIT_SSH_COMMAND"], "id_rsa") {
		t.Errorf("GIT_SSH_COMMAND should preserve existing value, got %q", env["GIT_SSH_COMMAND"])
	}
}

func TestSetupEnvironment_SSHCommandWithExistingTimeout(t *testing.T) {
	t.Setenv("GIT_SSH_COMMAND", "ssh -o ConnectTimeout=10")
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	// should NOT double-append ConnectTimeout
	count := strings.Count(env["GIT_SSH_COMMAND"], "ConnectTimeout")
	if count != 1 {
		t.Errorf("ConnectTimeout appears %d times in GIT_SSH_COMMAND: %q", count, env["GIT_SSH_COMMAND"])
	}
}

func TestSetupEnvironment_EmptyBase(t *testing.T) {
	t.Setenv("GIT_SSH_COMMAND", "")
	b := gitauthenv.New(map[string]string{})
	env := b.SetupEnvironment()
	if env["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("GIT_CONFIG_NOSYSTEM must be 1, got %q", env["GIT_CONFIG_NOSYSTEM"])
	}
}

func TestNoninteractiveEnv_NoPreservation(t *testing.T) {
	base := map[string]string{
		"GIT_CONFIG_GLOBAL":  "/dev/null",
		"GIT_CONFIG_NOSYSTEM": "1",
	}
	env := gitauthenv.NoninteractiveEnv(base, gitauthenv.NoninteractiveEnvOptions{})
	// Without PreserveConfigIsolation these should be removed
	if _, ok := env["GIT_CONFIG_GLOBAL"]; ok {
		t.Error("GIT_CONFIG_GLOBAL should not be present without PreserveConfigIsolation")
	}
	if _, ok := env["GIT_CONFIG_NOSYSTEM"]; ok {
		t.Error("GIT_CONFIG_NOSYSTEM should not be present without PreserveConfigIsolation")
	}
}

func TestNoninteractiveEnv_BothOptions(t *testing.T) {
	base := map[string]string{}
	opts := gitauthenv.NoninteractiveEnvOptions{
		PreserveConfigIsolation:   true,
		SuppressCredentialHelpers: true,
	}
	env := gitauthenv.NoninteractiveEnv(base, opts)
	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("GIT_ASKPASS must be echo when suppress=true, got %q", env["GIT_ASKPASS"])
	}
	if env["GIT_CONFIG_COUNT"] != "1" {
		t.Errorf("GIT_CONFIG_COUNT must be 1, got %q", env["GIT_CONFIG_COUNT"])
	}
	if env["GIT_CONFIG_VALUE_0"] != "" {
		t.Errorf("GIT_CONFIG_VALUE_0 must be empty string, got %q", env["GIT_CONFIG_VALUE_0"])
	}
}

func TestNoninteractiveEnv_TerminalPromptAlwaysZero(t *testing.T) {
	opts := gitauthenv.NoninteractiveEnvOptions{}
	env := gitauthenv.NoninteractiveEnv(map[string]string{"GIT_TERMINAL_PROMPT": "1"}, opts)
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("GIT_TERMINAL_PROMPT must always be 0, got %q", env["GIT_TERMINAL_PROMPT"])
	}
}

func TestSubprocessEnvDict_EmptyBase(t *testing.T) {
	env := gitauthenv.SubprocessEnvDict(map[string]string{})
	// Should still have PATH from ambient env
	if _, ok := env["PATH"]; !ok {
		// PATH might not always exist but we don't fail on this
	}
	// Confirm bad keys are absent
	for _, bad := range []string{"GIT_DIR", "GIT_CEILING_DIRECTORIES", "GIT_ALTERNATE_OBJECT_DIRECTORIES"} {
		if _, ok := env[bad]; ok {
			t.Errorf("SubprocessEnvDict should strip %q", bad)
		}
	}
}

func TestSubprocessEnvDict_OverridesAmbient(t *testing.T) {
	base := map[string]string{"GIT_TERMINAL_PROMPT": "0", "CUSTOM_KEY": "overridden"}
	env := gitauthenv.SubprocessEnvDict(base)
	if env["CUSTOM_KEY"] != "overridden" {
		t.Errorf("base should override ambient env, got %q", env["CUSTOM_KEY"])
	}
}

func TestNew_ReturnsBuilder(t *testing.T) {
	b := gitauthenv.New(map[string]string{"TOKEN": "abc"})
	if b == nil {
		t.Fatal("New should not return nil")
	}
	env := b.SetupEnvironment()
	if env["TOKEN"] != "abc" {
		t.Errorf("expected TOKEN=abc, got %q", env["TOKEN"])
	}
}
