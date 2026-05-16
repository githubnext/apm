package packagemanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaseManager_Name(t *testing.T) {
	m := NewBaseManager("mymanager")
	if m.Name() != "mymanager" {
		t.Errorf("expected 'mymanager', got %s", m.Name())
	}
}

func TestBaseManager_InstallReturnsError(t *testing.T) {
	m := NewBaseManager("base")
	err := m.Install("src", "dst")
	if err == nil {
		t.Error("expected Install to return error")
	}
}

func TestBaseManager_UninstallReturnsError(t *testing.T) {
	m := NewBaseManager("base")
	err := m.Uninstall("pkg", "dir")
	if err == nil {
		t.Error("expected Uninstall to return error")
	}
}

func TestBaseManager_ListReturnsError(t *testing.T) {
	m := NewBaseManager("base")
	_, err := m.List("dir")
	if err == nil {
		t.Error("expected List to return error")
	}
}

func TestBaseManager_IsSupportedFalse(t *testing.T) {
	m := NewBaseManager("base")
	if m.IsSupported("any") {
		t.Error("expected IsSupported to return false")
	}
}

func TestDefaultManager_Name(t *testing.T) {
	m := NewDefaultManager()
	if m.Name() != "default" {
		t.Errorf("expected 'default', got %s", m.Name())
	}
}

func TestDefaultManager_IsSupported(t *testing.T) {
	m := NewDefaultManager()
	if !m.IsSupported("anything") {
		t.Error("expected IsSupported to return true for DefaultManager")
	}
}

func TestDefaultManager_InstallMissingPackage(t *testing.T) {
	m := NewDefaultManager()
	dir := t.TempDir()
	err := m.Install("/nonexistent/path/pkg", dir)
	if err == nil {
		t.Error("expected error for missing package path")
	}
}

func TestDefaultManager_InstallAndList(t *testing.T) {
	src := t.TempDir()
	installDir := t.TempDir()

	// Create a package directory to install
	pkgDir := filepath.Join(src, "mypkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewDefaultManager()
	if err := m.Install(pkgDir, installDir); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	names, err := m.List(installDir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	found := false
	for _, n := range names {
		if n == "mypkg" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'mypkg' in list, got %v", names)
	}
}

func TestDefaultManager_Uninstall(t *testing.T) {
	installDir := t.TempDir()
	pkgDir := filepath.Join(installDir, "mypkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	m := NewDefaultManager()
	if err := m.Uninstall("mypkg", installDir); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	if _, err := os.Stat(pkgDir); !os.IsNotExist(err) {
		t.Error("expected package directory to be removed")
	}
}

func TestDefaultManager_ListEmptyDir(t *testing.T) {
	dir := t.TempDir()
	m := NewDefaultManager()
	names, err := m.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestDefaultManager_ListNonexistentDir(t *testing.T) {
	m := NewDefaultManager()
	names, err := m.List("/nonexistent/path/does/not/exist")
	if err != nil {
		t.Fatalf("List on nonexistent dir should return nil error, got %v", err)
	}
	if names != nil {
		t.Errorf("expected nil names for nonexistent dir, got %v", names)
	}
}
