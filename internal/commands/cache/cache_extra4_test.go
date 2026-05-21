package cachecmd

import "testing"

func TestFormatSize_GB_Extra4(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{1024 * 1024 * 1024, "1.0 GB"},
		{2 * 1024 * 1024 * 1024, "2.0 GB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_KB_Extra4(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{512 * 1024, "512.0 KB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_SmallBytes_Extra4(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{100, "100 B"},
		{512, "512 B"},
		{999, "999 B"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_FractionalMB_Extra4(t *testing.T) {
	got := formatSize(int64(1.5 * 1024 * 1024))
	if got == "" {
		t.Error("formatSize should return non-empty string")
	}
}

func TestFormatSize_LargeGB_Extra4(t *testing.T) {
	got := formatSize(10 * 1024 * 1024 * 1024)
	if got == "" {
		t.Error("formatSize should return non-empty string for large values")
	}
}
