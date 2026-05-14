// Package integrity verifies cached git checkout integrity.
package integrity

import (
"os"
"path/filepath"
"strings"
)

// ReadHeadSHA returns the resolved 40-char SHA at HEAD, or empty string on failure.
func ReadHeadSHA(checkoutDir string) string {
gitPath := filepath.Join(checkoutDir, ".git")
info, err := os.Stat(gitPath)
if err != nil {
return ""
}

var gitDir string
if !info.IsDir() {
content, err := os.ReadFile(gitPath)
if err != nil {
return ""
}
line := strings.TrimSpace(string(content))
if !strings.HasPrefix(line, "gitdir:") {
return ""
}
target := strings.TrimSpace(line[len("gitdir:"):])
abs, err := filepath.Abs(filepath.Join(checkoutDir, target))
if err != nil {
return ""
}
gitDir = abs
} else {
gitDir = gitPath
}

headPath := filepath.Join(gitDir, "HEAD")
headContent, err := os.ReadFile(headPath)
if err != nil {
return ""
}
head := strings.TrimSpace(string(headContent))
if strings.HasPrefix(head, "ref: ") {
refName := strings.TrimPrefix(head, "ref: ")
refFile := filepath.Join(gitDir, refName)
data, err := os.ReadFile(refFile)
if err != nil {
// Try packed-refs
return resolvePackedRef(gitDir, refName)
}
return strings.TrimSpace(string(data))
}
return head
}

func resolvePackedRef(gitDir, refName string) string {
data, err := os.ReadFile(filepath.Join(gitDir, "packed-refs"))
if err != nil {
return ""
}
for _, line := range strings.Split(string(data), "\n") {
if strings.HasSuffix(line, " "+refName) {
parts := strings.Fields(line)
if len(parts) >= 1 {
return parts[0]
}
}
}
return ""
}

// VerifyCheckout checks that the checkout's HEAD matches expectedSHA.
func VerifyCheckout(checkoutDir, expectedSHA string) bool {
actual := ReadHeadSHA(checkoutDir)
return actual != "" && actual == expectedSHA
}
