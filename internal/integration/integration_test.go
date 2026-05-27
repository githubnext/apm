package integration_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration"
)

// TestParityIntegrationResultFields verifies field parity with Python IntegrationResult.
func TestParityIntegrationResultFields(t *testing.T) {
	r := integration.IntegrationResult{
		FilesIntegrated: 3,
		FilesSkipped:    1,
		FilesAdopted:    2,
		LinksResolved:   5,
	}
	if r.FilesIntegrated != 3 {
		t.Fatalf("expected 3, got %d", r.FilesIntegrated)
	}
	if r.FilesSkipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", r.FilesSkipped)
	}
}

// TestParityIntegrationResultTotal verifies Total() helper.
func TestParityIntegrationResultTotal(t *testing.T) {
	r := integration.IntegrationResult{
		FilesIntegrated: 3,
		FilesSkipped:    1,
		FilesAdopted:    2,
	}
	if r.Total() != 6 {
		t.Fatalf("expected total 6, got %d", r.Total())
	}
}

// TestParityIntegrationResultZero verifies zero value is valid.
func TestParityIntegrationResultZero(t *testing.T) {
	var r integration.IntegrationResult
	if r.Total() != 0 {
		t.Fatalf("expected 0, got %d", r.Total())
	}
}

// TestParityIntegrationErrors verifies sentinel errors.
func TestParityIntegrationErrors(t *testing.T) {
	if integration.ErrIntegrationConflict == nil {
		t.Fatal("ErrIntegrationConflict should not be nil")
	}
	if integration.ErrIntegrationSkipped == nil {
		t.Fatal("ErrIntegrationSkipped should not be nil")
	}
}

// TestParityIntegrateOptionsFields verifies option struct.
func TestParityIntegrateOptionsFields(t *testing.T) {
	opts := integration.IntegrateOptions{DryRun: true, Force: false, Global: true}
	if !opts.DryRun {
		t.Fatal("expected DryRun=true")
	}
	if !opts.Global {
		t.Fatal("expected Global=true")
	}
}

// TestParityIntegratorInterface verifies the Integrator interface.
func TestParityIntegratorInterface(t *testing.T) {
	var _ integration.Integrator = (*mockIntegrator)(nil)
}

type mockIntegrator struct{}

func (m *mockIntegrator) Integrate(opts integration.IntegrateOptions) (integration.IntegrationResult, error) {
	return integration.IntegrationResult{FilesIntegrated: 1}, nil
}
func (m *mockIntegrator) Name() string { return "mock" }

// TestParityIntegratorMockRun verifies the interface can be called.
func TestParityIntegratorMockRun(t *testing.T) {
	var i integration.Integrator = &mockIntegrator{}
	result, err := i.Integrate(integration.IntegrateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FilesIntegrated != 1 {
		t.Fatalf("expected 1 file integrated, got %d", result.FilesIntegrated)
	}
}

// TestParityIntegrationResultSkillCreated verifies SkillCreated field.
func TestParityIntegrationResultSkillCreated(t *testing.T) {
	r := integration.IntegrationResult{SkillCreated: true}
	if !r.SkillCreated {
		t.Fatal("expected SkillCreated=true")
	}
}

// TestParityIntegrationResultSubSkills verifies SubSkillsPromoted field.
func TestParityIntegrationResultSubSkills(t *testing.T) {
	r := integration.IntegrationResult{SubSkillsPromoted: 3}
	if r.SubSkillsPromoted != 3 {
		t.Fatalf("expected 3, got %d", r.SubSkillsPromoted)
	}
}
