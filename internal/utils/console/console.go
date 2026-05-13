// Package console provides console utility functions for formatted CLI output.
//
// All output is within printable ASCII (U+0020-U+007E). Color codes use ANSI
// escape sequences, disabled automatically when NO_COLOR is set or TERM=dumb.
package console

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// StatusSymbols maps semantic names to ASCII bracket notation.
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

// ANSI color codes.
const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
	ansiBold   = "\033[1m"
)

// colorEnabled returns true when ANSI color output is supported.
func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return true
}

// Echo writes a message to w (defaults to os.Stdout) with optional color and
// symbol prefix. color may be "red", "green", "yellow", "blue", "cyan", or
// empty for default terminal color.
func Echo(w io.Writer, message, color, symbol string, bold bool) {
	if w == nil {
		w = os.Stdout
	}
	if sym, ok := StatusSymbols[symbol]; ok && symbol != "" {
		message = sym + " " + message
	}
	if colorEnabled() && color != "" {
		code := colorCode(color)
		if bold {
			fmt.Fprintf(w, "%s%s%s%s\n", ansiBold, code, message, ansiReset)
		} else {
			fmt.Fprintf(w, "%s%s%s\n", code, message, ansiReset)
		}
	} else {
		fmt.Fprintln(w, message)
	}
}

func colorCode(color string) string {
	switch strings.ToLower(color) {
	case "red":
		return ansiRed
	case "green":
		return ansiGreen
	case "yellow":
		return ansiYellow
	case "blue":
		return ansiBlue
	case "cyan":
		return ansiCyan
	default:
		return ""
	}
}

// Success prints a success message (green, bold).
func Success(message, symbol string) {
	Echo(os.Stdout, message, "green", symbol, true)
}

// Error prints an error message (red).
func Error(message, symbol string) {
	Echo(os.Stderr, message, "red", symbol, false)
}

// Warning prints a warning message (yellow).
func Warning(message, symbol string) {
	Echo(os.Stdout, message, "yellow", symbol, false)
}

// Info prints an info message (blue).
func Info(message, symbol string) {
	Echo(os.Stdout, message, "blue", symbol, false)
}

// Panel prints content framed by a simple ASCII border with an optional title.
func Panel(content, title, style string) {
	if title != "" {
		fmt.Printf("\n--- %s ---\n", title)
	}
	fmt.Println(content)
	if title != "" {
		fmt.Println(strings.Repeat("-", len(title)+8))
	}
}

// PrintFilesTable prints a simple two-column table of file name + description.
func PrintFilesTable(files [][]string, tableTitle string) {
	if tableTitle != "" {
		fmt.Println(tableTitle)
	}
	for _, row := range files {
		name := ""
		desc := ""
		if len(row) > 0 {
			name = row[0]
		}
		if len(row) > 1 {
			desc = row[1]
		}
		fmt.Printf("  %-40s %s\n", name, desc)
	}
}

// DownloadSpinner prints a simple download-in-progress message and calls fn.
// Unlike Python's context-manager spinner, this is a function-based helper.
func DownloadSpinner(repoName string, fn func()) {
	fmt.Printf("[>] Downloading %s...\n", repoName)
	fn()
}
