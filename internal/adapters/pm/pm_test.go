package pm_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/pm"
)

type mockPM struct{}

func (m *mockPM) Install(name, ver string) error    { return nil }
func (m *mockPM) Uninstall(name string) error       { return nil }
func (m *mockPM) ListInstalled() ([]string, error)  { return []string{"pkg-a"}, nil }
func (m *mockPM) Search(q string) ([]string, error) { return []string{"pkg-a", "pkg-b"}, nil }

// TestParityPMAdapterInterface verifies the interface type exists.
func TestParityPMAdapterInterface(t *testing.T) {
	var _ pm.MCPPackageManagerAdapter = (*mockPM)(nil)
}

// TestParityPMErrors verifies sentinel errors are defined.
func TestParityPMErrors(t *testing.T) {
	if pm.ErrPackageNotFound == nil {
		t.Fatal("ErrPackageNotFound should not be nil")
	}
	if pm.ErrInstallFailed == nil {
		t.Fatal("ErrInstallFailed should not be nil")
	}
}

// TestParityPMListInstalled verifies list via mock.
func TestParityPMListInstalled(t *testing.T) {
	var p pm.MCPPackageManagerAdapter = &mockPM{}
	pkgs, err := p.ListInstalled()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 1 || pkgs[0] != "pkg-a" {
		t.Fatalf("unexpected packages: %v", pkgs)
	}
}
