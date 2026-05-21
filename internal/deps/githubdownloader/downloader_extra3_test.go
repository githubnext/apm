package githubdownloader

import (
	"testing"
)

func TestRemoteRef_NameAndSHA_Extra3(t *testing.T) {
	r := RemoteRef{
		Name: "refs/tags/v1.0.0",
		SHA:  "abc1234567890abcdef1234567890abcdef123456",
	}
	if r.Name != "refs/tags/v1.0.0" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.SHA != "abc1234567890abcdef1234567890abcdef123456" {
		t.Errorf("SHA = %q", r.SHA)
	}
}

func TestRemoteRef_ZeroFields_Extra3(t *testing.T) {
	var r RemoteRef
	if r.Name != "" || r.SHA != "" {
		t.Error("zero value RemoteRef should have empty fields")
	}
}

func TestDownloadResult_DestDir_Extra3(t *testing.T) {
	d := DownloadResult{
		DestDir:   "/tmp/pkg",
		SHA:       "deadbeef12345678deadbeef12345678deadbeef",
		Ref:       "v1.0.0",
		Transport: "https",
	}
	if d.DestDir != "/tmp/pkg" {
		t.Errorf("DestDir = %q", d.DestDir)
	}
	if d.Transport != "https" {
		t.Errorf("Transport = %q, want https", d.Transport)
	}
}

func TestDownloadResult_SSHTransport_Extra3(t *testing.T) {
	d := DownloadResult{Transport: "ssh"}
	if d.Transport != "ssh" {
		t.Errorf("Transport = %q, want ssh", d.Transport)
	}
}

func TestRawFileResult_Fields_Extra3(t *testing.T) {
	r := RawFileResult{
		Content:     []byte("hello world"),
		ContentType: "text/plain",
		ETag:        `"abc123"`,
	}
	if string(r.Content) != "hello world" {
		t.Errorf("Content = %q", string(r.Content))
	}
	if r.ContentType != "text/plain" {
		t.Errorf("ContentType = %q", r.ContentType)
	}
	if r.ETag != `"abc123"` {
		t.Errorf("ETag = %q", r.ETag)
	}
}

func TestRawFileResult_Zero_Extra3(t *testing.T) {
	var r RawFileResult
	if len(r.Content) != 0 || r.ContentType != "" || r.ETag != "" {
		t.Error("zero value RawFileResult should have empty fields")
	}
}

func TestProtocolPreference_Values_Extra3(t *testing.T) {
	if ProtocolPreferSSH == ProtocolPreferHTTPS {
		t.Error("SSH and HTTPS protocol constants should differ")
	}
}

func TestDefaultOptions_ConcurrencyPositive_Extra3(t *testing.T) {
	opts := DefaultOptions()
	if opts.Concurrency <= 0 {
		t.Errorf("Concurrency = %d, want > 0", opts.Concurrency)
	}
}

func TestSemverSortKey_MajorMinorPatch_Extra3(t *testing.T) {
	k := SemverSortKey("v2.10.3")
	if k[0] != 2 || k[1] != 10 || k[2] != 3 {
		t.Errorf("SemverSortKey(v2.10.3) = %v", k)
	}
}

func TestSemverSortKey_NoVPrefix_Extra3(t *testing.T) {
	k := SemverSortKey("3.0.0")
	if k[0] != 3 || k[1] != 0 || k[2] != 0 {
		t.Errorf("SemverSortKey(3.0.0) = %v", k)
	}
}

func TestSortRemoteRefs_OrderBySemver_Extra3(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0"},
		{Name: "v2.0.0"},
		{Name: "v1.5.0"},
	}
	sorted := SortRemoteRefs(refs)
	if len(sorted) != 3 {
		t.Fatalf("sorted len = %d, want 3", len(sorted))
	}
	if sorted[0].Name != "v2.0.0" {
		t.Errorf("first sorted ref = %q, want v2.0.0", sorted[0].Name)
	}
}

func TestParseLsRemoteOutput_MultipleTags_Extra3(t *testing.T) {
	output := "abc1234\trefs/tags/v1.0.0\ndef5678\trefs/tags/v2.0.0\n"
	refs := ParseLsRemoteOutput(output)
	if len(refs) != 2 {
		t.Fatalf("parsed %d refs, want 2", len(refs))
	}
}
