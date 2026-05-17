package cachecmd

import "testing"

func TestFormatSize(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{5 * 1024 * 1024, "5.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}
	for _, tc := range tests {
		got := formatSize(tc.in)
		if got != tc.want {
			t.Errorf("formatSize(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatSizeBoundaries(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{1, "1 B"},
		{999, "999 B"},
		{1000, "1000 B"},
		{1023, "1023 B"},
		{1025, "1.0 KB"},
		{10 * 1024, "10.0 KB"},
		{100 * 1024, "100.0 KB"},
		{1023 * 1024, "1023.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{2 * 1024 * 1024, "2.0 MB"},
		{100 * 1024 * 1024, "100.0 MB"},
		{2 * 1024 * 1024 * 1024, "2.0 GB"},
		{10 * 1024 * 1024 * 1024, "10.0 GB"},
	}
	for _, tc := range tests {
		got := formatSize(tc.in)
		if got != tc.want {
			t.Errorf("formatSize(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatSizeUnits(t *testing.T) {
	// Verify each unit tier is reachable.
	cases := []struct {
		in   int64
		unit string
	}{
		{500, "B"},
		{2048, "KB"},
		{2 * 1024 * 1024, "MB"},
		{2 * 1024 * 1024 * 1024, "GB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		found := false
		for i := len(got) - 1; i >= 0; i-- {
			if got[i] == ' ' {
				if got[i+1:] == c.unit {
					found = true
				}
				break
			}
		}
		if !found {
			t.Errorf("formatSize(%d) = %q, expected unit %q", c.in, got, c.unit)
		}
	}
}

func TestFormatSizeLargeValues(t *testing.T) {
	// Very large values should still return GB.
	got := formatSize(1 << 40) // 1 TiB -- still formatted as GB
	if len(got) == 0 {
		t.Error("expected non-empty result for large value")
	}
}

func TestFormatSizeNotEmpty(t *testing.T) {
	for _, b := range []int64{0, 1, 100, 10000, 1000000, 1000000000} {
		got := formatSize(b)
		if got == "" {
			t.Errorf("formatSize(%d) returned empty string", b)
		}
	}
}
