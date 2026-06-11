// cmd_lockfile.go provides minimal apm.lock.yaml read/write operations.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LockDep is one installed dependency entry from apm.lock.yaml.
type LockDep struct {
	Name          string
	Version       string
	RepoURL       string
	InstallPath   string
	DeployedFiles []string
}

// writeLockfile writes the canonical minimal lockfile format.
func writeLockfile(path string, deps []LockDep) error {
	var sb strings.Builder
	sb.WriteString("lockfile_version: \"1\"\n")
	sb.WriteString("dependencies:\n")
	if len(deps) == 0 {
		sb.WriteString("  []\n")
	}
	for _, d := range deps {
		sb.WriteString("  - name: " + d.Name + "\n")
		sb.WriteString("    version: " + d.Version + "\n")
		sb.WriteString("    repo_url: " + d.RepoURL + "\n")
		sb.WriteString("    install_path: " + d.InstallPath + "\n")
		if len(d.DeployedFiles) > 0 {
			sb.WriteString("    deployed_files:\n")
			for _, f := range d.DeployedFiles {
				sb.WriteString("      - " + f + "\n")
			}
		} else {
			sb.WriteString("    deployed_files: []\n")
		}
		sb.WriteString("    deployed_file_hashes: {}\n")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lockfile dir: %w", err)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// readLockfileDeps parses apm.lock.yaml and returns dependency entries.
// Returns nil, nil if the file does not exist.
func readLockfileDeps(path string) ([]LockDep, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return parseLockfileDeps(string(data)), nil
}

// parseLockfileDeps extracts LockDep entries from lockfile YAML content.
func parseLockfileDeps(content string) []LockDep {
	var deps []LockDep
	var cur *LockDep
	inDeps := false
	inDeployed := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		indent := lockfileIndent(line)

		if trimmed == "dependencies:" || trimmed == "dependencies: []" {
			inDeps = true
			continue
		}
		if indent == 0 && trimmed != "" && !strings.HasPrefix(trimmed, "-") {
			if cur != nil {
				deps = append(deps, *cur)
				cur = nil
			}
			inDeps = false
			inDeployed = false
			continue
		}

		if !inDeps {
			continue
		}

		if strings.HasPrefix(trimmed, "- name:") {
			if cur != nil {
				deps = append(deps, *cur)
			}
			cur = &LockDep{Name: strings.TrimSpace(strings.TrimPrefix(trimmed, "- name:"))}
			inDeployed = false
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "version:"):
			cur.Version = strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
		case strings.HasPrefix(trimmed, "repo_url:"):
			cur.RepoURL = strings.TrimSpace(strings.TrimPrefix(trimmed, "repo_url:"))
		case strings.HasPrefix(trimmed, "install_path:"):
			cur.InstallPath = strings.TrimSpace(strings.TrimPrefix(trimmed, "install_path:"))
		case strings.HasPrefix(trimmed, "deployed_files:"):
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "deployed_files:"))
			inDeployed = val != "[]" && val != ""
		case inDeployed && indent >= 6 && strings.HasPrefix(trimmed, "- "):
			f := strings.TrimSpace(trimmed[2:])
			if f != "" {
				cur.DeployedFiles = append(cur.DeployedFiles, f)
			}
		case !strings.HasPrefix(trimmed, "-") && !strings.HasPrefix(trimmed, "deployed_file_hashes"):
			inDeployed = false
		}
	}
	if inDeps && cur != nil {
		deps = append(deps, *cur)
	}
	return deps
}

// lockfileIndent returns the number of leading spaces in a line.
func lockfileIndent(line string) int {
	count := 0
	for _, r := range line {
		switch r {
		case ' ':
			count++
		case '\t':
			count += 4
		default:
			return count
		}
	}
	return count
}

// copyDirTree copies all files from src directory tree to dst.
func copyDirTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// walkDeployedFiles returns all file paths under dir relative to base.
func walkDeployedFiles(dir, base string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(base, path)
		if relErr != nil {
			return relErr
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	return files, err
}

// appendToApmYML appends content to apm.yml (or creates it).
func appendToApmYML(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, content)
	return err
}

// writeConfigKey writes key: value pairs to a YAML config file.
// Existing file is read and the key is updated or appended.
func writeConfigKey(path, key, value string) error {
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}
	// Simple key format: top.nested -> top:\n  nested: value
	parts := strings.SplitN(key, ".", 2)
	var newContent string
	if len(parts) == 2 {
		top, sub := parts[0], parts[1]
		topKey := top + ":"
		subEntry := "  " + sub + ": " + value
		if strings.Contains(existing, topKey) {
			lines := strings.Split(existing, "\n")
			var out []string
			inTop := false
			replaced := false
			for _, l := range lines {
				trimmed := strings.TrimSpace(l)
				if strings.HasPrefix(l, topKey) {
					inTop = true
					out = append(out, l)
					continue
				}
				if inTop {
					if strings.HasPrefix(trimmed, sub+":") {
						out = append(out, subEntry)
						replaced = true
						inTop = false
						continue
					}
					if trimmed != "" && !strings.HasPrefix(l, " ") && !strings.HasPrefix(l, "\t") {
						if !replaced {
							out = append(out, subEntry)
							replaced = true
						}
						inTop = false
					}
				}
				out = append(out, l)
			}
			if !replaced {
				out = append(out, subEntry)
			}
			newContent = strings.Join(out, "\n")
		} else {
			newContent = strings.TrimRight(existing, "\n") + "\n" + topKey + "\n" + subEntry + "\n"
		}
	} else {
		entry := key + ": " + value
		if strings.Contains(existing, key+":") {
			lines := strings.Split(existing, "\n")
			var out []string
			for _, l := range lines {
				if strings.HasPrefix(strings.TrimSpace(l), key+":") {
					out = append(out, entry)
				} else {
					out = append(out, l)
				}
			}
			newContent = strings.Join(out, "\n")
		} else {
			newContent = strings.TrimRight(existing, "\n") + "\n" + entry + "\n"
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(newContent), 0o644)
}

// readConfigKey reads a top-level key from a simple YAML config file.
// Returns the value and true if found, or empty string and false if not.
func readConfigKey(path, key string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	prefix := key + ":"
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			val := strings.TrimSpace(trimmed[len(prefix):])
			return val, true
		}
	}
	return "", false
}

// removeConfigKey removes a top-level key line from a simple YAML config file.
func removeConfigKey(path, key string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil // nothing to remove if file doesn't exist
	}
	prefix := key + ":"
	lines := strings.Split(string(data), "\n")
	var out []string
	for _, l := range lines {
		if !strings.HasPrefix(strings.TrimSpace(l), prefix) {
			out = append(out, l)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}
