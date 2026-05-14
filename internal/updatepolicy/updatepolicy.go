// Package updatepolicy provides build-time policy for APM self-update behavior.
// Package maintainers can patch constants during build to disable self-update
// and show users a package-manager-specific update command.
package updatepolicy

// DefaultSelfUpdateDisabledMessage is the default guidance when self-update is disabled.
const DefaultSelfUpdateDisabledMessage = "Self-update is disabled for this APM distribution. Update APM using your package manager."

// Build-time policy values. Packagers can override these at link time via
// -ldflags "-X updatepolicy.SelfUpdateEnabled=false".
var (
	// SelfUpdateEnabled controls whether self-update is allowed.
	SelfUpdateEnabled = true
	// SelfUpdateDisabledMessage is shown when self-update is disabled.
	SelfUpdateDisabledMessage = DefaultSelfUpdateDisabledMessage
)

// isPrintableASCII returns true when s contains only printable ASCII characters.
func isPrintableASCII(s string) bool {
	for _, c := range s {
		if c < ' ' || c > '~' {
			return false
		}
	}
	return true
}

// IsSelfUpdateEnabled returns true when this build allows self-update.
func IsSelfUpdateEnabled() bool {
	return SelfUpdateEnabled
}

// GetSelfUpdateDisabledMessage returns the guidance message shown when self-update is disabled.
func GetSelfUpdateDisabledMessage() string {
	if SelfUpdateDisabledMessage == "" {
		return DefaultSelfUpdateDisabledMessage
	}
	msg := SelfUpdateDisabledMessage
	if msg == "" {
		return DefaultSelfUpdateDisabledMessage
	}
	if !isPrintableASCII(msg) {
		return DefaultSelfUpdateDisabledMessage
	}
	return msg
}

// GetUpdateHintMessage returns the update hint used in startup notifications.
func GetUpdateHintMessage() string {
	if IsSelfUpdateEnabled() {
		return "Run apm update to upgrade"
	}
	return GetSelfUpdateDisabledMessage()
}
