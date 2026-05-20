package semver_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/semver"
)

func TestParse_ZeroPatch(t *testing.T) {
	v, err := semver.Parse("1.2.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 0 {
		t.Errorf("unexpected version: %+v", v)
	}
}

func TestParse_AllZero(t *testing.T) {
	v, err := semver.Parse("0.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 0 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("unexpected version: %+v", v)
	}
}

func TestParse_LargeNumbers(t *testing.T) {
	v, err := semver.Parse("100.200.300")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 100 || v.Minor != 200 || v.Patch != 300 {
		t.Errorf("unexpected version: %+v", v)
	}
}

func TestCompare_PatchDifference(t *testing.T) {
	a, _ := semver.Parse("1.0.1")
	b, _ := semver.Parse("1.0.0")
	if a.Compare(b) != 1 {
		t.Error("1.0.1 should be > 1.0.0")
	}
	if b.Compare(a) != -1 {
		t.Error("1.0.0 should be < 1.0.1")
	}
}

func TestCompare_SameVersion(t *testing.T) {
	a, _ := semver.Parse("2.3.4")
	b, _ := semver.Parse("2.3.4")
	if a.Compare(b) != 0 {
		t.Error("equal versions should compare 0")
	}
}

func TestSatisfiesRange_ExactVersionMatch(t *testing.T) {
	v, _ := semver.Parse("1.2.3")
	if !semver.SatisfiesRange(v, "1.2.3") {
		t.Error("1.2.3 should satisfy range 1.2.3")
	}
}

func TestSatisfiesRange_ExactVersionNoMatch(t *testing.T) {
	v, _ := semver.Parse("1.2.4")
	if semver.SatisfiesRange(v, "1.2.3") {
		t.Error("1.2.4 should not satisfy range 1.2.3")
	}
}

func TestSatisfiesRange_GreaterThan(t *testing.T) {
	v, _ := semver.Parse("2.0.0")
	if !semver.SatisfiesRange(v, ">1.0.0") {
		t.Error("2.0.0 should satisfy >1.0.0")
	}
}

func TestSatisfiesRange_LessThan(t *testing.T) {
	v, _ := semver.Parse("0.9.0")
	if !semver.SatisfiesRange(v, "<1.0.0") {
		t.Error("0.9.0 should satisfy <1.0.0")
	}
}

func TestSatisfiesRange_GreaterOrEqual(t *testing.T) {
	v, _ := semver.Parse("1.0.0")
	if !semver.SatisfiesRange(v, ">=1.0.0") {
		t.Error("1.0.0 should satisfy >=1.0.0")
	}
}

func TestSatisfiesRange_LessOrEqual(t *testing.T) {
	v, _ := semver.Parse("1.0.0")
	if !semver.SatisfiesRange(v, "<=1.0.0") {
		t.Error("1.0.0 should satisfy <=1.0.0")
	}
}

func TestSatisfiesRange_CaretMinorRange(t *testing.T) {
	v, _ := semver.Parse("1.5.0")
	if !semver.SatisfiesRange(v, "^1.0.0") {
		t.Error("1.5.0 should satisfy ^1.0.0")
	}
}

func TestSatisfiesRange_CaretOutOfMajor(t *testing.T) {
	v, _ := semver.Parse("2.0.0")
	if semver.SatisfiesRange(v, "^1.0.0") {
		t.Error("2.0.0 should not satisfy ^1.0.0")
	}
}

func TestParse_InvalidEmpty(t *testing.T) {
	_, err := semver.Parse("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestParse_InvalidAlpha(t *testing.T) {
	_, err := semver.Parse("not-a-version")
	if err == nil {
		t.Error("expected error for non-semver string")
	}
}
