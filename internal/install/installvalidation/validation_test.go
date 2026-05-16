package installvalidation_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/installvalidation"
)

func TestAuthenticationError_WithHost(t *testing.T) {
	err := &installvalidation.AuthenticationError{Host: "github.com", Message: "bad token"}
	msg := err.Error()
	if !strings.Contains(msg, "github.com") {
		t.Fatalf("expected host in error, got %q", msg)
	}
	if !strings.Contains(msg, "bad token") {
		t.Fatalf("expected message, got %q", msg)
	}
}

func TestAuthenticationError_NoHost(t *testing.T) {
	err := &installvalidation.AuthenticationError{Message: "no creds"}
	msg := err.Error()
	if !strings.Contains(msg, "no creds") {
		t.Fatalf("expected message, got %q", msg)
	}
}

func TestTLSError_WithHost(t *testing.T) {
	cause := errors.New("x509: cert invalid")
	err := &installvalidation.TLSError{Host: "ghe.example.com", Cause: cause}
	msg := err.Error()
	if !strings.Contains(msg, installvalidation.TLSErrorPrefix) {
		t.Fatalf("expected TLS prefix, got %q", msg)
	}
	if !strings.Contains(msg, "ghe.example.com") {
		t.Fatalf("expected host, got %q", msg)
	}
	if err.Unwrap() != cause {
		t.Fatal("Unwrap should return cause")
	}
}

func TestIsTLSFailure_True(t *testing.T) {
	err := &installvalidation.TLSError{Host: "x", Cause: errors.New("cert")}
	if !installvalidation.IsTLSFailure(err) {
		t.Fatal("IsTLSFailure should be true")
	}
}

func TestIsTLSFailure_False(t *testing.T) {
	if installvalidation.IsTLSFailure(errors.New("other")) {
		t.Fatal("IsTLSFailure should be false for non-TLS error")
	}
	if installvalidation.IsTLSFailure(nil) {
		t.Fatal("IsTLSFailure(nil) should be false")
	}
}

func TestLocalPathFailureReason_Missing(t *testing.T) {
	reason := installvalidation.LocalPathFailureReason("/nonexistent/path/to/pkg")
	if reason == "" {
		t.Fatal("expected a failure reason for missing path")
	}
}
