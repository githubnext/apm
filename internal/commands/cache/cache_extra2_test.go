package cachecmd

import "testing"

func TestFormatSize_KB(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{1536, "1.5 KB"},
		{1024*1024 - 1, "1024.0 KB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_MB(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{1024 * 1024, "1.0 MB"},
		{2 * 1024 * 1024, "2.0 MB"},
		{int64(1.5 * 1024 * 1024), "1.5 MB"},
	}
	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSize_GB(t *testing.T) {
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

func TestFormatSize_NegativeIsBytes(t *testing.T) {
	got := formatSize(-1)
	if got != "-1 B" {
		t.Errorf("negative should be treated as bytes, got %q", got)
	}
}

func TestFormatSize_Zero(t *testing.T) {
	got := formatSize(0)
	if got != "0 B" {
		t.Errorf("zero should be '0 B', got %q", got)
	}
}

func TestFormatSize_LargeGB(t *testing.T) {
	got := formatSize(100 * 1024 * 1024 * 1024)
	if got != "100.0 GB" {
		t.Errorf("100 GB expected '100.0 GB', got %q", got)
	}
}
