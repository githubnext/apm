// Package console provides utility functions for formatting and CLI output.
//
// Provides STATUS_SYMBOLS and rich-echo helpers with ANSI colour output.
// All output stays within printable ASCII (U+0020-U+007E).
package console

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// StatusSymbols maps symbol names to ASCII bracket notation.
var StatusSymbols = map[string]string{
	"success":  "[*]",
	"sparkles": "[*]",
	"running":  "[>]",
	"gear":     "[*]",
	"info":     "[i]",
	"warning":  "[!]",
	"error":    "[x]",
	"check":    "[+]",
	"cross":    "[x]",
	"list":     "[#]",
	"preview":  "[>]",
	"robot":    "[>]",
	"metrics":  "[#]",
	"default":  "[>]",
	"eyes":     "[>]",
	"folder":   "[>]",
	"cogs":     "[*]",
	"plugin":   "[>]",
	"search":   "[>]",
	"download": "[>]",
	"update":   "[~]",
	"remove":   "[-]",
	"equal":    "[=]",
}

// ANSI colour codes (ASCII-safe; guarded by NO_COLOR / TERM=dumb).
var ansiColors = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"muted":   "\033[37m",
	"info":    "\033[34m",
	"reset":   "\033[0m",
	"bold":    "\033[1m",
}

var (
	colorEnabled     bool
	colorEnabledOnce sync.Once
)

func isColorEnabled() bool {
	colorEnabledOnce.Do(func() {
		if os.Getenv("NO_COLOR") != "" {
			colorEnabled = false
			return
		}
		term := os.Getenv("TERM")
		if term == "dumb" {
			colorEnabled = false
			return
		}
		colorEnabled = true
	})
	return colorEnabled
}

// DefaultOut is the output writer (defaults to os.Stdout).
var DefaultOut io.Writer = os.Stdout

// RichEcho writes a coloured message to DefaultOut.
func RichEcho(message, color string, bold bool, symbol string) {
	if sym, ok := StatusSymbols[symbol]; ok {
		message = sym + " " + message
	}
	if isColorEnabled() {
		prefix := ""
		if c, ok := ansiColors[color]; ok {
			prefix += c
		}
		if bold {
			prefix += ansiColors["bold"]
		}
		if prefix != "" {
			fmt.Fprintf(DefaultOut, "%s%s%s\n", prefix, message, ansiColors["reset"])
			return
		}
	}
	fmt.Fprintln(DefaultOut, message)
}

// RichSuccess displays a success message (green, bold).
func RichSuccess(message, symbol string) {
	RichEcho(message, "green", true, symbol)
}

// RichError displays an error message (red).
func RichError(message, symbol string) {
	RichEcho(message, "red", false, symbol)
}

// RichWarning displays a warning message (yellow).
func RichWarning(message, symbol string) {
	RichEcho(message, "yellow", false, symbol)
}

// RichInfo displays an info message (blue).
func RichInfo(message, symbol string) {
	RichEcho(message, "blue", false, symbol)
}

// RichPanel displays content in a simple text panel.
func RichPanel(content, title, style string) {
	if title != "" {
		fmt.Fprintf(DefaultOut, "\n--- %s ---\n", title)
	}
	fmt.Fprintln(DefaultOut, content)
	if title != "" {
		border := ""
		for i := 0; i < len(title)+8; i++ {
			border += "-"
		}
		fmt.Fprintln(DefaultOut, border)
	}
}
