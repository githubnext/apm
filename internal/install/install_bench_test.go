package install_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
)

var installBenchSink int

var manifestBenchmarkPrefixes = []string{
	".github/prompts/",
	".github/agents/",
	".claude/agents/",
	".claude/commands/",
	".github/skills/",
	".github/hooks/",
}

func buildBenchmarkManagedFiles(packages, filesPerPackage int) map[string]struct{} {
	managed := make(map[string]struct{}, packages*filesPerPackage)
	for i := 0; i < packages; i++ {
		prefix := manifestBenchmarkPrefixes[i%len(manifestBenchmarkPrefixes)]
		for j := 0; j < filesPerPackage; j++ {
			managed[fmt.Sprintf("%spkg-%d-file-%d.md", prefix, i, j)] = struct{}{}
		}
	}
	return managed
}

func BenchmarkManifestNormalizeManagedFiles50x5(b *testing.B) {
	managed := buildBenchmarkManagedFiles(50, 5)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		installBenchSink = len(baseintegrator.NormalizeManagedFiles(managed))
	}
}

func BenchmarkManifestCheckCollision50x5(b *testing.B) {
	projectRoot := b.TempDir()
	relPath := ".github/prompts/pkg-0-file-0.md"
	targetPath := filepath.Join(projectRoot, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("managed"), 0o644); err != nil {
		b.Fatal(err)
	}
	managed := baseintegrator.NormalizeManagedFiles(buildBenchmarkManagedFiles(50, 5))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if baseintegrator.CheckCollision(targetPath, relPath, managed, false, nil) {
			installBenchSink++
		}
	}
}

func BenchmarkManifestPartitionManagedFiles50x5(b *testing.B) {
	managed := baseintegrator.NormalizeManagedFiles(buildBenchmarkManagedFiles(50, 5))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buckets := baseintegrator.PartitionManagedFiles(managed, nil)
		installBenchSink = len(buckets)
	}
}

func BenchmarkManifestSyncRemoveFiles50x5(b *testing.B) {
	projectRoot := b.TempDir()
	managed := baseintegrator.NormalizeManagedFiles(buildBenchmarkManagedFiles(50, 5))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stats := baseintegrator.SyncRemoveFiles(projectRoot, managed, ".github/prompts/", "", "", nil, nil)
		installBenchSink = stats.FilesRemoved + stats.Errors
	}
}

func BenchmarkManifestCleanupEmptyParents50x5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root := b.TempDir()
		deleted := make([]string, 0, 100)
		for n := 0; n < 100; n++ {
			dir := filepath.Join(root, fmt.Sprintf("d%d", n%6), "sub0", "sub1", "sub2")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				b.Fatal(err)
			}
			p := filepath.Join(dir, fmt.Sprintf("file-%d.md", n))
			if err := os.WriteFile(p, nil, 0o644); err != nil {
				b.Fatal(err)
			}
			if err := os.Remove(p); err != nil {
				b.Fatal(err)
			}
			deleted = append(deleted, p)
		}

		b.StartTimer()
		baseintegrator.CleanupEmptyParents(deleted, root)
		b.StopTimer()
	}
}

func BenchmarkManifestScopedUninstallSet50x5(b *testing.B) {
	packageFiles := make(map[int][]string, 50)
	for i := 0; i < 50; i++ {
		prefix := manifestBenchmarkPrefixes[i%len(manifestBenchmarkPrefixes)]
		for j := 0; j < 5; j++ {
			packageFiles[i] = append(packageFiles[i], fmt.Sprintf("%spkg-%d-file-%d.md", prefix, i, j))
		}
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		removedFiles := make([]string, 0, 25)
		for pkg := 0; pkg < 5; pkg++ {
			removedFiles = append(removedFiles, packageFiles[pkg]...)
		}
		count := 0
		for _, prefix := range manifestBenchmarkPrefixes {
			for _, relPath := range removedFiles {
				if strings.HasPrefix(relPath, prefix) {
					count++
				}
			}
		}
		installBenchSink = count
	}
}
