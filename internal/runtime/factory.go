package runtime

import (
	"fmt"
	"os/exec"
)

// knownAdapterNames lists the order of preference for runtime adapters.
var knownAdapterNames = []string{"copilot", "codex", "llm"}

// GetAvailableRuntimes returns info for each runtime that is installed on PATH or ~/.apm/runtimes/.
func GetAvailableRuntimes() []RuntimeInfo {
	var available []RuntimeInfo
	for _, name := range knownAdapterNames {
		entry, ok := SupportedRuntimes[name]
		if !ok {
			continue
		}
		path, _ := FindRuntimeBinary(entry.Binary)
		if path == "" {
			continue
		}
		available = append(available, RuntimeInfo{
			Name:        name,
			Available:   true,
			Description: entry.Description,
		})
	}
	return available
}

// GetRuntimeByName returns a RuntimeInfo for a named runtime, or an error if
// not found or not available.
func GetRuntimeByName(name string) (RuntimeInfo, error) {
	entry, ok := SupportedRuntimes[name]
	if !ok {
		return RuntimeInfo{}, fmt.Errorf("%w: %s", ErrRuntimeNotFound, name)
	}
	path, _ := FindRuntimeBinary(entry.Binary)
	if path == "" {
		return RuntimeInfo{}, fmt.Errorf("%w: %s", ErrRuntimeUnavailable, name)
	}
	return RuntimeInfo{
		Name:        name,
		Available:   true,
		Description: entry.Description,
	}, nil
}

// GetBestAvailableRuntime returns the highest-preference available runtime.
func GetBestAvailableRuntime() (RuntimeInfo, error) {
	for _, name := range knownAdapterNames {
		info, err := GetRuntimeByName(name)
		if err == nil {
			return info, nil
		}
	}
	return RuntimeInfo{}, fmt.Errorf("no supported runtime found (tried: copilot, codex, llm)")
}

// IsRuntimeBinaryAvailable returns true if the given runtime binary is reachable.
func IsRuntimeBinaryAvailable(binaryName string) bool {
	p, _ := FindRuntimeBinary(binaryName)
	if p != "" {
		return true
	}
	_, err := exec.LookPath(binaryName)
	return err == nil
}
