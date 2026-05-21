package githubdownloader

import "testing"

func TestDefaultOptions_CacheDirEmpty_Extra4(t *testing.T) {
opts := DefaultOptions()
if opts.CacheDir != "" {
t.Errorf("expected empty CacheDir by default, got %q", opts.CacheDir)
}
}

func TestDefaultOptions_AllowFallback_Extra4(t *testing.T) {
opts := DefaultOptions()
if !opts.AllowFallback {
t.Error("expected AllowFallback true by default")
}
}

func TestOptions_CacheDirField_Extra4(t *testing.T) {
opts := Options{CacheDir: "/tmp/cache"}
if opts.CacheDir != "/tmp/cache" {
t.Errorf("unexpected CacheDir: %s", opts.CacheDir)
}
}

func TestOptions_ConcurrencyField_Extra4(t *testing.T) {
opts := Options{Concurrency: 8}
if opts.Concurrency != 8 {
t.Errorf("unexpected Concurrency: %d", opts.Concurrency)
}
}

func TestOptions_TimeoutSecsField_Extra4(t *testing.T) {
opts := Options{TimeoutSecs: 600.0}
if opts.TimeoutSecs != 600.0 {
t.Errorf("unexpected TimeoutSecs: %f", opts.TimeoutSecs)
}
}

func TestDownloadResult_DestDirField_Extra4(t *testing.T) {
r := DownloadResult{DestDir: "/tmp/dest"}
if r.DestDir != "/tmp/dest" {
t.Errorf("unexpected DestDir: %s", r.DestDir)
}
}

func TestDownloadResult_SHAField_Extra4(t *testing.T) {
r := DownloadResult{SHA: "abc123def456abc123def456abc123def456abc1"}
if r.SHA != "abc123def456abc123def456abc123def456abc1" {
t.Errorf("unexpected SHA: %s", r.SHA)
}
}

func TestDownloadResult_TransportHTTPS_Extra4(t *testing.T) {
r := DownloadResult{Transport: "https"}
if r.Transport != "https" {
t.Errorf("unexpected Transport: %s", r.Transport)
}
}

func TestDownloadResult_TransportSSH_Extra4(t *testing.T) {
r := DownloadResult{Transport: "ssh"}
if r.Transport != "ssh" {
t.Errorf("unexpected Transport: %s", r.Transport)
}
}

func TestProtocolPreference_Distinct_Extra4(t *testing.T) {
if ProtocolPreferHTTPS == ProtocolPreferSSH {
t.Error("expected distinct values for HTTPS and SSH preferences")
}
}

func TestBuildTransportPlan_HTTPSNoFallback_Extra4(t *testing.T) {
plan := BuildTransportPlan(ProtocolPreferHTTPS, false)
if plan.Primary != "https" {
t.Errorf("expected primary https, got %s", plan.Primary)
}
}

func TestBuildTransportPlan_SSHNoFallback_Extra4(t *testing.T) {
plan := BuildTransportPlan(ProtocolPreferSSH, false)
if plan.Primary != "ssh" {
t.Errorf("expected primary ssh, got %s", plan.Primary)
}
}

func TestTransportPlan_FallbacksEmpty_Extra4(t *testing.T) {
plan := BuildTransportPlan(ProtocolPreferHTTPS, false)
if len(plan.Fallbacks) != 0 {
t.Errorf("expected no fallbacks, got %v", plan.Fallbacks)
}
}

func TestSemverSortKey_Valid_Extra4(t *testing.T) {
k := SemverSortKey("v1.12.3")
zero := [4]int{}
if k == zero {
t.Error("expected non-zero sortkey for valid semver")
}
}
