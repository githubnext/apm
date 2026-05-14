// Package heal implements the heal-chain dispatcher for per-dep mid-flow
// correction during the install pipeline.
// Mirrors src/apm_cli/install/phases/heal.py.
package heal

// HealMessageLevel indicates the severity of a heal diagnostic message.
type HealMessageLevel int

const (
	HealMessageInfo HealMessageLevel = iota
	HealMessageWarn
)

// HealMessage is a user-facing message emitted by a healer.
type HealMessage struct {
	Level      HealMessageLevel
	Text       string
	PackageKey string
}

// HealContext holds the per-dep state threaded through the heal chain.
type HealContext struct {
	PackageKey                       string
	LockfileMatch                    bool
	LockfileMatchViaContentHashOnly  bool
	UpdateRefs                       bool
	RefChanged                       bool
	BypassKeys                       map[string]bool
	FiredGroups                      map[string]bool
	Messages                         []HealMessage
}

// NewHealContext creates an initialised HealContext for one dependency.
func NewHealContext(
	packageKey string,
	lockfileMatch bool,
	lockfileMatchViaContentHashOnly bool,
	updateRefs bool,
	refChanged bool,
) HealContext {
	return HealContext{
		PackageKey:                      packageKey,
		LockfileMatch:                   lockfileMatch,
		LockfileMatchViaContentHashOnly: lockfileMatchViaContentHashOnly,
		UpdateRefs:                      updateRefs,
		RefChanged:                      refChanged,
		BypassKeys:                      map[string]bool{},
		FiredGroups:                     map[string]bool{},
	}
}

// AddWarn appends a WARN-level message to the context.
func (h *HealContext) AddWarn(text, packageKey string) {
	h.Messages = append(h.Messages, HealMessage{Level: HealMessageWarn, Text: text, PackageKey: packageKey})
}

// AddInfo appends an INFO-level message to the context.
func (h *HealContext) AddInfo(text, packageKey string) {
	h.Messages = append(h.Messages, HealMessage{Level: HealMessageInfo, Text: text, PackageKey: packageKey})
}

// Healer is implemented by each individual heal rule.
type Healer interface {
	// ExclusiveGroup returns a group tag; at most one healer per group fires
	// per dep. Empty string means no group.
	ExclusiveGroup() string
	// Applies returns true when this healer should run for the current context.
	Applies(hctx *HealContext) bool
	// Execute mutates hctx to apply the heal.
	Execute(hctx *HealContext)
}

// RunHealChain executes each healer in chain order, honouring exclusive groups.
// Returns the post-heal (lockfileMatch, refChanged) pair.
func RunHealChain(chain []Healer, hctx *HealContext) (lockfileMatch bool, refChanged bool) {
	for _, healer := range chain {
		if g := healer.ExclusiveGroup(); g != "" && hctx.FiredGroups[g] {
			continue
		}
		if !healer.Applies(hctx) {
			continue
		}
		healer.Execute(hctx)
		if g := healer.ExclusiveGroup(); g != "" {
			hctx.FiredGroups[g] = true
		}
	}
	return hctx.LockfileMatch, hctx.RefChanged
}
