// Package scriptrunner implements APM NPM-like script execution.
package scriptrunner

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// RuntimeKind identifies a supported AI runtime.
type RuntimeKind string

const (
	RuntimeCopilot RuntimeKind = "copilot"
	RuntimeCodex   RuntimeKind = "codex"
	RuntimeLLM     RuntimeKind = "llm"
	RuntimeGemini  RuntimeKind = "gemini"
	RuntimeUnknown RuntimeKind = "unknown"
)

// ScriptRunner executes APM scripts with auto-compilation of .prompt.md files.
type ScriptRunner struct {
	Compiler *PromptCompiler
	UseColor bool
}

// New returns a ScriptRunner with default settings.
func New(useColor bool) *ScriptRunner {
	return &ScriptRunner{
		Compiler: NewPromptCompiler(),
		UseColor: useColor,
	}
}

// RunScript runs a script from apm.yml with parameter substitution.
//
// Execution priority:
//  1. Explicit scripts in apm.yml
//  2. Auto-discovered prompt files
//  3. Error if not found
func (s *ScriptRunner) RunScript(scriptName string, params map[string]string) error {
	headerLines := formatScriptHeader(scriptName, params)
	for _, l := range headerLines {
		fmt.Println(l)
	}

	isVirtual := isVirtualPackageReference(scriptName)

	config, err := loadConfig()
	if err != nil || config == nil {
		if isVirtual {
			fmt.Println("  [i]  Creating minimal apm.yml for zero-config execution...")
			if createErr := createMinimalConfig(); createErr != nil {
				return createErr
			}
			config, err = loadConfig()
			if err != nil {
				return err
			}
		} else {
			return errors.New("No apm.yml found in current directory")
		}
	}

	// 1. Check explicit scripts first.
	if scripts, ok := config["scripts"].(map[string]any); ok {
		if cmdVal, found := scripts[scriptName]; found {
			if command, ok := cmdVal.(string); ok {
				return s.executeScriptCommand(command, params)
			}
		}
	}

	// 2. Auto-discover prompt file.
	discovered := s.discoverPromptFile(scriptName)
	if discovered != "" {
		fmt.Printf("[i] Auto-discovered: %s\n", filepath.ToSlash(discovered))
		rtKind, rtErr := detectInstalledRuntime()
		if rtErr != nil {
			return rtErr
		}
		command := generateRuntimeCommand(rtKind, discovered)
		return s.executeScriptCommand(command, params)
	}

	// 2.5 Try auto-install if it looks like a virtual package reference.
	if isVirtual {
		fmt.Printf("\n Auto-installing virtual package: %s\n", scriptName)
		if s.autoInstallVirtualPackage(scriptName) {
			discovered = s.discoverPromptFile(scriptName)
			if discovered != "" {
				fmt.Print("\n* Package installed and ready to run\n\n")
				rtKind, rtErr := detectInstalledRuntime()
				if rtErr != nil {
					return rtErr
				}
				command := generateRuntimeCommand(rtKind, discovered)
				return s.executeScriptCommand(command, params)
			}
			return errors.New("Package installed successfully but prompt not found.\n" +
				"The package may not contain the expected prompt file.\n" +
				"Check apm_modules for installed files.")
		}
	}

	// 3. Not found.
	var available string
	if scripts, ok := config["scripts"].(map[string]any); ok && len(scripts) > 0 {
		keys := make([]string, 0, len(scripts))
		for k := range scripts {
			keys = append(keys, k)
		}
		available = strings.Join(keys, ", ")
	} else {
		available = "none"
	}

	return fmt.Errorf(
		"Script or prompt '%s' not found.\n"+
			"Available scripts in apm.yml: %s\n\n"+
			"To find available prompts, check:\n"+
			"  - Local: .apm/prompts/, .github/prompts/, or project root\n"+
			"  - Dependencies: apm_modules/*/.apm/prompts/\n\n"+
			"Or install a prompt package:\n"+
			"  apm install <owner>/<repo>/path/to/prompt.prompt.md",
		scriptName, available,
	)
}

// executeScriptCommand executes a script command with parameter substitution.
func (s *ScriptRunner) executeScriptCommand(command string, params map[string]string) error {
	compiledCommand, compiledPromptFiles, runtimeContent := s.autoCompilePrompts(command, params)

	if len(compiledPromptFiles) > 0 {
		for _, line := range formatCompilationProgress(compiledPromptFiles) {
			fmt.Println(line)
		}
	}

	rtKind := detectRuntime(compiledCommand)

	if runtimeContent != "" {
		for _, line := range formatRuntimeExecution(rtKind, compiledCommand, len(runtimeContent)) {
			fmt.Println(line)
		}
		for _, line := range formatContentPreview(runtimeContent) {
			fmt.Println(line)
		}
	}

	env := setupRuntimeEnvironment()

	var envVarsSet []string
	if env["GITHUB_TOKEN"] != "" {
		envVarsSet = append(envVarsSet, "GITHUB_TOKEN")
	}
	if env["GITHUB_APM_PAT"] != "" {
		envVarsSet = append(envVarsSet, "GITHUB_APM_PAT")
	}
	if len(envVarsSet) > 0 {
		for _, line := range formatEnvironmentSetup(rtKind, envVarsSet) {
			fmt.Println(line)
		}
	}

	var cmdErr error
	if runtimeContent != "" {
		cmdErr = s.executeRuntimeCommand(compiledCommand, runtimeContent, env)
	} else {
		cmdErr = runShellCommand(compiledCommand, env)
	}

	if cmdErr != nil {
		for _, line := range formatExecutionError(rtKind) {
			fmt.Println(line)
		}
		var exitErr *exec.ExitError
		if errors.As(cmdErr, &exitErr) {
			return fmt.Errorf("Script execution failed with exit code %d", exitErr.ExitCode())
		}
		return fmt.Errorf("Script execution failed: %w", cmdErr)
	}

	for _, line := range formatExecutionSuccess(rtKind) {
		fmt.Println(line)
	}
	return nil
}

// ListScripts returns all available scripts from apm.yml.
func (s *ScriptRunner) ListScripts() map[string]string {
	config, err := loadConfig()
	if err != nil || config == nil {
		return nil
	}
	scripts, ok := config["scripts"].(map[string]any)
	if !ok {
		return nil
	}
	result := make(map[string]string, len(scripts))
	for k, v := range scripts {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

// autoCompilePrompts finds .prompt.md files in the command and compiles them.
// Returns (compiledCommand, compiledPromptFiles, runtimeContent).
func (s *ScriptRunner) autoCompilePrompts(command string, params map[string]string) (string, []string, string) {
	re := regexp.MustCompile(`(\S+\.prompt\.md)`)
	promptFiles := re.FindAllString(command, -1)

	var compiledPromptFiles []string
	var runtimeContent string
	compiledCommand := command

	runtimeCommands := []string{"copilot", "codex", "llm", "gemini"}

	for _, pf := range promptFiles {
		compiledPath, err := s.Compiler.Compile(pf, params)
		if err != nil {
			continue
		}
		compiledPromptFiles = append(compiledPromptFiles, pf)

		data, err := os.ReadFile(compiledPath)
		if err != nil {
			continue
		}
		compiledContent := strings.TrimSpace(string(data))

		// Check if this is a runtime command.
		isRuntimeCmd := false
		for _, rt := range runtimeCommands {
			re2 := regexp.MustCompile(`(?:^|\s)` + rt + `(?:\s|$)`)
			if re2.MatchString(command) && strings.Contains(command, pf) {
				isRuntimeCmd = true
				break
			}
		}

		compiledCommand = transformRuntimeCommand(compiledCommand, pf, compiledContent, compiledPath)

		if isRuntimeCmd {
			runtimeContent = compiledContent
		}
	}

	return compiledCommand, compiledPromptFiles, runtimeContent
}

// transformRuntimeCommand rewrites a command containing a .prompt.md reference
// to use the appropriate runtime invocation.
func transformRuntimeCommand(command, promptFile, compiledContent, compiledPath string) string {
	runtimeCommands := []string{"codex", "copilot", "llm", "gemini"}

	// Try env-var prefix pattern first.
	for _, rt := range runtimeCommands {
		rtPattern := " " + rt + " "
		if strings.Contains(command, rtPattern) && strings.Contains(command, promptFile) {
			parts := strings.SplitN(command, rtPattern, 2)
			potentialEnvPart := parts[0]
			runtimePart := rt + " " + parts[1]

			if strings.Contains(potentialEnvPart, "=") && !strings.HasPrefix(potentialEnvPart, rt) {
				result := parseAndBuildRuntimeCommand(rt, runtimePart, promptFile, potentialEnvPart)
				if result != "" {
					return result
				}
			}
		}
	}

	// Try individual runtime patterns without env-var prefix.
	for _, rt := range runtimeCommands {
		re := regexp.MustCompile(`^` + rt + `\s+.*` + regexp.QuoteMeta(promptFile))
		if re.MatchString(command) {
			result := parseAndBuildRuntimeCommand(rt, command, promptFile, "")
			if result != "" {
				return result
			}
		}
	}

	// Bare prompt file -> codex exec.
	if strings.TrimSpace(command) == promptFile {
		return "codex exec"
	}

	// Fallback: replace file path with compiled path.
	return strings.ReplaceAll(command, promptFile, compiledPath)
}

func parseAndBuildRuntimeCommand(rtCmd, commandPart, promptFile, envPrefix string) string {
	pattern := regexp.MustCompile(rtCmd + `\s+(.*?)(` + regexp.QuoteMeta(promptFile) + `)(.*?)$`)
	m := pattern.FindStringSubmatch(commandPart)
	if m == nil {
		return ""
	}
	argsBefore := strings.TrimSpace(m[1])
	argsAfter := strings.TrimSpace(m[3])

	if envPrefix != "" && rtCmd != "codex" {
		argsBefore = strings.TrimSpace(strings.ReplaceAll(argsBefore, "-p", ""))
	}

	prefix := ""
	if envPrefix != "" {
		prefix = envPrefix + " "
	}

	switch rtCmd {
	case "codex":
		result := prefix + "codex exec"
		if argsBefore != "" {
			result += " " + argsBefore
		}
		if argsAfter != "" {
			result += " " + argsAfter
		}
		return result
	case "copilot":
		cleaned := strings.TrimSpace(strings.ReplaceAll(argsBefore, "-p", ""))
		result := prefix + "copilot"
		if cleaned != "" {
			result += " " + cleaned
		}
		if argsAfter != "" {
			result += " " + argsAfter
		}
		return result
	case "llm":
		result := prefix + "llm"
		if argsBefore != "" {
			result += " " + argsBefore
		}
		if argsAfter != "" {
			result += " " + argsAfter
		}
		return result
	case "gemini":
		re := regexp.MustCompile(`(^|\s)-p(\s|$)`)
		cleaned := strings.TrimSpace(re.ReplaceAllString(argsBefore, "$1$2"))
		result := prefix + "gemini"
		if cleaned != "" {
			result += " " + cleaned
		}
		if argsAfter != "" {
			result += " " + argsAfter
		}
		return result
	}
	return ""
}

// detectRuntime detects which runtime is referenced in a command.
func detectRuntime(command string) RuntimeKind {
	lower := strings.ToLower(strings.TrimSpace(command))
	patterns := []struct {
		rt  RuntimeKind
		pat string
	}{
		{RuntimeCopilot, `(?:^|\s)copilot(?:\s|$)`},
		{RuntimeCodex, `(?:^|\s)codex(?:\s|$)`},
		{RuntimeLLM, `(?:^|\s)llm(?:\s|$)`},
		{RuntimeGemini, `(?:^|\s)gemini(?:\s|$)`},
	}
	for _, p := range patterns {
		if matched, _ := regexp.MatchString(p.pat, lower); matched {
			return p.rt
		}
	}
	return RuntimeUnknown
}

// executeRuntimeCommand runs a runtime command passing content as an argument.
func (s *ScriptRunner) executeRuntimeCommand(command, content string, env map[string]string) error {
	args := splitArgs(command)

	// Extract env-var prefixes from the front of args.
	envVars := copyEnv(env)
	var actualArgs []string
	for _, arg := range args {
		if strings.Contains(arg, "=") && len(actualArgs) == 0 {
			kv := strings.SplitN(arg, "=", 2)
			if isValidEnvVarName(kv[0]) {
				envVars[kv[0]] = kv[1]
				continue
			}
		}
		actualArgs = append(actualArgs, arg)
	}

	rtKind := detectRuntime(strings.Join(actualArgs, " "))
	switch rtKind {
	case RuntimeCopilot:
		actualArgs = append(actualArgs, "-p", content)
	case RuntimeCodex:
		actualArgs = append(actualArgs, content)
	case RuntimeLLM:
		actualArgs = append(actualArgs, content)
	case RuntimeGemini:
		actualArgs = append(actualArgs, "-p", content)
	default:
		actualArgs = append(actualArgs, content)
	}

	// On Windows, resolve via PATH to find .cmd / .ps1 wrappers.
	if len(actualArgs) > 0 && runtime.GOOS == "windows" {
		if resolved, err := exec.LookPath(actualArgs[0]); err == nil {
			actualArgs[0] = resolved
		}
	}

	cmd := exec.Command(actualArgs[0], actualArgs[1:]...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envMapToSlice(envVars)
	return cmd.Run()
}

// runShellCommand executes a command via the system shell.
func runShellCommand(command string, env map[string]string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) //nolint:gosec
	} else {
		cmd = exec.Command("sh", "-c", command) //nolint:gosec
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envMapToSlice(env)
	return cmd.Run()
}

// discoverPromptFile discovers a prompt file by name.
func (s *ScriptRunner) discoverPromptFile(name string) string {
	if strings.Contains(name, "/") {
		return s.discoverQualifiedPrompt(name)
	}

	searchName := name
	if !strings.HasSuffix(searchName, ".prompt.md") {
		searchName = name + ".prompt.md"
	}

	// Local search paths.
	localPaths := []string{
		searchName,
		filepath.Join(".apm", "prompts", searchName),
		filepath.Join(".github", "prompts", searchName),
	}
	for _, p := range localPaths {
		fi, err := os.Lstat(p)
		if err == nil && !fi.IsDir() && fi.Mode()&fs.ModeSymlink == 0 {
			return p
		}
	}

	// Search in apm_modules.
	apmModules := "apm_modules"
	if _, err := os.Stat(apmModules); err != nil {
		return ""
	}

	var matches []string
	_ = filepath.WalkDir(apmModules, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.Name() == searchName {
			matches = append(matches, path)
		}
		// Also look for SKILL.md in a directory matching `name`.
		if d.IsDir() && d.Name() == name {
			skillFile := filepath.Join(path, "SKILL.md")
			if fi, err2 := os.Lstat(skillFile); err2 == nil && !fi.IsDir() {
				matches = append(matches, skillFile)
			}
		}
		return nil
	})

	if len(matches) == 1 {
		return matches[0]
	}
	if len(matches) > 1 {
		// Collision — build error message and print it; callers check empty string.
		fmt.Fprint(os.Stderr, buildCollisionError(name, matches))
		return ""
	}
	return ""
}

func buildCollisionError(name string, matches []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Multiple prompts found for '%s':\n", name)
	for _, m := range matches {
		parts := strings.Split(filepath.ToSlash(m), "/")
		idx := -1
		for i, p := range parts {
			if p == "apm_modules" {
				idx = i
				break
			}
		}
		if idx >= 0 && idx+2 < len(parts) {
			fmt.Fprintf(&b, "  - %s/%s (%s)\n", parts[idx+1], parts[idx+2], m)
		} else {
			fmt.Fprintf(&b, "  - %s\n", m)
		}
	}
	fmt.Fprintln(&b, "\nPlease specify using qualified path:")
	for _, m := range matches {
		parts := strings.Split(filepath.ToSlash(m), "/")
		idx := -1
		for i, p := range parts {
			if p == "apm_modules" {
				idx = i
				break
			}
		}
		if idx >= 0 && idx+2 < len(parts) {
			fmt.Fprintf(&b, "  apm run %s/%s/%s\n", parts[idx+1], parts[idx+2], name)
		}
	}
	fmt.Fprintln(&b, "\nOr add an explicit script to apm.yml:")
	fmt.Fprintln(&b, "  scripts:")
	fmt.Fprintf(&b, "    my-%s: \"copilot -p <path-to-preferred-prompt>\"\n", name)
	return b.String()
}

// discoverQualifiedPrompt discovers a prompt using owner/repo/name format.
func (s *ScriptRunner) discoverQualifiedPrompt(qualifiedPath string) string {
	parts := strings.Split(qualifiedPath, "/")
	if len(parts) < 2 {
		return ""
	}

	promptName := parts[len(parts)-1]
	if !strings.HasSuffix(promptName, ".prompt.md") {
		promptName = promptName + ".prompt.md"
	}

	apmModules := "apm_modules"
	if _, err := os.Stat(apmModules); err != nil {
		return ""
	}

	// For 3+ part qualified paths, check subdirectory SKILL.md first.
	if len(parts) >= 3 {
		subdirPath := filepath.Join(append([]string{apmModules}, parts...)...)
		skillFile := filepath.Join(subdirPath, "SKILL.md")
		if fi, err := os.Lstat(skillFile); err == nil && !fi.IsDir() {
			return skillFile
		}
	}

	owner := parts[0]
	ownerDir := filepath.Join(apmModules, owner)
	if _, err := os.Stat(ownerDir); err != nil {
		return ""
	}

	entries, err := os.ReadDir(ownerDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgDir := filepath.Join(ownerDir, entry.Name())
		var found string
		_ = filepath.WalkDir(pkgDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || found != "" {
				return nil
			}
			if d.Name() == promptName {
				// Check qualified path match.
				pathSlash := filepath.ToSlash(path)
				qParts := strings.Split(qualifiedPath, "/")
				if qParts[0] != "" && strings.Contains(pathSlash, qParts[0]) {
					expectedName := qParts[len(qParts)-1]
					if !strings.HasSuffix(expectedName, ".prompt.md") {
						expectedName += ".prompt.md"
					}
					if d.Name() == expectedName {
						found = path
					}
				}
			}
			return nil
		})
		if found != "" {
			return found
		}
	}
	return ""
}

// isVirtualPackageReference returns true if name looks like owner/repo/... syntax.
func isVirtualPackageReference(name string) bool {
	return strings.Count(name, "/") >= 2
}

// autoInstallVirtualPackage is a stub — actual install requires network access.
func (s *ScriptRunner) autoInstallVirtualPackage(packageRef string) bool {
	fmt.Printf("  [x] Auto-install not supported in Go runtime: %s\n", packageRef)
	return false
}

// detectInstalledRuntime detects an installed AI runtime CLI.
func detectInstalledRuntime() (RuntimeKind, error) {
	for _, rt := range []struct {
		name RuntimeKind
		bin  string
	}{
		{RuntimeCopilot, "copilot"},
		{RuntimeCodex, "codex"},
		{RuntimeGemini, "gemini"},
	} {
		if _, err := exec.LookPath(rt.bin); err == nil {
			return rt.name, nil
		}
	}
	return RuntimeUnknown, errors.New("No compatible runtime found.\n" +
		"Install GitHub Copilot CLI with:\n" +
		"  apm runtime setup copilot\n" +
		"Or install Codex CLI with:\n" +
		"  apm runtime setup codex\n" +
		"Or install Gemini CLI with:\n" +
		"  apm runtime setup gemini")
}

// generateRuntimeCommand generates a default runtime invocation for a discovered prompt.
func generateRuntimeCommand(rt RuntimeKind, promptFile string) string {
	switch rt {
	case RuntimeCopilot:
		return fmt.Sprintf("copilot --log-level all --log-dir copilot-logs --allow-all-tools -p %s", promptFile)
	case RuntimeCodex:
		return fmt.Sprintf("codex -s workspace-write --skip-git-repo-check %s", promptFile)
	case RuntimeGemini:
		return fmt.Sprintf("gemini -p %s", promptFile)
	default:
		return fmt.Sprintf("copilot -p %s", promptFile)
	}
}

// setupRuntimeEnvironment builds the environment map for script execution.
func setupRuntimeEnvironment() map[string]string {
	env := make(map[string]string)
	for _, kv := range os.Environ() {
		idx := strings.IndexByte(kv, '=')
		if idx >= 0 {
			env[kv[:idx]] = kv[idx+1:]
		}
	}
	return env
}

// loadConfig loads apm.yml from the current directory using a minimal YAML parser.
func loadConfig() (map[string]any, error) {
	data, err := os.ReadFile("apm.yml")
	if err != nil {
		return nil, err
	}
	return parseSimpleYAML(string(data)), nil
}

// parseSimpleYAML is a minimal single-level YAML parser sufficient for apm.yml.
func parseSimpleYAML(content string) map[string]any {
	result := make(map[string]any)
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentKey string
	var currentList []any
	var currentMap map[string]any
	inMap := false
	inList := false

	flush := func() {
		if currentKey == "" {
			return
		}
		if inMap && currentMap != nil {
			result[currentKey] = currentMap
		} else if inList && currentList != nil {
			result[currentKey] = currentList
		}
		currentKey = ""
		currentMap = nil
		currentList = nil
		inMap = false
		inList = false
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Top-level key: value pair
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			flush()
			currentKey = key

			if val == "" {
				// Could be start of map or list — wait for next line
				inMap = false
				inList = false
			} else {
				result[key] = unquoteYAML(val)
				currentKey = ""
			}
			continue
		}

		// Indented list item
		if strings.HasPrefix(strings.TrimLeft(line, " \t"), "- ") && currentKey != "" {
			item := strings.TrimSpace(strings.TrimLeft(line, " \t")[2:])
			if !inList {
				flush()
				currentKey = strings.Split(line, ":")[0] // recover key — but we lost it
			}
			inList = true
			currentList = append(currentList, unquoteYAML(item))
			continue
		}

		// Indented key: value (sub-map)
		trimmed := strings.TrimLeft(line, " \t")
		if strings.Contains(trimmed, ":") && currentKey != "" {
			parts := strings.SplitN(trimmed, ":", 2)
			subKey := strings.TrimSpace(parts[0])
			subVal := strings.TrimSpace(parts[1])
			if !inMap {
				currentMap = make(map[string]any)
				inMap = true
			}
			currentMap[subKey] = unquoteYAML(subVal)
			continue
		}
	}
	flush()
	return result
}

func unquoteYAML(s string) string {
	if len(s) >= 2 &&
		((s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

// createMinimalConfig creates a minimal apm.yml for zero-config usage.
func createMinimalConfig() error {
	cwd, _ := os.Getwd()
	name := filepath.Base(cwd)
	content := fmt.Sprintf("name: %s\nversion: 1.0.0\ndescription: Auto-generated for zero-config virtual package execution\n", name)
	return os.WriteFile("apm.yml", []byte(content), 0o644)
}

// -- Helpers -----------------------------------------------------------------

func splitArgs(command string) []string {
	// Simple POSIX-style tokenizer: handle quoted strings.
	var args []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(command); i++ {
		c := command[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case c == ' ' && !inSingle && !inDouble:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func isValidEnvVarName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, c := range s {
		if i == 0 && !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_') {
			return false
		}
		if i > 0 && !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_') {
			return false
		}
	}
	return true
}

func copyEnv(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func envMapToSlice(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

// -- Formatter stubs (plain-text, no Rich dependency) -----------------------

func formatScriptHeader(scriptName string, params map[string]string) []string {
	lines := []string{fmt.Sprintf("[*] Running script: %s", scriptName)}
	if len(params) > 0 {
		parts := make([]string, 0, len(params))
		for k, v := range params {
			parts = append(parts, k+"="+v)
		}
		lines = append(lines, "    Parameters: "+strings.Join(parts, ", "))
	}
	return lines
}

func formatCompilationProgress(files []string) []string {
	return []string{fmt.Sprintf("[*] Compiled: %s", strings.Join(files, ", "))}
}

func formatRuntimeExecution(rt RuntimeKind, command string, contentLen int) []string {
	return []string{fmt.Sprintf("[>] Executing via %s (%d bytes)", rt, contentLen)}
}

func formatContentPreview(content string) []string {
	preview := content
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	return []string{"    " + strings.ReplaceAll(preview, "\n", "\n    ")}
}

func formatEnvironmentSetup(rt RuntimeKind, vars []string) []string {
	return []string{fmt.Sprintf("[i] Environment: %s", strings.Join(vars, ", "))}
}

func formatExecutionSuccess(rt RuntimeKind) []string {
	return []string{fmt.Sprintf("[+] Script completed successfully via %s", rt)}
}

func formatExecutionError(rt RuntimeKind) []string {
	return []string{fmt.Sprintf("[x] Script failed via %s", rt)}
}
