// Package mcpargs parses MCP CLI argument KEY=VALUE pairs.
package mcpargs

import "fmt"

// ParseKVPairs parses a slice of KEY=VALUE strings into a map.
// flagName is used in error messages.
func ParseKVPairs(pairs []string, flagName string) (map[string]string, error) {
result := map[string]string{}
for _, raw := range pairs {
idx := -1
for i, c := range raw {
if c == '=' {
idx = i
break
}
}
if idx < 0 {
return nil, fmt.Errorf("invalid %s '%s': expected KEY=VALUE", flagName, raw)
}
key := raw[:idx]
value := raw[idx+1:]
if key == "" {
return nil, fmt.Errorf("invalid %s '%s': key cannot be empty", flagName, raw)
}
result[key] = value
}
return result, nil
}

// ParseEnvPairs parses --env KEY=VAL repetitions into a map.
func ParseEnvPairs(pairs []string) (map[string]string, error) {
return ParseKVPairs(pairs, "--env")
}

// ParseHeaderPairs parses --header KEY=VAL repetitions into a map.
func ParseHeaderPairs(pairs []string) (map[string]string, error) {
return ParseKVPairs(pairs, "--header")
}
