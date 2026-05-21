package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var gateBenchSink ScanVerdict

func benchmarkMixedContent(size int) string {
	block := strings.Repeat("Hello world. ", 5) + "\u200b" + strings.Repeat("More text. ", 5) + "\u200c"
	var sb strings.Builder
	for sb.Len() < size {
		sb.WriteString(block)
	}
	return sb.String()[:size]
}

func BenchmarkGateCheckMixedContent100KB(b *testing.B) {
	root := b.TempDir()
	path := filepath.Join(root, "bench.md")
	if err := os.WriteFile(path, []byte(benchmarkMixedContent(100_000)), 0o644); err != nil {
		b.Fatal(err)
	}
	gate := New(BlockPolicy, false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gateBenchSink = gate.CheckFile(path)
	}
}

func BenchmarkGateCheckWideFileSet50(b *testing.B) {
	root := b.TempDir()
	paths := make([]string, 0, 50)
	for i := 0; i < 50; i++ {
		path := filepath.Join(root, "pkg", "file.md")
		if i > 0 {
			path = filepath.Join(root, "pkg", fmt.Sprintf("file-%02d.md", i))
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("safe benchmark content\n"), 0o644); err != nil {
			b.Fatal(err)
		}
		paths = append(paths, path)
	}
	gate := New(BlockPolicy, false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gateBenchSink = gate.Check(paths)
	}
}
