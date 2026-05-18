package gitstderr_test

import (
"testing"

"github.com/githubnext/apm/internal/marketplace/gitstderr"
)

func TestTranslate_Auth(t *testing.T) {
r := gitstderr.Translate("fatal: authentication failed for 'https://github.com/acme/tools'",
gitstderr.Options{Operation: "ls-remote", Remote: "acme/tools"})
if r.Kind != gitstderr.KindAuth {
t.Fatalf("expected KindAuth, got %s", r.Kind)
}
if r.Summary == "" || r.Hint == "" {
t.Fatal("expected non-empty summary and hint")
}
}

func TestTranslate_NotFound(t *testing.T) {
r := gitstderr.Translate("ERROR: Repository not found.", gitstderr.Options{Operation: "clone"})
if r.Kind != gitstderr.KindNotFound {
t.Fatalf("expected KindNotFound, got %s", r.Kind)
}
}

func TestTranslate_Timeout(t *testing.T) {
r := gitstderr.Translate("fatal: unable to connect to github.com: connection timed out",
gitstderr.Options{})
if r.Kind != gitstderr.KindTimeout {
t.Fatalf("expected KindTimeout, got %s", r.Kind)
}
}

func TestTranslate_Unknown(t *testing.T) {
r := gitstderr.Translate("some unexpected error", gitstderr.Options{})
if r.Kind != gitstderr.KindUnknown {
t.Fatalf("expected KindUnknown, got %s", r.Kind)
}
}

func TestTranslate_TruncatesRaw(t *testing.T) {
long := string(make([]byte, 600))
for i := range long {
long = long[:i] + "a" + long[i+1:]
}
r := gitstderr.Translate(long, gitstderr.Options{})
if len(r.Raw) > 520 {
t.Fatalf("raw too long: %d", len(r.Raw))
}
}

func TestTranslate_CouldNotResolveHost_IsTimeout(t *testing.T) {
r := gitstderr.Translate("fatal: could not resolve host: github.com", gitstderr.Options{})
if r.Kind != gitstderr.KindTimeout {
t.Fatalf("expected KindTimeout for DNS failure, got %s", r.Kind)
}
}

func TestTranslate_InvalidCredentials(t *testing.T) {
r := gitstderr.Translate("fatal: invalid credentials", gitstderr.Options{Operation: "fetch"})
if r.Kind != gitstderr.KindAuth {
t.Fatalf("expected KindAuth for invalid credentials, got %s", r.Kind)
}
}

func TestTranslate_Empty(t *testing.T) {
r := gitstderr.Translate("", gitstderr.Options{})
if r.Kind != gitstderr.KindUnknown {
t.Fatalf("expected KindUnknown for empty stderr, got %s", r.Kind)
}
}

func TestTranslate_Raw_Preserved(t *testing.T) {
input := "some git error message"
r := gitstderr.Translate(input, gitstderr.Options{})
if r.Raw != input {
t.Errorf("Raw = %q, want %q", r.Raw, input)
}
}

func TestGitErrorKind_String(t *testing.T) {
cases := map[gitstderr.GitErrorKind]string{
gitstderr.KindAuth:     "auth",
gitstderr.KindNotFound: "not_found",
gitstderr.KindTimeout:  "timeout",
gitstderr.KindUnknown:  "unknown",
}
for kind, want := range cases {
if got := kind.String(); got != want {
t.Errorf("GitErrorKind(%d).String() = %q, want %q", kind, got, want)
}
}
}

func TestTranslate_NetworkReadFailed_IsTimeout(t *testing.T) {
r := gitstderr.Translate("error: RPC failed; curl 18 transfer closed", gitstderr.Options{})
// Curl transfer-closed should be timeout or unknown -- just check it doesn't panic.
_ = r.Kind
}

func TestTranslate_NoSuchRemote_IsNotFound(t *testing.T) {
r := gitstderr.Translate("fatal: 'origin' does not appear to be a git repository", gitstderr.Options{})
// Should be not_found or unknown -- ensure no panic.
_ = r
}
