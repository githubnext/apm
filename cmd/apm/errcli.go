// errcli.go rewrites 2-line Go unknown-option error messages into the
// 4-line Click 8.x format that Python emits, ensuring output parity for
// test_every_python_command_rejects_unknown_option_consistently.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// cmdUsageSuffix maps "apm CMD" paths to the suffix after the command path
// in the Usage line.  Commands not listed use the default " [OPTIONS]".
var cmdUsageSuffix = map[string]string{
	"apm":                            " [OPTIONS] COMMAND [ARGS]...",
	"apm cache":                      " [OPTIONS] COMMAND [ARGS]...",
	"apm config":                     " [OPTIONS] [COMMAND] [ARGS]...",
	"apm deps":                       " [OPTIONS] COMMAND [ARGS]...",
	"apm experimental":               " [OPTIONS] [COMMAND] [ARGS]...",
	"apm marketplace":                " [OPTIONS] COMMAND [ARGS]...",
	"apm marketplace package":        " [OPTIONS] COMMAND [ARGS]...",
	"apm mcp":                        " [OPTIONS] COMMAND [ARGS]...",
	"apm plugin":                     " [OPTIONS] COMMAND [ARGS]...",
	"apm plugin init":                " [OPTIONS] [PROJECT_NAME]",
	"apm policy":                     " [OPTIONS] COMMAND [ARGS]...",
	"apm runtime":                    " [OPTIONS] COMMAND [ARGS]...",
	"apm targets":                    " [OPTIONS] [COMMAND] [ARGS]...",
	"apm audit":                      " [OPTIONS] [PACKAGE]",
	"apm config get":                 " [OPTIONS] [KEY]",
	"apm config set":                 " [OPTIONS] KEY VALUE",
	"apm config unset":               " [OPTIONS] KEY",
	"apm deps info":                  " [OPTIONS] PACKAGE",
	"apm deps update":                " [OPTIONS] [PACKAGES]...",
	"apm experimental disable":       " [OPTIONS] NAME",
	"apm experimental enable":        " [OPTIONS] NAME",
	"apm experimental reset":         " [OPTIONS] [NAME]",
	"apm init":                       " [OPTIONS] [PROJECT_NAME]",
	"apm info":                       " [OPTIONS] PACKAGE [FIELD]",
	"apm install":                    " [OPTIONS] [PACKAGES]...",
	"apm marketplace add":            " [OPTIONS] REPO",
	"apm marketplace browse":         " [OPTIONS] NAME",
	"apm marketplace package add":    " [OPTIONS] SOURCE",
	"apm marketplace package remove": " [OPTIONS] NAME",
	"apm marketplace package set":    " [OPTIONS] NAME",
	"apm marketplace remove":         " [OPTIONS] NAME",
	"apm marketplace update":         " [OPTIONS] [NAME]",
	"apm marketplace validate":       " [OPTIONS] NAME",
	"apm mcp install":                " [OPTIONS] NAME",
	"apm mcp search":                 " [OPTIONS] QUERY",
	"apm mcp show":                   " [OPTIONS] SERVER_NAME",
	"apm preview":                    " [OPTIONS] [SCRIPT_NAME]",
	"apm run":                        " [OPTIONS] [SCRIPT_NAME]",
	"apm runtime remove":             " [OPTIONS] {copilot|codex|llm|gemini}",
	"apm runtime setup":              " [OPTIONS] {copilot|codex|llm|gemini}",
	"apm search":                     " [OPTIONS] QUERY@MARKETPLACE",
	"apm uninstall":                  " [OPTIONS] PACKAGES...",
	"apm unpack":                     " [OPTIONS] BUNDLE_PATH",
	"apm view":                       " [OPTIONS] PACKAGE [FIELD]",
}

// usageLine returns the full Usage line for a command path.
func usageLine(cmdPath string) string {
	if suf, ok := cmdUsageSuffix[cmdPath]; ok {
		return "Usage: " + cmdPath + suf
	}
	return "Usage: " + cmdPath + " [OPTIONS]"
}

// clickErrWriter intercepts stderr writes and converts the 2-line Go error
// format ("Error: No such option: X\n", "Try '...' for help.\n") into the
// 4-line Click 8.x format: usage, try, blank, error.
type clickErrWriter struct {
	w       io.Writer
	pending string // "Error: No such option: X\n" waiting for its Try line
	lineBuf string // incomplete line waiting for its terminating \n
}

func (w *clickErrWriter) processLine(line string) {
	if w.pending != "" {
		if strings.HasPrefix(line, "Try '") {
			tryContent := strings.TrimPrefix(line, "Try '")
			tryContent = strings.TrimSuffix(strings.TrimRight(tryContent, "\n"), "' for help.")
			cmdPath := strings.TrimSuffix(tryContent, " --help")
			// Convert Go "Error: No such option: --X" to Click 8.x "Error: No such option '--X'."
			errLine := w.pending
			const errPrefix = "Error: No such option: "
			if strings.HasPrefix(errLine, errPrefix) {
				opt := strings.TrimRight(strings.TrimPrefix(errLine, errPrefix), "\n")
				errLine = "Error: No such option '" + opt + "'.\n"
			}
			fmt.Fprintf(w.w, "%s\n%s\n%s", usageLine(cmdPath), line, errLine)
			w.pending = ""
			return
		}
		fmt.Fprint(w.w, w.pending)
		w.pending = ""
	}
	if strings.HasPrefix(line, "Error: No such option: ") {
		w.pending = line
		return
	}
	fmt.Fprint(w.w, line)
}

func (w *clickErrWriter) flush() {
	if w.lineBuf != "" {
		fmt.Fprint(w.w, w.lineBuf)
		w.lineBuf = ""
	}
	if w.pending != "" {
		fmt.Fprint(w.w, w.pending)
		w.pending = ""
	}
}

func (w *clickErrWriter) Write(p []byte) (int, error) {
	s := w.lineBuf + string(p)
	w.lineBuf = ""
	for {
		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			w.lineBuf = s
			break
		}
		w.processLine(s[:idx+1])
		s = s[idx+1:]
	}
	return len(p), nil
}

// wrapStderr replaces os.Stderr with a pipe whose reader runs a
// clickErrWriter goroutine.  The returned flush function must be called
// (and allowed to return) before os.Exit to drain the pipe.
func wrapStderr() func() {
	r, w, err := os.Pipe()
	if err != nil {
		return func() {}
	}
	orig := os.Stderr
	os.Stderr = w
	done := make(chan struct{})
	go func() {
		defer close(done)
		ew := &clickErrWriter{w: orig}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			ew.processLine(scanner.Text() + "\n")
		}
		ew.flush()
		r.Close()
	}()
	return func() {
		w.Close()
		<-done
	}
}
