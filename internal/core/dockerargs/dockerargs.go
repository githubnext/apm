// Package dockerargs handles Docker argument processing with deduplication.
package dockerargs

// ProcessDockerArgs processes Docker arguments with environment variable deduplication.
func ProcessDockerArgs(baseArgs []string, envVars map[string]string) []string {
result := []string{}
envVarsAdded := map[string]bool{}
hasInteractive := false
hasRM := false

for _, arg := range baseArgs {
if arg == "-i" || arg == "--interactive" {
hasInteractive = true
}
if arg == "--rm" {
hasRM = true
}
}

for _, arg := range baseArgs {
result = append(result, arg)
if arg == "run" {
if !hasInteractive {
result = append(result, "-i")
}
if !hasRM {
result = append(result, "--rm")
}
for name, val := range envVars {
if !envVarsAdded[name] {
result = append(result, "-e", name+"="+val)
envVarsAdded[name] = true
}
}
}
}
return result
}

// ExtractEnvVars extracts -e flags from Docker args.
func ExtractEnvVars(args []string) (cleanArgs []string, envVars map[string]string) {
envVars = map[string]string{}
i := 0
for i < len(args) {
if args[i] == "-e" && i+1 < len(args) {
spec := args[i+1]
idx := -1
for j, c := range spec {
if c == '=' {
idx = j
break
}
}
if idx >= 0 {
envVars[spec[:idx]] = spec[idx+1:]
} else {
envVars[spec] = "${" + spec + "}"
}
i += 2
} else {
cleanArgs = append(cleanArgs, args[i])
i++
}
}
return cleanArgs, envVars
}

// MergeEnvVars merges environment variables, with newEnv taking precedence.
func MergeEnvVars(existing, newEnv map[string]string) map[string]string {
merged := map[string]string{}
for k, v := range existing {
merged[k] = v
}
for k, v := range newEnv {
merged[k] = v
}
return merged
}
