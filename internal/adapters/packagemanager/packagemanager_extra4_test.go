package packagemanager_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/adapters/packagemanager"
)

func TestDefaultManager_Name_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	if m.Name() != "default" {
		t.Fatalf("expected default, got %s", m.Name())
	}
}

func TestDefaultManager_IsSupported_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	if !m.IsSupported("anything") {
		t.Fatal("expected IsSupported to return true")
	}
}

func TestDefaultManager_InstallMissingPackage_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	dir := t.TempDir()
	err := m.Install("/nonexistent/pkg", dir)
	if err == nil {
		t.Fatal("expected error for missing package")
	}
}

func TestDefaultManager_List_EmptyDir_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	dir := t.TempDir()
	names, err := m.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected empty list, got %v", names)
	}
}

func TestDefaultManager_List_NonExistent_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	names, err := m.List("/nonexistent/dir/xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if names != nil {
		t.Fatalf("expected nil slice, got %v", names)
	}
}

func TestDefaultManager_Install_CopiesFile_Extra4(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	pkgDir := filepath.Join(src, "mypkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := packagemanager.NewDefaultManager()
	if err := m.Install(pkgDir, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultManager_Uninstall_NonExistent_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	dir := t.TempDir()
	// Uninstalling a package that doesn't exist should not error.
	if err := m.Uninstall("nothere", dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultManager_List_WithSubdirs_Extra4(t *testing.T) {
	m := packagemanager.NewDefaultManager()
	dir := t.TempDir()
	for _, sub := range []string{"alpha", "beta", "gamma"} {
		if err := os.Mkdir(filepath.Join(dir, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	names, err := m.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 entries, got %d: %v", len(names), names)
	}
}

func TestBaseManager_ErrorContainsName_Extra4(t *testing.T) {
	m := packagemanager.NewBaseManager("mypm")
	err := m.Install("x", "y")
	if err == nil || !strings.Contains(err.Error(), "mypm") {
		t.Fatalf("expected error containing 'mypm', got %v", err)
	}
}
