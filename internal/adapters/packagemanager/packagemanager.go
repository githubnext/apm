// Package packagemanager provides the base package manager abstraction and
// the default package manager implementation for APM.
//
// Corresponds to src/apm_cli/adapters/package_manager/base.py and
// src/apm_cli/adapters/package_manager/default_manager.py.
package packagemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager defines the interface that every package manager adapter must implement.
type Manager interface {
	// Name returns the human-readable name of this package manager.
	Name() string

	// Install installs a dependency at the given path.
	// packagePath is the source; installDir is the target root.
	Install(packagePath, installDir string) error

	// Uninstall removes an installed package from installDir.
	Uninstall(packageName, installDir string) error

	// List returns installed package names under installDir.
	List(installDir string) ([]string, error)

	// IsSupported returns true if this manager can handle the given package.
	IsSupported(packagePath string) bool
}

// BaseManager is an embeddable no-op implementation of Manager.
// Concrete adapters embed this and override only the methods they need.
type BaseManager struct {
	name string
}

// NewBaseManager creates a BaseManager with the given name.
func NewBaseManager(name string) *BaseManager {
	return &BaseManager{name: name}
}

func (b *BaseManager) Name() string { return b.name }

func (b *BaseManager) Install(_, _ string) error {
	return fmt.Errorf("%s: Install not implemented", b.name)
}

func (b *BaseManager) Uninstall(_, _ string) error {
	return fmt.Errorf("%s: Uninstall not implemented", b.name)
}

func (b *BaseManager) List(_ string) ([]string, error) {
	return nil, fmt.Errorf("%s: List not implemented", b.name)
}

func (b *BaseManager) IsSupported(_ string) bool { return false }

// DefaultManager is the built-in file-copy package manager for APM.
// It copies package contents into the APM modules directory using os.Rename
// where possible (same filesystem) or a full copy otherwise.
type DefaultManager struct {
	*BaseManager
}

// NewDefaultManager creates a DefaultManager.
func NewDefaultManager() *DefaultManager {
	return &DefaultManager{BaseManager: NewBaseManager("default")}
}

// IsSupported always returns true; the default manager handles every package.
func (d *DefaultManager) IsSupported(_ string) bool { return true }

// Install copies packagePath into installDir/<basename of packagePath>.
func (d *DefaultManager) Install(packagePath, installDir string) error {
	if _, err := os.Stat(packagePath); err != nil {
		return fmt.Errorf("defaultmanager: package not found: %s", packagePath)
	}
	dest := filepath.Join(installDir, filepath.Base(packagePath))
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return fmt.Errorf("defaultmanager: mkdir %s: %w", installDir, err)
	}
	// Attempt rename first (cheap on same FS), fall back to copy.
	if err := os.Rename(packagePath, dest); err != nil {
		return copyDir(packagePath, dest)
	}
	return nil
}

// Uninstall removes installDir/<packageName>.
func (d *DefaultManager) Uninstall(packageName, installDir string) error {
	target := filepath.Join(installDir, packageName)
	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("defaultmanager: uninstall %s: %w", packageName, err)
	}
	return nil
}

// List returns the names of directories under installDir.
func (d *DefaultManager) List(installDir string) ([]string, error) {
	entries, err := os.ReadDir(installDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("defaultmanager: list %s: %w", installDir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// copyDir recursively copies src to dst.
func copyDir(src, dst string) error {
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
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
