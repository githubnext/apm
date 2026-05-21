package exclude

import (
	"strings"
	"testing"
)

var excludeBenchSink bool

func benchmarkPathParts(depth int) []string {
	parts := make([]string, 0, depth)
	for i := 0; i < depth-1; i++ {
		parts = append(parts, string(rune('a'+i%26)))
	}
	return append(parts, "test.py")
}

func benchmarkDoubleStarPattern(starSegments int) []string {
	labels := []string{"a", "b", "c", "d", "e"}
	parts := make([]string, 0, starSegments*2+1)
	for i := 0; i < starSegments; i++ {
		parts = append(parts, "**", labels[i%len(labels)])
	}
	return append(parts, "*.py")
}

func BenchmarkMatchDoubleStar3SegmentsDepth20(b *testing.B) {
	pathParts := benchmarkPathParts(20)
	patternParts := benchmarkDoubleStarPattern(3)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		excludeBenchSink = matchDoubleStar(pathParts, patternParts)
	}
}

func BenchmarkShouldExcludeSimpleGlob1000Shape(b *testing.B) {
	base := b.TempDir()
	path := base + "/src/module/deep/nested/file.py"
	patterns, err := ValidateExcludePatterns([]string{"**/*.py"})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		excludeBenchSink = ShouldExclude(path, base, patterns)
	}
}

func BenchmarkMatchPatternNoDoubleStar(b *testing.B) {
	path := strings.Repeat("module/", 4) + "test_example.md"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		excludeBenchSink = matchesPattern(path, "**/*.md")
	}
}
