package packagemanager_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/packagemanager"
)

func TestNewBaseManager_Name_Extra3(t *testing.T) {
	b := packagemanager.NewBaseManager("mypm")
	if b.Name() != "mypm" {
		t.Errorf("expected mypm, got %q", b.Name())
	}
}

func TestBaseManager_Install_ReturnsError_Extra3(t *testing.T) {
	b := packagemanager.NewBaseManager("base")
	err := b.Install("/src", "/dst")
	if err == nil {
		t.Error("expected error from base Install")
	}
}

func TestBaseManager_Uninstall_ReturnsError_Extra3(t *testing.T) {
	b := packagemanager.NewBaseManager("base")
	err := b.Uninstall("pkg", "/dst")
	if err == nil {
		t.Error("expected error from base Uninstall")
	}
}

func TestBaseManager_List_ReturnsError_Extra3(t *testing.T) {
	b := packagemanager.NewBaseManager("base")
	_, err := b.List("/dir")
	if err == nil {
		t.Error("expected error from base List")
	}
}

func TestBaseManager_IsSupported_False_Extra3(t *testing.T) {
	b := packagemanager.NewBaseManager("base")
	if b.IsSupported("/any/path") {
		t.Error("expected false from BaseManager.IsSupported")
	}
}

func TestNewDefaultManager_NotNil_Extra3(t *testing.T) {
	d := packagemanager.NewDefaultManager()
	if d == nil {
		t.Error("expected non-nil")
	}
}

func TestDefaultManager_Name_Extra3(t *testing.T) {
	d := packagemanager.NewDefaultManager()
	if d.Name() == "" {
		t.Error("expected non-empty name")
	}
}

func TestDefaultManager_IsSupported_Extra3(t *testing.T) {
	d := packagemanager.NewDefaultManager()
	if !d.IsSupported("/any/path") {
		t.Error("DefaultManager.IsSupported should return true")
	}
}

func TestDefaultManager_List_EmptyDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := packagemanager.NewDefaultManager()
	pkgs, err := d.List(dir)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(pkgs) != 0 {
		t.Errorf("expected empty list, got %v", pkgs)
	}
}

func TestDefaultManager_Install_ThenList_Extra3(t *testing.T) {
	srcDir := t.TempDir()
	pkgDir := filepath.Join(srcDir, "mypkg")
	_ = os.MkdirAll(pkgDir, 0o755)
	_ = os.WriteFile(filepath.Join(pkgDir, "file.txt"), []byte("content"), 0o644)

	installDir := t.TempDir()
	d := packagemanager.NewDefaultManager()
	if err := d.Install(pkgDir, installDir); err != nil {
		t.Fatalf("Install: %v", err)
	}
	pkgs, err := d.List(installDir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(pkgs) == 0 {
		t.Error("expected at least one package after install")
	}
}

func TestDefaultManager_Uninstall_NotExist_Extra3(t *testing.T) {
	d := packagemanager.NewDefaultManager()
	// RemoveAll does not error on non-existent paths; Uninstall should succeed
	err := d.Uninstall("nonexistent", t.TempDir())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBaseManager_ImplementsManager_Extra3(t *testing.T) {
	var _ packagemanager.Manager = packagemanager.NewBaseManager("test")
}

func TestDefaultManager_ImplementsManager_Extra3(t *testing.T) {
	var _ packagemanager.Manager = packagemanager.NewDefaultManager()
}
