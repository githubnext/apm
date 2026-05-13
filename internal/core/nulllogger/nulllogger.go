// Package nulllogger provides a console-fallback logger for integrator contexts.
package nulllogger

// NullCommandLogger is a partial CommandLogger facade for MCPIntegrator contexts.
// Every implemented method produces visible terminal output via fmt.Print.
type NullCommandLogger struct {
Verbose bool
}

// Start logs a start message.
func (l *NullCommandLogger) Start(message, symbol string) {
if symbol == "" {
symbol = "running"
}
log("[i]", message)
}

// Progress logs a progress message.
func (l *NullCommandLogger) Progress(message, symbol string) {
log("[i]", message)
}

// Success logs a success message.
func (l *NullCommandLogger) Success(message, symbol string) {
log("[+]", message)
}

// Warning logs a warning message.
func (l *NullCommandLogger) Warning(message, symbol string) {
log("[!]", message)
}

// Error logs an error message.
func (l *NullCommandLogger) Error(message, symbol string) {
log("[x]", message)
}

// VerboseDetail discards verbose details (Verbose is always false).
func (l *NullCommandLogger) VerboseDetail(message string) {}

// TreeItem logs a tree item.
func (l *NullCommandLogger) TreeItem(message string) {
log("  -", message)
}

// PackageInlineWarning discards inline warnings.
func (l *NullCommandLogger) PackageInlineWarning(message string) {}

// MCPLookupHeartbeat mirrors CommandLogger.MCPLookupHeartbeat.
func (l *NullCommandLogger) MCPLookupHeartbeat(count int) {
if count <= 0 {
return
}
noun := "servers"
if count == 1 {
noun = "server"
}
log("[>]", "Looking up "+itoa(count)+" MCP "+noun+" in registry...")
}

func log(symbol, msg string) {
println(symbol + " " + msg)
}

func itoa(n int) string {
if n < 0 {
return "-" + itoa(-n)
}
if n < 10 {
return string(rune('0' + n))
}
return itoa(n/10) + string(rune('0'+n%10))
}
