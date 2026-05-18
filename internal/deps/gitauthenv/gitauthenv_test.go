package gitauthenv_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/gitauthenv"
)

func TestSetupEnvironment_BaseKeys(t *testing.T) {
	b := gitauthenv.New(map[string]string{"GITHUB_TOKEN": "tok123"})
	env := b.SetupEnvironment()

	for _, key := range []string{"GIT_TERMINAL_PROMPT", "GIT_ASKPASS", "GIT_CONFIG_NOSYSTEM"} {
		if _, ok := env[key]; !ok {
			t.Errorf("expected key %q in SetupEnvironment output", key)
		}
	}
	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("GIT_TERMINAL_PROMPT must be 0")
	}
	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("GIT_ASKPASS must be echo")
	}
}

func TestSetupEnvironment_PreservesBaseKeys(t *testing.T) {
	b := gitauthenv.New(map[string]string{"CUSTOM_VAR": "hello"})
	env := b.SetupEnvironment()
	if env["CUSTOM_VAR"] != "hello" {
		t.Errorf("base env key CUSTOM_VAR not preserved")
	}
}

func TestSetupEnvironment_SSHCommand(t *testing.T) {
	b := gitauthenv.New(nil)
	env := b.SetupEnvironment()
	if !strings.Contains(env["GIT_SSH_COMMAND"], "ConnectTimeout") {
		t.Errorf("GIT_SSH_COMMAND should contain ConnectTimeout, got %q", env["GIT_SSH_COMMAND"])
	}
}

func TestNoninteractiveEnv_Defaults(t *testing.T) {
	base := map[string]string{"GITHUB_TOKEN": "tok"}
	env := gitauthenv.NoninteractiveEnv(base, gitauthenv.NoninteractiveEnvOptions{})

	if env["GIT_TERMINAL_PROMPT"] != "0" {
		t.Errorf("GIT_TERMINAL_PROMPT must be 0")
	}
	if _, ok := env["GIT_ASKPASS"]; ok {
		t.Errorf("GIT_ASKPASS should be absent in default NoninteractiveEnv")
	}
}

func TestNoninteractiveEnv_SuppressCredentials(t *testing.T) {
	base := map[string]string{}
	env := gitauthenv.NoninteractiveEnv(base, gitauthenv.NoninteractiveEnvOptions{SuppressCredentialHelpers: true})

	if env["GIT_ASKPASS"] != "echo" {
		t.Errorf("GIT_ASKPASS must be echo when SuppressCredentialHelpers=true")
	}
	if env["GIT_CONFIG_KEY_0"] != "credential.helper" {
		t.Errorf("GIT_CONFIG_KEY_0 must be credential.helper")
	}
}

func TestNoninteractiveEnv_PreserveConfigIsolation(t *testing.T) {
	base := map[string]string{"GIT_CONFIG_GLOBAL": "/dev/null"}
	env := gitauthenv.NoninteractiveEnv(base, gitauthenv.NoninteractiveEnvOptions{PreserveConfigIsolation: true})
	if env["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("GIT_CONFIG_NOSYSTEM must be 1 when PreserveConfigIsolation=true")
	}
	if env["GIT_CONFIG_GLOBAL"] != "/dev/null" {
		t.Errorf("GIT_CONFIG_GLOBAL should be preserved")
	}
}

func TestSubprocessEnvDict_MergesBase(t *testing.T) {
	base := map[string]string{"MY_TOKEN": "abc"}
	env := gitauthenv.SubprocessEnvDict(base)
	if env["MY_TOKEN"] != "abc" {
		t.Errorf("base env key MY_TOKEN not present in SubprocessEnvDict output")
	}
}

func TestSubprocessEnvDict_StripsBadKeys(t *testing.T) {
	base := map[string]string{}
	env := gitauthenv.SubprocessEnvDict(base)
	for _, bad := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE"} {
		if _, ok := env[bad]; ok {
			t.Errorf("SubprocessEnvDict should strip %q but found it", bad)
		}
	}
}
