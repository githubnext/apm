// Package localcontent implements local-content integration helpers.
// Mirrors src/apm_cli/install/phases/local_content.py.
package localcontent

import (
	"os"
	"path/filepath"
)

// primitiveDirs are the recognized subdirectory names under .apm/.
var primitiveDirs = []string{
	"skills",
	"instructions",
	"chatmodes",
	"agents",
	"prompts",
	"hooks",
	"commands",
}

// ProjectHasRootPrimitives returns true when projectRoot contains a .apm/ directory.
func ProjectHasRootPrimitives(projectRoot string) bool {
	info, err := os.Stat(filepath.Join(projectRoot, ".apm"))
	return err == nil && info.IsDir()
}

// HasLocalApmContent returns true when .apm/ exists and contains at least one
// primitive file in a recognized subdirectory.
func HasLocalApmContent(projectRoot string) bool {
	apmDir := filepath.Join(projectRoot, ".apm")
	fi, err := os.Stat(apmDir)
	if err != nil || !fi.IsDir() {
		return false
	}
	for _, subdir := range primitiveDirs {
		subdirPath := filepath.Join(apmDir, subdir)
		si, err := os.Stat(subdirPath)
		if err != nil || !si.IsDir() {
			continue
		}
		// Walk for any file.
		hasFile := false
		_ = filepath.WalkDir(subdirPath, func(_ string, d os.DirEntry, err error) error {
			if err != nil || hasFile {
				return nil
			}
			if !d.IsDir() {
				hasFile = true
			}
			return nil
		})
		if hasFile {
			return true
		}
	}
	return false
}
