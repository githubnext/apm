// Package heals implements the heal chain for install-time self-correction.
// Mirrors src/apm_cli/install/heals/base.py, branch_ref_drift.py, and buggy_lockfile_recovery.py.
package heals

// HealMessageLevel indicates the severity of a heal diagnostic message.
type HealMessageLevel int

const (
	HealMessageInfo HealMessageLevel = iota
	HealMessageWarn
)

// HealMessage is a user-facing message emitted by a heal.
type HealMessage struct {
	Level      HealMessageLevel
	Text       string
	PackageKey string
}

// HealContext holds per-dep state threaded through the heal chain.
type HealContext struct {
	PackageKey                      string
	ResolvedRefType                 string // "BRANCH", "TAG", "SHA", ""
	ResolvedCommit                  string // remote HEAD SHA; "" or "cached" if unavailable
	ExistingLockfileApmVersion      string // e.g. "0.12.2"; "" if unknown
	ExistingLockedCommit            string // commit in existing lockfile; "" if none
	LockfileMatch                   bool
	LockfileMatchViaContentHashOnly bool
	UpdateRefs                      bool
	RefChanged                      bool
	BypassKeys                      map[string]bool
	FiredGroups                     map[string]bool
	Messages                        []HealMessage
}

// NewHealContext creates an initialised HealContext for one dependency.
func NewHealContext(
	packageKey string,
	lockfileMatch bool,
	lockfileMatchViaContentHashOnly bool,
	updateRefs bool,
) HealContext {
	return HealContext{
		PackageKey:                      packageKey,
		LockfileMatch:                   lockfileMatch,
		LockfileMatchViaContentHashOnly: lockfileMatchViaContentHashOnly,
		UpdateRefs:                      updateRefs,
		BypassKeys:                      make(map[string]bool),
		FiredGroups:                     make(map[string]bool),
	}
}

// AddBypassKey marks a dep key as having a legitimate hash change.
func (h *HealContext) AddBypassKey(key string) {
	h.BypassKeys[key] = true
}

// Emit appends a user-facing message to the context.
func (h *HealContext) Emit(level HealMessageLevel, text string) {
	h.Messages = append(h.Messages, HealMessage{
		Level:      level,
		Text:       text,
		PackageKey: h.PackageKey,
	})
}

// Heal is the interface each heal struct implements.
type Heal interface {
	Name() string
	Order() int
	ExclusiveGroup() string
	Applies(hctx *HealContext) bool
	Execute(hctx *HealContext)
}

// RunHealChain runs all heals in order, respecting exclusive groups.
func RunHealChain(hctx *HealContext, chain []Heal) {
	for _, h := range chain {
		if eg := h.ExclusiveGroup(); eg != "" {
			if hctx.FiredGroups[eg] {
				continue
			}
		}
		if !h.Applies(hctx) {
			continue
		}
		h.Execute(hctx)
		if eg := h.ExclusiveGroup(); eg != "" {
			hctx.FiredGroups[eg] = true
		}
	}
}

// ----- BranchRefDriftHeal -----

// BranchRefDriftHeal re-downloads when a branch ref's remote SHA has
// advanced past the lockfile-recorded SHA.
// Mirrors src/apm_cli/install/heals/branch_ref_drift.py.
type BranchRefDriftHeal struct{}

func (BranchRefDriftHeal) Name() string           { return "branch_ref_drift" }
func (BranchRefDriftHeal) Order() int              { return 10 }
func (BranchRefDriftHeal) ExclusiveGroup() string  { return "branch_drift" }

func (BranchRefDriftHeal) Applies(hctx *HealContext) bool {
	if !hctx.LockfileMatch || hctx.UpdateRefs {
		return false
	}
	if hctx.ResolvedRefType != "BRANCH" {
		return false
	}
	remoteSHA := hctx.ResolvedCommit
	if remoteSHA == "" || remoteSHA == "cached" {
		return false
	}
	if hctx.ExistingLockedCommit == "" || hctx.ExistingLockedCommit == "cached" {
		return false
	}
	return remoteSHA != hctx.ExistingLockedCommit
}

func (BranchRefDriftHeal) Execute(hctx *HealContext) {
	lockedSHA := hctx.ExistingLockedCommit
	remoteSHA := hctx.ResolvedCommit
	shortLocked := lockedSHA
	if len(shortLocked) > 8 {
		shortLocked = shortLocked[:8]
	}
	shortRemote := remoteSHA
	if len(shortRemote) > 8 {
		shortRemote = shortRemote[:8]
	}
	hctx.LockfileMatch = false
	hctx.RefChanged = true
	hctx.AddBypassKey(hctx.PackageKey)
	hctx.Emit(
		HealMessageInfo,
		"  Branch ref drift: "+hctx.PackageKey+" remote @"+shortRemote+
			" != locked @"+shortLocked+" -- forcing re-download",
	)
}

// ----- BuggyLockfileRecoveryHeal -----

// buggyBranchRefDriftVersions lists APM versions known to produce
// phantom resolved_commit values in branch-ref deps.
var buggyBranchRefDriftVersions = map[string]bool{
	"0.10.0": true, "0.10.1": true, "0.10.2": true,
	"0.11.0": true, "0.11.1": true, "0.11.2": true,
	"0.12.0": true, "0.12.1": true, "0.12.2": true,
}

// BuggyLockfileRecoveryHeal recovers from the v<=0.12.2 branch-ref cache drift bug.
// Mirrors src/apm_cli/install/heals/buggy_lockfile_recovery.py.
type BuggyLockfileRecoveryHeal struct{}

func (BuggyLockfileRecoveryHeal) Name() string           { return "buggy_lockfile_recovery" }
func (BuggyLockfileRecoveryHeal) Order() int              { return 20 }
func (BuggyLockfileRecoveryHeal) ExclusiveGroup() string  { return "branch_drift" }

func (BuggyLockfileRecoveryHeal) Applies(hctx *HealContext) bool {
	if !hctx.LockfileMatch {
		return false
	}
	if !hctx.LockfileMatchViaContentHashOnly {
		return false
	}
	if hctx.UpdateRefs {
		return false
	}
	if hctx.ResolvedRefType != "BRANCH" {
		return false
	}
	return buggyBranchRefDriftVersions[hctx.ExistingLockfileApmVersion]
}

func (BuggyLockfileRecoveryHeal) Execute(hctx *HealContext) {
	hctx.LockfileMatch = false
	hctx.RefChanged = true
	hctx.AddBypassKey(hctx.PackageKey)
	hctx.Emit(
		HealMessageWarn,
		"Recovering "+hctx.PackageKey+" from "+
			"branch-ref cache drift in lockfile generated by APM <= 0.12.2 "+
			"-- forcing re-download to restore consistency. "+
			"Upgrade APM (>= 0.13.0) to prevent recurrence.",
	)
}

// DefaultHealChain returns the standard ordered heal chain.
func DefaultHealChain() []Heal {
	return []Heal{
		BranchRefDriftHeal{},
		BuggyLockfileRecoveryHeal{},
	}
}
