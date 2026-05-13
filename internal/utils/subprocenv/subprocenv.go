// Package subprocenv provides environment sanitisation for spawning external
// processes from a PyInstaller-frozen binary.
//
// When APM ships as a PyInstaller --onedir binary the bootloader prepends
// the bundle's _internal directory to LD_LIBRARY_PATH (Linux) and the
// DYLD_* variables (macOS) so that the main Python process can find its own
// shared libraries. Child processes inherit that environment by default,
// which causes system binaries (git, curl, the install script) to resolve
// their dependencies against the bundled libraries. This package centralises
// the restoration logic that mirrors the Python subprocess_env module.
package subprocenv

import (
	"os"
	"runtime"
	"strings"
)

// pyinstallerManagedVars are the library-path variables that PyInstaller's
// bootloader rewrites at launch. Each has a sibling <NAME>_ORIG holding the
// pre-launch value that must be restored before handing the environment to a
// child process.
var pyinstallerManagedVars = []string{
	"LD_LIBRARY_PATH",    // Linux and most Unixes
	"DYLD_LIBRARY_PATH",  // macOS dynamic library search path
	"DYLD_FRAMEWORK_PATH", // macOS framework search path
}

// isFrozen returns true when the process was started by PyInstaller. This is
// detected by checking for the _MEIPASS environment variable that PyInstaller
// always sets in a frozen binary.
func isFrozen() bool {
	_, ok := os.LookupEnv("_MEIPASS")
	return ok
}

// ExternalProcessEnv returns an environment map safe for spawning external
// system binaries.
//
// When not running as a PyInstaller-frozen binary the current os.Environ() is
// returned as a fresh map with no modifications.
//
// When frozen, every library-path variable in pyinstallerManagedVars is
// restored from its <NAME>_ORIG sibling. If no _ORIG sibling exists the
// variable is removed entirely so the child does not inherit the bundle's
// _internal path. The _ORIG keys themselves are stripped.
//
// If base is non-nil it is used as the source mapping instead of os.Environ().
func ExternalProcessEnv(base map[string]string) map[string]string {
	env := envToMap(base)

	if !isFrozen() {
		return env
	}

	for _, key := range pyinstallerManagedVars {
		origKey := key + "_ORIG"
		if origVal, ok := env[origKey]; ok {
			env[key] = origVal
			delete(env, origKey)
		} else {
			delete(env, key)
		}
	}
	return env
}

// envToMap converts a []string slice (KEY=VALUE pairs) or an existing map into
// a fresh map[string]string copy. When base is nil os.Environ() is used.
func envToMap(base map[string]string) map[string]string {
	if base != nil {
		out := make(map[string]string, len(base))
		for k, v := range base {
			out[k] = v
		}
		return out
	}
	pairs := os.Environ()
	out := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			out[pair] = ""
			continue
		}
		out[pair[:idx]] = pair[idx+1:]
	}
	return out
}

// MapToSlice converts a map[string]string into a []string of KEY=VALUE pairs
// suitable for exec.Cmd.Env.
func MapToSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

// IsWindows reports whether the current OS is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
