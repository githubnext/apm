package cachecmd

import "testing"

func TestFormatSize_Bytes_Extra3(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_MB_Extra3(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{1024 * 1024, "1.0 MB"},
		{2 * 1024 * 1024, "2.0 MB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_GB_Extra3(t *testing.T) {
	got := formatSize(1024 * 1024 * 1024)
	if got != "1.0 GB" {
		t.Errorf("formatSize(1 GiB) = %q, want %q", got, "1.0 GB")
	}
}

func TestFormatSize_Negative_Extra3(t *testing.T) {
	// Negative bytes are unusual but should not panic.
	got := formatSize(-1)
	if got == "" {
		t.Error("formatSize(-1) should return a non-empty string")
	}
}
