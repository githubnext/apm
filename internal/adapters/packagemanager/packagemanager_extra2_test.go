package packagemanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaseManager_NameEmpty_Extra2(t *testing.T) {
	b := NewBaseManager("")
	if b.Name() != "" {
		t.Errorf("expected empty name, got %q", b.Name())
	}
}

func TestBaseManager_InstallNotImplemented_Extra2(t *testing.T) {
	b := NewBaseManager("mymanager")
	err := b.Install("/some/pkg", "/some/dir")
	if err == nil {
		t.Error("expected error from BaseManager.Install")
	}
}

func TestBaseManager_UninstallNotImplemented_Extra2(t *testing.T) {
	b := NewBaseManager("mymanager")
	err := b.Uninstall("pkg", "/some/dir")
	if err == nil {
		t.Error("expected error from BaseManager.Uninstall")
	}
}

func TestBaseManager_ListNotImplemented_Extra2(t *testing.T) {
	b := NewBaseManager("mymanager")
	_, err := b.List("/some/dir")
	if err == nil {
		t.Error("expected error from BaseManager.List")
	}
}

func TestBaseManager_IsSupportedFalse_Extra2(t *testing.T) {
	b := NewBaseManager("mymanager")
	if b.IsSupported("/any/path") {
		t.Error("expected IsSupported=false for BaseManager")
	}
}

func TestDefaultManager_IsSupportedAny_Extra2(t *testing.T) {
	d := NewDefaultManager()
	if !d.IsSupported("") {
		t.Error("expected IsSupported=true for any path")
	}
	if !d.IsSupported("/some/random/path.tar.gz") {
		t.Error("expected IsSupported=true for any path")
	}
}

func TestDefaultManager_Name_Extra2(t *testing.T) {
	d := NewDefaultManager()
	if d.Name() != "default" {
		t.Errorf("expected 'default', got %q", d.Name())
	}
}

func TestDefaultManager_ListEmptyDir_Extra2(t *testing.T) {
	dir := t.TempDir()
	d := NewDefaultManager()
	list, err := d.List(dir)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %v", list)
	}
}

func TestDefaultManager_ListNonexistent_Extra2(t *testing.T) {
	d := NewDefaultManager()
	list, err := d.List("/nonexistent/path/xyz123")
	if err != nil {
		t.Fatalf("unexpected error for nonexistent dir: %v", err)
	}
	if list != nil {
		t.Errorf("expected nil list, got %v", list)
	}
}

func TestDefaultManager_InstallAndList_Extra2(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	pkgDir := filepath.Join(src, "mypkg")
	_ = os.MkdirAll(pkgDir, 0o755)
	_ = os.WriteFile(filepath.Join(pkgDir, "file.txt"), []byte("hello"), 0o644)

	d := NewDefaultManager()
	if err := d.Install(pkgDir, dst); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	list, err := d.List(dst)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	found := false
	for _, name := range list {
		if name == "mypkg" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected mypkg in list, got %v", list)
	}
}

func TestDefaultManager_UninstallAndList_Extra2(t *testing.T) {
	dst := t.TempDir()
	pkgDir := filepath.Join(dst, "mypkg")
	_ = os.MkdirAll(pkgDir, 0o755)

	d := NewDefaultManager()
	if err := d.Uninstall("mypkg", dst); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}

	list, err := d.List(dst)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	for _, name := range list {
		if name == "mypkg" {
			t.Error("mypkg should have been uninstalled")
		}
	}
}

func TestDefaultManager_InstallMissingSrc_Extra2(t *testing.T) {
	dst := t.TempDir()
	d := NewDefaultManager()
	err := d.Install("/nonexistent/mypkg", dst)
	if err == nil {
		t.Error("expected error when source does not exist")
	}
}

func TestDefaultManager_ListOnlyDirs_Extra2(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hi"), 0o644)
	sub := filepath.Join(dir, "subpkg")
	_ = os.MkdirAll(sub, 0o755)

	d := NewDefaultManager()
	list, err := d.List(dir)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	for _, name := range list {
		if name == "file.txt" {
			t.Error("List should not include files, only dirs")
		}
	}
	found := false
	for _, name := range list {
		if name == "subpkg" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected subpkg in list, got %v", list)
	}
}

func TestManagerInterface_Extra2(t *testing.T) {
	var _ Manager = NewDefaultManager()
	var _ Manager = NewBaseManager("x")
}
