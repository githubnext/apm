// Package outcomerouting is the single source of truth for the 9-outcome
// policy-discovery routing table.
// Migrated from src/apm_cli/policy/outcome_routing.py.
package outcomerouting

import (
	"fmt"

	"github.com/githubnext/apm/internal/policy/schema"
)

// PolicyViolationError is raised when a policy demands fail-closed behaviour.
type PolicyViolationError struct {
	Message      string
	PolicySource string
}

func (e *PolicyViolationError) Error() string {
	return e.Message
}

// PolicyFetchResult holds the result of a discover_policy call.
type PolicyFetchResult struct {
	Outcome         string
	Source          string
	Cached          bool
	Error           string
	FetchError      string
	CacheAgeSeconds int
	Policy          *schema.ApmPolicy
}

// PolicyLogger is the minimal interface expected of a logger for routing.
type PolicyLogger interface {
	PolicyResolved(source string, cached bool, enforcement string, ageSeconds int)
	PolicyDiscoveryMiss(outcome string, source string, err string)
}

// outcomesHonoringFetchFailureDefault is the set of outcomes that respect the
// project-side policy.fetch_failure_default knob.
var outcomesHonoringFetchFailureDefault = map[string]bool{
	"malformed":           true,
	"cache_miss_fetch_fail": true,
	"garbage_response":    true,
	"no_git_remote":       true,
	"absent":              true,
	"empty":               true,
}

// nonFoundLoggedOutcomes is the set of outcomes routed through the canonical
// policy_discovery_miss logger helper.
var nonFoundLoggedOutcomes = map[string]bool{
	"absent":              true,
	"no_git_remote":       true,
	"empty":               true,
	"malformed":           true,
	"cache_miss_fetch_fail": true,
	"garbage_response":    true,
}

// RouteDiscoveryOutcome routes a PolicyFetchResult to logging and fail-closed
// decisions.
//
// Parameters:
//   - fetchResult: result of discover_policy_with_chain
//   - logger: logger implementing PolicyLogger (nil is tolerated)
//   - fetchFailureDefault: project-side policy.fetch_failure_default ("warn" or "block")
//   - raiseBlockingErrors: when true, return a PolicyViolationError for blocking outcomes
//
// Returns the effective ApmPolicy when enforcement should proceed, nil otherwise.
// When raiseBlockingErrors is true and a blocking condition is met, a non-nil error
// is returned alongside a nil policy.
func RouteDiscoveryOutcome(
	fetchResult PolicyFetchResult,
	logger PolicyLogger,
	fetchFailureDefault string,
	raiseBlockingErrors bool,
) (*schema.ApmPolicy, error) {
	outcome := fetchResult.Outcome
	source := fetchResult.Source

	if outcome == "disabled" {
		return nil, nil
	}

	// hash_mismatch: ALWAYS fail closed regardless of fetch_failure_default.
	if outcome == "hash_mismatch" {
		errStr := fetchResult.Error
		if errStr == "" {
			errStr = fetchResult.FetchError
		}
		if logger != nil {
			logger.PolicyDiscoveryMiss("hash_mismatch", source, errStr)
		}
		if raiseBlockingErrors {
			return nil, &PolicyViolationError{
				Message: fmt.Sprintf(
					"Install blocked: policy hash mismatch -- pinned policy.hash "+
						"does not match fetched policy bytes (source=%s). "+
						"Update apm.yml policy.hash or contact your org admin.",
					sourceOrUnknown(source),
				),
				PolicySource: sourceOrUnknown(source),
			}
		}
		return nil, nil
	}

	// 6 of 9 non-found outcomes route through the canonical logger helper.
	if nonFoundLoggedOutcomes[outcome] {
		errStr := fetchResult.Error
		if errStr == "" {
			errStr = fetchResult.FetchError
		}
		if logger != nil {
			logger.PolicyDiscoveryMiss(outcome, source, errStr)
		}
		if raiseBlockingErrors &&
			outcomesHonoringFetchFailureDefault[outcome] &&
			fetchFailureDefault == "block" {
			return nil, &PolicyViolationError{
				Message: fmt.Sprintf(
					"Install blocked: no enforceable org policy was resolved "+
						"(outcome=%s) and project apm.yml has "+
						"policy.fetch_failure_default=block (source=%s)",
					outcome,
					sourceOrUnknown(source),
				),
				PolicySource: sourceOrUnknown(source),
			}
		}
		return nil, nil
	}

	// cached_stale: log, enforce with the cached policy, potentially fail closed.
	if outcome == "cached_stale" {
		policy := fetchResult.Policy
		if logger != nil {
			if policy != nil {
				enforcement := policy.Enforcement
				if enforcement == "" {
					enforcement = "warn"
				}
				logger.PolicyResolved(source, true, enforcement, fetchResult.CacheAgeSeconds)
			}
			logger.PolicyDiscoveryMiss("cached_stale", source, fetchResult.FetchError)
		}
		if raiseBlockingErrors && policy != nil {
			ff := policy.FetchFailure
			if ff == "" {
				ff = "warn"
			}
			if ff == "block" {
				return nil, &PolicyViolationError{
					Message: fmt.Sprintf(
						"Install blocked: org policy refresh failed and the cached "+
							"policy declares fetch_failure=block (source=%s)",
						sourceOrUnknown(source),
					),
					PolicySource: sourceOrUnknown(source),
				}
			}
		}
		return policy, nil
	}

	// found: normal path
	if outcome == "found" {
		policy := fetchResult.Policy
		if logger != nil && policy != nil {
			enforcement := policy.Enforcement
			if enforcement == "" {
				enforcement = "warn"
			}
			logger.PolicyResolved(source, fetchResult.Cached, enforcement, fetchResult.CacheAgeSeconds)
		}
		return policy, nil
	}

	// Defensive: unrecognised outcome -- skip enforcement.
	return nil, nil
}

func sourceOrUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
