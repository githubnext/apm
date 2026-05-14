// Package localbundle provides helpers for installing local APM bundles.
//
// Migrated from src/apm_cli/install/local_bundle_handler.py
package localbundle

import (
"encoding/json"
"os"
"path/filepath"
"strings"
)

// MCPServerSpec represents a single MCP server entry from .mcp.json.
type MCPServerSpec struct {
Name      string
Transport string
Command   string
Args      []string
Env       map[string]string
URL       string
Registry  bool
Raw       map[string]interface{}
}

// ParseBundleMCPServers parses <bundle>/.mcp.json into MCPServerSpec entries.
// Returns an empty slice when the file is missing or malformed.
func ParseBundleMCPServers(bundleDir string) []MCPServerSpec {
var mcpPath string
entries, err := os.ReadDir(bundleDir)
if err != nil {
return nil
}
for _, e := range entries {
if !e.IsDir() && strings.ToLower(e.Name()) == ".mcp.json" {
mcpPath = filepath.Join(bundleDir, e.Name())
break
}
}
if mcpPath == "" {
return nil
}

data, err := os.ReadFile(mcpPath)
if err != nil {
return nil
}
var root map[string]interface{}
if err := json.Unmarshal(data, &root); err != nil {
return nil
}

serversRaw, ok := root["mcpServers"]
if !ok {
return nil
}
serversMap, ok := serversRaw.(map[string]interface{})
if !ok {
return nil
}

var out []MCPServerSpec
for name, cfgRaw := range serversMap {
cfg, ok := cfgRaw.(map[string]interface{})
if !ok {
continue
}
spec := MCPServerSpec{
Name:    name,
Raw:     cfg,
Command: strVal(cfg["command"]),
URL:     strVal(cfg["url"]),
}
// transport / type
if t := strVal(cfg["type"]); t != "" {
spec.Transport = t
} else {
spec.Transport = strVal(cfg["transport"])
}
// args
if argsRaw, ok := cfg["args"]; ok {
if argsSlice, ok := argsRaw.([]interface{}); ok {
for _, a := range argsSlice {
spec.Args = append(spec.Args, strVal(a))
}
}
}
// env
spec.Env = strMapVal(cfg["env"])
out = append(out, spec)
}
return out
}

// BundleMCPPresent returns true if the bundle directory contains a .mcp.json file.
func BundleMCPPresent(bundleDir string) bool {
entries, err := os.ReadDir(bundleDir)
if err != nil {
return false
}
for _, e := range entries {
if !e.IsDir() && strings.ToLower(e.Name()) == ".mcp.json" {
return true
}
}
return false
}

func strVal(v interface{}) string {
if v == nil {
return ""
}
if s, ok := v.(string); ok {
return s
}
return ""
}

func strMapVal(v interface{}) map[string]string {
if v == nil {
return nil
}
switch m := v.(type) {
case map[string]interface{}:
result := make(map[string]string, len(m))
for k, val := range m {
result[k] = strVal(val)
}
return result
case map[string]string:
return m
}
return nil
}
