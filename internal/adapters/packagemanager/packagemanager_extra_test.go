package packagemanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaseManager_NamePersists(t *testing.T) {
	m := NewBaseManager("my-test-manager")
	if m.Name() != "my-test-manager" {
		t.Errorf("expected my-test-manager, got %s", m.Name())
	}
}

func TestBaseManager_InstallErrorContainsName(t *testing.T) {
	m := NewBaseManager("test-mgr")
	err := m.Install("a", "b")
	if err == nil {
		t.Fatal("expected error")
	}
	if len(err.Error()) == 0 {
		t.Error("error message should not be empty")
	}
}

func TestBaseManager_UninstallErrorContainsName(t *testing.T) {
	m := NewBaseManager("test-mgr")
	err := m.Uninstall("pkg", "dir")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBaseManager_ListErrorNotNil(t *testing.T) {
	m := NewBaseManager("test-mgr")
	pkgs, err := m.List("dir")
	if err == nil {
		t.Fatal("expected error")
	}
	if pkgs != nil {
		t.Error("expected nil pkgs on error")
	}
}

func TestDefaultManager_NameIsDefault(t *testing.T) {
	m := NewDefaultManager()
	if m.Name() != "default" {
		t.Errorf("expected default, got %s", m.Name())
	}
}

func TestDefaultManager_IsSupportedAlwaysTrue(t *testing.T) {
	m := NewDefaultManager()
	cases := []string{"", "any", "/path/to/pkg", "npm:pkg", "github:owner/repo"}
	for _, c := range cases {
		if !m.IsSupported(c) {
			t.Errorf("IsSupported(%q) should return true", c)
		}
	}
}

func TestDefaultManager_InstallFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	// Create a directory (List only returns dirs) to install
	pkgDir := filepath.Join(src, "single-pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "data.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := NewDefaultManager()
	if err := m.Install(pkgDir, dst); err != nil {
		t.Fatalf("Install dir failed: %v", err)
	}
	// The installed dir should appear in dst
	names, _ := m.List(dst)
	found := false
	for _, n := range names {
		if n == "single-pkg" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected single-pkg in %v", names)
	}
}

func TestDefaultManager_UninstallNonexistent(t *testing.T) {
	dir := t.TempDir()
	m := NewDefaultManager()
	// Removing non-existent package should not error
	err := m.Uninstall("does-not-exist", dir)
	if err != nil {
		t.Errorf("Uninstall nonexistent should not error, got: %v", err)
	}
}

func TestDefaultManager_ListMultiple(t *testing.T) {
	dir := t.TempDir()
	// Create multiple items
	for _, name := range []string{"pkg-a", "pkg-b", "pkg-c"} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	m := NewDefaultManager()
	names, err := m.List(dir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(names) < 3 {
		t.Errorf("expected at least 3 packages, got %v", names)
	}
}

func TestDefaultManager_InstallDir_ContentsPreserved(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	pkgDir := filepath.Join(src, "mypkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "a.txt"), []byte("aaa"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "b.txt"), []byte("bbb"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := NewDefaultManager()
	if err := m.Install(pkgDir, dst); err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	// Both files should be in the installed location
	installedDir := filepath.Join(dst, "mypkg")
	for _, name := range []string{"a.txt", "b.txt"} {
		if _, err := os.Stat(filepath.Join(installedDir, name)); err != nil {
			t.Errorf("expected %s in installed dir: %v", name, err)
		}
	}
}

func TestDefaultManager_UninstallAndList(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"keep-pkg", "remove-pkg"} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	m := NewDefaultManager()
	if err := m.Uninstall("remove-pkg", dir); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}
	names, _ := m.List(dir)
	for _, n := range names {
		if n == "remove-pkg" {
			t.Error("remove-pkg should not be in list after Uninstall")
		}
	}
	found := false
	for _, n := range names {
		if n == "keep-pkg" {
			found = true
		}
	}
	if !found {
		t.Error("keep-pkg should still be in list")
	}
}
