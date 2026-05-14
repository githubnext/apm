// Package gitstderr translates git stderr into actionable, ASCII-only error messages.
//
// Callers pass captured stderr text, an optional exit code, and context
// (operation name, remote). This package classifies the failure into one
// of four known modes and returns a structured TranslatedGitError with a
// one-line summary, an actionable hint, and the (truncated) raw stderr.
//
// No subprocess, network, filesystem, or logging side effects -- this is
// a pure function package.
package gitstderr

import (
"fmt"
"strings"
)

const (
rawMaxLen     = 500
summaryMaxLen = 80
)

// GitErrorKind enumerates known git failure modes.
type GitErrorKind int

const (
// KindAuth indicates an authentication failure.
KindAuth GitErrorKind = iota
// KindNotFound indicates a ref or repository not found failure.
KindNotFound
// KindTimeout indicates a network timeout or connectivity failure.
KindTimeout
// KindUnknown indicates an unclassified failure.
KindUnknown
)

// String returns the value string for GitErrorKind.
func (k GitErrorKind) String() string {
switch k {
case KindAuth:
return "auth"
case KindNotFound:
return "not_found"
case KindTimeout:
return "timeout"
default:
return "unknown"
}
}

// TranslatedGitError is the structured result of translating git stderr.
type TranslatedGitError struct {
Kind    GitErrorKind
Summary string
Hint    string
Raw     string
}

var authPatterns = []string{
"authentication failed",
"invalid credentials",
"could not read password",
"permission denied (publickey)",
"403 forbidden",
"401 unauthorized",
"fatal: authentication",
"remote: write access",
"please make sure you have the correct access rights",
"the requested url returned error: 401",
"the requested url returned error: 403",
}

var notFoundPatterns = []string{
"repository not found",
"does not appear to be a git repository",
"not a valid ref",
"couldn't find remote ref",
"could not resolve",
"the requested url returned error: 404",
"no such ref",
"unknown ref",
}

var timeoutPatterns = []string{
"operation timed out",
"connection timed out",
"could not resolve host",
"connection refused",
"network is unreachable",
"temporary failure in name resolution",
"ssl_read: connection reset",
"early eof",
"rpc failed",
}

func truncateRaw(stderr string) string {
if len(stderr) <= rawMaxLen {
return stderr
}
return stderr[:rawMaxLen] + "... (truncated)"
}

func classify(stderrLower string) GitErrorKind {
for _, p := range authPatterns {
if strings.Contains(stderrLower, p) {
return KindAuth
}
}
for _, p := range notFoundPatterns {
if strings.Contains(stderrLower, p) {
// "could not resolve host" is a DNS/network issue, not not-found.
if p == "could not resolve" && strings.Contains(stderrLower, "could not resolve host") {
continue
}
return KindNotFound
}
}
for _, p := range timeoutPatterns {
if strings.Contains(stderrLower, p) {
return KindTimeout
}
}
return KindUnknown
}

func buildSummary(kind GitErrorKind, operation string, exitCode *int) string {
var text string
switch kind {
case KindAuth:
text = fmt.Sprintf("Git authentication failed during %s.", operation)
case KindNotFound:
text = fmt.Sprintf("Git ref or repository not found during %s.", operation)
case KindTimeout:
text = fmt.Sprintf("Git network timeout during %s.", operation)
default:
if exitCode != nil {
text = fmt.Sprintf("Git failed during %s (exit %d).", operation, *exitCode)
} else {
text = fmt.Sprintf("Git failed during %s.", operation)
}
}
if len(text) > summaryMaxLen {
text = text[:summaryMaxLen-3] + "..."
}
return text
}

func buildHint(kind GitErrorKind, operation string, remote string) string {
switch kind {
case KindAuth:
return "Check your GITHUB_TOKEN / gh auth / SSH key. Run 'apm marketplace doctor' to diagnose."
case KindNotFound:
remoteLabel := "the remote"
if remote != "" {
remoteLabel = "'" + remote + "'"
}
return fmt.Sprintf("Verify the remote %s exists and the ref is spelled correctly.", remoteLabel)
case KindTimeout:
return "Network issue contacting the remote. Retry or check your connection."
default:
return fmt.Sprintf("Git failed during %s. See raw stderr above.", operation)
}
}

// Options configures a Translate call.
type Options struct {
// ExitCode is the optional exit code from git. Pass nil if unknown.
ExitCode *int
// Operation names the git operation (e.g. "ls-remote"). Defaults to "git operation".
Operation string
// Remote is the optional remote name or URL for the hint.
Remote string
}

// Translate classifies git stderr text into a known failure mode and produces an actionable hint.
func Translate(stderr string, opts Options) TranslatedGitError {
if opts.Operation == "" {
opts.Operation = "git operation"
}
kind := classify(strings.ToLower(stderr))
return TranslatedGitError{
Kind:    kind,
Summary: buildSummary(kind, opts.Operation, opts.ExitCode),
Hint:    buildHint(kind, opts.Operation, opts.Remote),
Raw:     truncateRaw(stderr),
}
}
