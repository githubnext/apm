// Package buildid stabilizes build IDs in compiled outputs.
package buildid

import (
"crypto/sha256"
"fmt"
"strings"

"github.com/githubnext/apm/internal/compilation/compilationconst"
)

// StabilizeBuildID replaces BuildIDPlaceholder with a deterministic 12-char SHA256 hash.
// It is idempotent: returns content unchanged if no placeholder is present.
func StabilizeBuildID(content string) string {
lines := strings.Split(content, "\n")
trailingNL := strings.HasSuffix(content, "\n")

// Remove trailing empty string from Split when content ends with newline.
if trailingNL && len(lines) > 0 && lines[len(lines)-1] == "" {
lines = lines[:len(lines)-1]
}

idx := -1
for i, line := range lines {
if line == compilationconst.BuildIDPlaceholder {
idx = i
break
}
}
if idx < 0 {
return content
}

hashLines := make([]string, 0, len(lines)-1)
for i, line := range lines {
if i != idx {
hashLines = append(hashLines, line)
}
}

sum := sha256.Sum256([]byte(strings.Join(hashLines, "\n")))
buildID := fmt.Sprintf("%x", sum)[:12]
lines[idx] = fmt.Sprintf("<!-- Build ID: %s -->", buildID)

result := strings.Join(lines, "\n")
if trailingNL {
result += "\n"
}
return result
}
