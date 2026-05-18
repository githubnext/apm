package cachecmd

import "testing"

func TestFormatSize_Bytes(t *testing.T) {
cases := []struct {
in   int64
want string
}{
{0, "0 B"},
{1, "1 B"},
{999, "999 B"},
{1023, "1023 B"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}

func TestFormatSize_Kilobytes(t *testing.T) {
cases := []struct {
in   int64
want string
}{
{1024, "1.0 KB"},
{1536, "1.5 KB"},
{2048, "2.0 KB"},
{1024 * 10, "10.0 KB"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}

func TestFormatSize_Megabytes(t *testing.T) {
cases := []struct {
in   int64
want string
}{
{1024 * 1024, "1.0 MB"},
{5 * 1024 * 1024, "5.0 MB"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}

func TestFormatSize_Gigabytes(t *testing.T) {
in := int64(1024 * 1024 * 1024)
want := "1.0 GB"
got := formatSize(in)
if got != want {
t.Errorf("formatSize(%d) = %q, want %q", in, got, want)
}
}

func TestFormatSize_BoundaryKB(t *testing.T) {
// 1024*1024 - 1 is still MB boundary
cases := []struct {
in   int64
want string
}{
{1024*1024 - 1, "1024.0 KB"},
{1024 * 512, "512.0 KB"},
{1024 * 100, "100.0 KB"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}

func TestFormatSize_BoundaryMB(t *testing.T) {
cases := []struct {
in   int64
want string
}{
{1024 * 1024 * 100, "100.0 MB"},
{1024 * 1024 * 500, "500.0 MB"},
{1024*1024*1024 - 1, "1024.0 MB"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}

func TestFormatSize_MultipleGB(t *testing.T) {
cases := []struct {
in   int64
want string
}{
{2 * 1024 * 1024 * 1024, "2.0 GB"},
{10 * 1024 * 1024 * 1024, "10.0 GB"},
}
for _, c := range cases {
got := formatSize(c.in)
if got != c.want {
t.Errorf("formatSize(%d) = %q, want %q", c.in, got, c.want)
}
}
}
