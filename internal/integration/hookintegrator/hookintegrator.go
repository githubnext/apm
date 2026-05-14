// Package hookintegrator provides hook integration for APM packages.
// Deploys hook JSON files and referenced scripts to target directories.
// Ported from src/apm_cli/integration/hook_integrator.py
package hookintegrator

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
	"github.com/githubnext/apm/internal/integration/targets"
)

// HookIntegrationResult holds results of a hook integration operation.
type HookIntegrationResult struct {
	FilesIntegrated int
	FilesUpdated    int
	FilesSkipped    int
	TargetPaths     []string
	ScriptsCopied   int
}

// HooksIntegrated is an alias for FilesIntegrated (backward compat).
func (r *HookIntegrationResult) HooksIntegrated() int {
	return r.FilesIntegrated
}

// mergeHookConfig describes a target that merges hooks into a single JSON file.
type mergeHookConfig struct {
	ConfigFilename string
	TargetKey      string
	RequireDir     bool
}

// hookEventMap maps source event names to target-specific names.
var hookEventMap = map[string]map[string]string{
	"claude": {
		"preToolUse":  "PreToolUse",
		"postToolUse": "PostToolUse",
	},
	"gemini": {
		"PreToolUse":  "BeforeTool",
		"preToolUse":  "BeforeTool",
		"PostToolUse": "AfterTool",
		"postToolUse": "AfterTool",
		"Stop":        "SessionEnd",
	},
}

// mergeHookTargets maps target names to merge configurations.
var mergeHookTargets = map[string]mergeHookConfig{
	"claude":   {ConfigFilename: "settings.json", TargetKey: "claude", RequireDir: false},
	"cursor":   {ConfigFilename: "hooks.json", TargetKey: "cursor", RequireDir: true},
	"codex":    {ConfigFilename: "hooks.json", TargetKey: "codex", RequireDir: true},
	"gemini":   {ConfigFilename: "settings.json", TargetKey: "gemini", RequireDir: true},
	"windsurf": {ConfigFilename: "hooks.json", TargetKey: "windsurf", RequireDir: true},
}

// hookFileTargetSuffixes maps hook file stem suffixes to target sets.
var hookFileTargetSuffixes = map[string]map[string]bool{
	"copilot-hooks":  {"copilot": true, "vscode": true},
	"cursor-hooks":   {"cursor": true},
	"claude-hooks":   {"claude": true},
	"codex-hooks":    {"codex": true},
	"gemini-hooks":   {"gemini": true},
	"windsurf-hooks": {"windsurf": true},
}

// hookCommandKeys lists all supported hook command keys.
var hookCommandKeys = []string{"command", "bash", "powershell", "windows", "linux", "osx"}

// pluginRootRe matches ${CLAUDE_PLUGIN_ROOT}/path and similar.
var pluginRootRe = regexp.MustCompile(`\$\{(?:CLAUDE_PLUGIN_ROOT|CURSOR_PLUGIN_ROOT|PLUGIN_ROOT)\}([\\/][^\s]+)`)

// relPathRe matches relative ./path or .\path references.
var relPathRe = regexp.MustCompile(`(\.[\\/][^\s]+)`)

// filterHookFilesForTarget returns only hook files intended for targetKey.
func filterHookFilesForTarget(hookFiles []string, targetKey string) []string {
	var result []string
	for _, hf := range hookFiles {
		stemLower := strings.ToLower(strings.TrimSuffix(filepath.Base(hf), filepath.Ext(hf)))
		matchedSuffix := ""
		matched := false
		for suffix, allowed := range hookFileTargetSuffixes {
			if stemLower == suffix || strings.HasSuffix(stemLower, "-"+suffix) {
				matchedSuffix = suffix
				if allowed[targetKey] {
					result = append(result, hf)
					matched = true
				}
				break
			}
		}
		if matchedSuffix == "" && !matched {
			// Universal -- deploy to all targets
			result = append(result, hf)
		}
	}
	return result
}

// toGeminiHookEntries transforms hook entries to Gemini CLI format.
func toGeminiHookEntries(entries []interface{}) []interface{} {
	var result []interface{}
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			result = append(result, raw)
			continue
		}
		// Already nested (Claude/Gemini format)
		if hooks, ok := entry["hooks"].([]interface{}); ok {
			for _, h := range hooks {
				if hm, ok := h.(map[string]interface{}); ok {
					copilotKeysToGemini(hm)
				}
			}
			result = append(result, entry)
			continue
		}
		// Flat Copilot entry -- wrap in nested format
		inner := shallowCopyMap(entry)
		copilotKeysToGemini(inner)
		apmSource, _ := inner["_apm_source"].(string)
		delete(inner, "_apm_source")
		outer := map[string]interface{}{"hooks": []interface{}{inner}}
		if apmSource != "" {
			outer["_apm_source"] = apmSource
		}
		result = append(result, outer)
	}
	return result
}

func shallowCopyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copilotKeysToGemini(hook map[string]interface{}) {
	if _, hasCmd := hook["command"]; !hasCmd {
		for _, key := range []string{"bash", "powershell", "windows"} {
			if v, ok := hook[key]; ok {
				hook["command"] = v
				delete(hook, key)
				break
			}
		}
	}
	if ts, ok := hook["timeoutSec"]; ok {
		switch v := ts.(type) {
		case float64:
			hook["timeout"] = v * 1000
		case int:
			hook["timeout"] = v * 1000
		}
		delete(hook, "timeoutSec")
	}
}

// HookIntegrator handles integration of APM package hooks.
type HookIntegrator struct{}

// New returns a new HookIntegrator.
func New() *HookIntegrator { return &HookIntegrator{} }

// FindHookFiles finds all hook JSON files in a package.
// Searches .apm/hooks/ and hooks/.
func (hi *HookIntegrator) FindHookFiles(packagePath string) []string {
	var hookFiles []string
	seen := map[string]bool{}

	for _, sub := range []string{".apm/hooks", "hooks"} {
		dir := filepath.Join(packagePath, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			p := filepath.Join(dir, e.Name())
			if info, err := os.Lstat(p); err != nil || (info.Mode()&fs.ModeSymlink) != 0 {
				continue
			}
			resolved, _ := filepath.EvalSymlinks(p)
			if resolved == "" {
				resolved = p
			}
			if !seen[resolved] {
				seen[resolved] = true
				hookFiles = append(hookFiles, p)
			}
		}
	}
	return hookFiles
}

// parseHookJSON parses a hook JSON file.
func parseHookJSON(hookFile string) (map[string]interface{}, bool) {
	data, err := os.ReadFile(hookFile)
	if err != nil {
		return nil, false
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, false
	}
	return result, true
}

type scriptCopy struct {
	Source    string
	TargetRel string
}

// rewriteCommandForTarget rewrites a hook command to use installed script paths.
func (hi *HookIntegrator) rewriteCommandForTarget(
	command, packagePath, packageName, targetKey string,
	hookFileDir, rootDir string,
) (string, []scriptCopy) {
	var scripts []scriptCopy
	newCommand := command

	var scriptsBase string
	if rootDir == "" {
		switch targetKey {
		case "vscode":
			rootDir = ".github"
		case "cursor":
			rootDir = ".cursor"
		case "codex":
			rootDir = ".codex"
		case "windsurf":
			rootDir = ".windsurf"
		default:
			rootDir = ".claude"
		}
	}
	switch targetKey {
	case "vscode":
		scriptsBase = rootDir + "/hooks/scripts/" + packageName
	case "cursor":
		scriptsBase = rootDir + "/hooks/" + packageName
	case "codex":
		scriptsBase = rootDir + "/hooks/" + packageName
	case "windsurf":
		scriptsBase = rootDir + "/hooks/" + packageName
	default:
		scriptsBase = rootDir + "/hooks/" + packageName
	}

	pkgResolved, _ := filepath.EvalSymlinks(packagePath)
	if pkgResolved == "" {
		pkgResolved = packagePath
	}

	// Handle plugin root variables
	for _, match := range pluginRootRe.FindAllStringSubmatchIndex(command, -1) {
		fullVar := command[match[0]:match[1]]
		relPart := command[match[2]:match[3]]
		relPart = strings.ReplaceAll(relPart, "\\", "/")
		relPart = strings.TrimPrefix(relPart, "/")
		srcFile := filepath.Join(packagePath, relPart)
		srcResolved, _ := filepath.EvalSymlinks(srcFile)
		if srcResolved == "" {
			srcResolved = srcFile
		}
		if !strings.HasPrefix(srcResolved, pkgResolved) {
			continue
		}
		if info, err := os.Stat(srcFile); err != nil || info.IsDir() {
			continue
		}
		targetRel := scriptsBase + "/" + relPart
		scripts = append(scripts, scriptCopy{Source: srcFile, TargetRel: targetRel})
		newCommand = strings.ReplaceAll(newCommand, fullVar, targetRel)
	}

	// Handle relative ./path references
	resolveBase := hookFileDir
	if resolveBase == "" {
		resolveBase = packagePath
	}
	for _, match := range relPathRe.FindAllStringIndex(newCommand, -1) {
		relRef := newCommand[match[0]:match[1]]
		relPath := strings.TrimPrefix(relRef, "./")
		relPath = strings.TrimPrefix(relPath, ".\\")
		relPath = strings.ReplaceAll(relPath, "\\", "/")
		srcFile := filepath.Join(resolveBase, relPath)
		srcResolved, _ := filepath.EvalSymlinks(srcFile)
		if srcResolved == "" {
			srcResolved = srcFile
		}
		if !strings.HasPrefix(srcResolved, pkgResolved) {
			continue
		}
		if info, err := os.Stat(srcFile); err != nil || info.IsDir() {
			continue
		}
		targetRel := scriptsBase + "/" + relPath
		scripts = append(scripts, scriptCopy{Source: srcFile, TargetRel: targetRel})
		newCommand = strings.ReplaceAll(newCommand, relRef, targetRel)
	}
	return newCommand, scripts
}

// rewriteHooksData rewrites all command paths in a hooks JSON structure.
func (hi *HookIntegrator) rewriteHooksData(
	data map[string]interface{},
	packagePath, packageName, targetKey string,
	hookFileDir, rootDir string,
) (map[string]interface{}, []scriptCopy) {
	rewritten := deepCopyMap(data)
	var allScripts []scriptCopy

	hooksRaw, _ := rewritten["hooks"].(map[string]interface{})
	if hooksRaw == nil {
		return rewritten, nil
	}

	for eventName, rawMatchers := range hooksRaw {
		matchers, ok := rawMatchers.([]interface{})
		if !ok {
			continue
		}
		for _, rawMatcher := range matchers {
			matcher, ok := rawMatcher.(map[string]interface{})
			if !ok {
				continue
			}
			// Rewrite flat-format keys
			for _, key := range hookCommandKeys {
				if cmd, ok := matcher[key].(string); ok {
					newCmd, sc := hi.rewriteCommandForTarget(cmd, packagePath, packageName, targetKey, hookFileDir, rootDir)
					matcher[key] = newCmd
					allScripts = append(allScripts, sc...)
				}
			}
			// Rewrite nested hooks array (Claude format)
			if innerHooks, ok := matcher["hooks"].([]interface{}); ok {
				for _, rawHook := range innerHooks {
					hook, ok := rawHook.(map[string]interface{})
					if !ok {
						continue
					}
					for _, key := range hookCommandKeys {
						if cmd, ok := hook[key].(string); ok {
							newCmd, sc := hi.rewriteCommandForTarget(cmd, packagePath, packageName, targetKey, hookFileDir, rootDir)
							hook[key] = newCmd
							allScripts = append(allScripts, sc...)
						}
					}
				}
			}
		}
		_ = eventName
	}

	// Deduplicate scripts by target path
	seen := map[string]string{}
	for _, sc := range allScripts {
		if _, ok := seen[sc.TargetRel]; !ok {
			seen[sc.TargetRel] = sc.Source
		}
	}
	var uniqueScripts []scriptCopy
	for tgt, src := range seen {
		uniqueScripts = append(uniqueScripts, scriptCopy{Source: src, TargetRel: tgt})
	}
	return rewritten, uniqueScripts
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	b, _ := json.Marshal(m)
	var out map[string]interface{}
	_ = json.Unmarshal(b, &out)
	return out
}

func portableRelpath(path, base string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return strings.ReplaceAll(rel, "\\", "/")
}

// IntegratePackageHooks integrates hooks for the Copilot/VSCode target (individual JSON files).
func (hi *HookIntegrator) IntegratePackageHooks(
	packageInstallPath, projectRoot string,
	packageName string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
	rootDir string,
) *HookIntegrationResult {
	hookFiles := hi.FindHookFiles(packageInstallPath)
	hookFiles = filterHookFilesForTarget(hookFiles, "copilot")
	if len(hookFiles) == 0 {
		return &HookIntegrationResult{}
	}

	if rootDir == "" {
		rootDir = ".github"
	}
	hooksDir := filepath.Join(projectRoot, rootDir, "hooks")
	_ = os.MkdirAll(hooksDir, 0o755)

	if packageName == "" {
		packageName = filepath.Base(packageInstallPath)
	}

	var result HookIntegrationResult
	for _, hookFile := range hookFiles {
		data, ok := parseHookJSON(hookFile)
		if !ok {
			continue
		}
		rewritten, scripts := hi.rewriteHooksData(data, packageInstallPath, packageName, "vscode", filepath.Dir(hookFile), rootDir)
		stem := strings.TrimSuffix(filepath.Base(hookFile), filepath.Ext(hookFile))
		targetFilename := packageName + "-" + stem + ".json"
		targetPath := filepath.Join(hooksDir, targetFilename)
		relPath := portableRelpath(targetPath, projectRoot)

		if baseintegrator.CheckCollision(targetPath, relPath, managedFiles, force, diag) {
			continue
		}

		b, err := json.MarshalIndent(rewritten, "", "  ")
		if err != nil {
			continue
		}
		if err := os.WriteFile(targetPath, append(b, '\n'), 0o644); err != nil {
			continue
		}
		result.FilesIntegrated++
		result.TargetPaths = append(result.TargetPaths, targetPath)

		for _, sc := range scripts {
			scriptTarget := filepath.Join(projectRoot, sc.TargetRel)
			if err := os.MkdirAll(filepath.Dir(scriptTarget), 0o755); err != nil {
				continue
			}
			if baseintegrator.CheckCollision(scriptTarget, sc.TargetRel, managedFiles, force, diag) {
				continue
			}
			srcData, err := os.ReadFile(sc.Source)
			if err != nil {
				continue
			}
			if err := os.WriteFile(scriptTarget, srcData, 0o755); err != nil {
				continue
			}
			result.ScriptsCopied++
			result.TargetPaths = append(result.TargetPaths, scriptTarget)
		}
	}
	return &result
}

// integrateMergedHooks integrates hooks by merging into a target-specific JSON config.
func (hi *HookIntegrator) integrateMergedHooks(
	config mergeHookConfig,
	packageInstallPath, projectRoot string,
	packageName string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
	rootDir string,
) *HookIntegrationResult {
	empty := &HookIntegrationResult{}
	if rootDir == "" {
		rootDir = "." + config.TargetKey
	}
	targetDir := filepath.Join(projectRoot, rootDir)
	if config.RequireDir {
		if _, err := os.Stat(targetDir); err != nil {
			return empty
		}
	}

	hookFiles := hi.FindHookFiles(packageInstallPath)
	hookFiles = filterHookFilesForTarget(hookFiles, config.TargetKey)
	if len(hookFiles) == 0 {
		return empty
	}

	if packageName == "" {
		packageName = filepath.Base(packageInstallPath)
	}

	jsonPath := filepath.Join(targetDir, config.ConfigFilename)
	jsonConfig := map[string]interface{}{}
	if data, err := os.ReadFile(jsonPath); err == nil {
		_ = json.Unmarshal(data, &jsonConfig)
	}
	if _, ok := jsonConfig["hooks"]; !ok {
		jsonConfig["hooks"] = map[string]interface{}{}
	}
	hooksMap := jsonConfig["hooks"].(map[string]interface{})

	eMap := hookEventMap[config.TargetKey]
	clearedEvents := map[string]bool{}

	var result HookIntegrationResult

	for _, hookFile := range hookFiles {
		data, ok := parseHookJSON(hookFile)
		if !ok {
			continue
		}
		rewritten, scripts := hi.rewriteHooksData(data, packageInstallPath, packageName, config.TargetKey, filepath.Dir(hookFile), rootDir)

		hooksRaw, _ := rewritten["hooks"].(map[string]interface{})
		if hooksRaw == nil {
			continue
		}

		for rawEventName, rawEntries := range hooksRaw {
			entries, ok := rawEntries.([]interface{})
			if !ok {
				continue
			}
			eventName := rawEventName
			if mapped, ok := eMap[rawEventName]; ok {
				eventName = mapped
			}
			if _, ok := hooksMap[eventName]; !ok {
				hooksMap[eventName] = []interface{}{}
			}
			existingEntries := toSlice(hooksMap[eventName])

			// Transform to Gemini format
			if config.TargetKey == "gemini" {
				entries = toGeminiHookEntries(entries)
			}
			// Mark with APM source
			for _, e := range entries {
				if em, ok := e.(map[string]interface{}); ok {
					em["_apm_source"] = packageName
				}
			}

			// Idempotent upsert: clear prior entries for this package
			if !clearedEvents[eventName] {
				filtered := make([]interface{}, 0, len(existingEntries))
				for _, e := range existingEntries {
					if em, ok := e.(map[string]interface{}); ok {
						if em["_apm_source"] == packageName {
							continue
						}
					}
					filtered = append(filtered, e)
				}
				existingEntries = filtered
				clearedEvents[eventName] = true
			}
			existingEntries = append(existingEntries, entries...)

			// Deduplicate same-package entries
			existingEntries = deduplicateHookEntries(existingEntries, packageName)
			hooksMap[eventName] = existingEntries
		}
		result.FilesIntegrated++

		for _, sc := range scripts {
			scriptTarget := filepath.Join(projectRoot, sc.TargetRel)
			_ = os.MkdirAll(filepath.Dir(scriptTarget), 0o755)
			if baseintegrator.CheckCollision(scriptTarget, sc.TargetRel, managedFiles, force, diag) {
				continue
			}
			srcData, err := os.ReadFile(sc.Source)
			if err != nil {
				continue
			}
			if err := os.WriteFile(scriptTarget, srcData, 0o755); err != nil {
				continue
			}
			result.ScriptsCopied++
			result.TargetPaths = append(result.TargetPaths, scriptTarget)
		}
	}

	_ = os.MkdirAll(targetDir, 0o755)
	b, err := json.MarshalIndent(jsonConfig, "", "  ")
	if err == nil {
		_ = os.WriteFile(jsonPath, append(b, '\n'), 0o644)
	}
	return &result
}

func toSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return nil
}

func deduplicateHookEntries(entries []interface{}, packageName string) []interface{} {
	type cmpKey struct {
		source string
		cmp    string
	}
	seen := map[cmpKey]bool{}
	var result []interface{}
	for _, e := range entries {
		em, ok := e.(map[string]interface{})
		if !ok {
			result = append(result, e)
			continue
		}
		src, _ := em["_apm_source"].(string)
		if src != packageName {
			result = append(result, e)
			continue
		}
		cmpMap := map[string]interface{}{}
		for k, v := range em {
			if k != "_apm_source" {
				cmpMap[k] = v
			}
		}
		cmpBytes, _ := json.Marshal(cmpMap)
		key := cmpKey{source: src, cmp: string(cmpBytes)}
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}
	return result
}

// IntegrateHooksForTarget integrates hooks for a single target profile.
func (hi *HookIntegrator) IntegrateHooksForTarget(
	tgt *targets.TargetProfile,
	packageInstallPath, projectRoot, packageName string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
) *HookIntegrationResult {
	if tgt.Name == "copilot" {
		return hi.IntegratePackageHooks(packageInstallPath, projectRoot, packageName, force, managedFiles, diag, tgt.RootDir)
	}
	if cfg, ok := mergeHookTargets[tgt.Name]; ok {
		return hi.integrateMergedHooks(cfg, packageInstallPath, projectRoot, packageName, force, managedFiles, diag, tgt.RootDir)
	}
	return &HookIntegrationResult{}
}

// SyncStats holds cleanup statistics.
type SyncStats struct {
	FilesRemoved int
	Errors       int
}

// SyncIntegration removes APM-managed hook files.
func (hi *HookIntegrator) SyncIntegration(
	projectRoot string,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
) SyncStats {
	var stats SyncStats
	if allTargets == nil {
		for _, t := range targets.KnownTargets {
			allTargets = append(allTargets, t)
		}
	}

	hookPrefixes := hookPrefixList(allTargets)

	if managedFiles != nil {
		var deleted []string
		for relPath := range managedFiles {
			norm := strings.ReplaceAll(relPath, "\\", "/")
			if strings.Contains(norm, "..") {
				continue
			}
			if !hasAnyPrefix(norm, hookPrefixes) {
				continue
			}
			target := filepath.Join(projectRoot, relPath)
			if info, err := os.Stat(target); err != nil || info.IsDir() {
				continue
			}
			if err := os.Remove(target); err != nil {
				stats.Errors++
			} else {
				stats.FilesRemoved++
				deleted = append(deleted, target)
			}
		}
		baseintegrator.CleanupEmptyParents(deleted, projectRoot)
	} else {
		// Legacy: glob for *-apm.json
		hooksDir := filepath.Join(projectRoot, ".github", "hooks")
		entries, err := os.ReadDir(hooksDir)
		if err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), "-apm.json") {
					if err := os.Remove(filepath.Join(hooksDir, e.Name())); err != nil {
						stats.Errors++
					} else {
						stats.FilesRemoved++
					}
				}
			}
		}
	}

	// Clean APM entries from merged-hook JSON configs
	for _, tgt := range allTargets {
		cfg, ok := mergeHookTargets[tgt.Name]
		if !ok {
			continue
		}
		jsonPath := filepath.Join(projectRoot, tgt.RootDir, cfg.ConfigFilename)
		cleanApmEntriesFromJSON(jsonPath, &stats)
	}
	return stats
}

func hookPrefixList(allTargets []*targets.TargetProfile) []string {
	var out []string
	for _, tgt := range allTargets {
		if !tgt.Supports("hooks") {
			continue
		}
		sm := tgt.Primitives["hooks"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		out = append(out, effectiveRoot+"/hooks/")
	}
	return out
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func cleanApmEntriesFromJSON(jsonPath string, stats *SyncStats) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		stats.Errors++
		return
	}
	hooksRaw, ok := cfg["hooks"]
	if !ok {
		return
	}
	hooksMap, ok := hooksRaw.(map[string]interface{})
	if !ok {
		return
	}
	modified := false
	for eventName := range hooksMap {
		entries := toSlice(hooksMap[eventName])
		filtered := make([]interface{}, 0, len(entries))
		for _, e := range entries {
			if em, ok := e.(map[string]interface{}); ok {
				if _, hasSource := em["_apm_source"]; hasSource {
					modified = true
					continue
				}
			}
			filtered = append(filtered, e)
		}
		if len(filtered) == 0 {
			delete(hooksMap, eventName)
			modified = true
		} else {
			hooksMap[eventName] = filtered
		}
	}
	if len(hooksMap) == 0 {
		delete(cfg, "hooks")
		modified = true
	}
	if modified {
		b, err := json.MarshalIndent(cfg, "", "  ")
		if err == nil {
			_ = os.WriteFile(jsonPath, append(b, '\n'), 0o644)
			stats.FilesRemoved++
		}
	}
}
