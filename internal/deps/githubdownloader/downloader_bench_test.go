package githubdownloader

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"testing"
)

var downloaderBenchSink int

func benchmarkSHA(index int) string {
	sum := sha1.Sum([]byte(fmt.Sprintf("ref-%d", index)))
	return fmt.Sprintf("%x", sum)
}

func generateLsRemoteOutput(refCount int) string {
	var sb strings.Builder
	tagCount := int(float64(refCount) * 0.6)
	branchCount := refCount - tagCount
	for i := 0; i < tagCount; i++ {
		name := fmt.Sprintf("v%d.%d.%d", i/100, (i/10)%10, i%10)
		sha := benchmarkSHA(i)
		if i%3 == 0 {
			sb.WriteString(benchmarkSHA(i + 10000))
			sb.WriteString("\trefs/tags/")
			sb.WriteString(name)
			sb.WriteByte('\n')
			sb.WriteString(sha)
			sb.WriteString("\trefs/tags/")
			sb.WriteString(name)
			sb.WriteString("^{}\n")
		} else {
			sb.WriteString(sha)
			sb.WriteString("\trefs/tags/")
			sb.WriteString(name)
			sb.WriteByte('\n')
		}
	}
	for i := 0; i < branchCount; i++ {
		branch := fmt.Sprintf("feature-%d", i)
		if i == 0 {
			branch = "main"
		} else if i == 1 {
			branch = "develop"
		}
		sb.WriteString(benchmarkSHA(i + 5000))
		sb.WriteString("\trefs/heads/")
		sb.WriteString(branch)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func BenchmarkParseLsRemoteOutput500Refs(b *testing.B) {
	output := generateLsRemoteOutput(500)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		refs := ParseLsRemoteOutput(output)
		downloaderBenchSink = len(refs)
	}
}

func BenchmarkSortRemoteRefs500Refs(b *testing.B) {
	refs := ParseLsRemoteOutput(generateLsRemoteOutput(500))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sorted := SortRemoteRefs(refs)
		downloaderBenchSink = len(sorted)
	}
}
