package apmresolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var resolverBenchSink int

func writeBenchmarkApmYML(b *testing.B, dir string, deps []string) {
	b.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		b.Fatal(err)
	}
	var sb strings.Builder
	sb.WriteString("name: bench-pkg\nversion: 1.0.0\n")
	if len(deps) > 0 {
		sb.WriteString("dependencies:\n")
		for _, dep := range deps {
			sb.WriteString("  - ")
			sb.WriteString(dep)
			sb.WriteByte('\n')
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(sb.String()), 0o644); err != nil {
		b.Fatal(err)
	}
}

func setupBenchmarkWideFan(b *testing.B, count int) string {
	b.Helper()
	root := b.TempDir()
	modules := filepath.Join(root, "apm_modules")
	deps := make([]string, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("pkg-%d", i)
		deps = append(deps, "org/"+name)
		writeBenchmarkApmYML(b, filepath.Join(modules, name), nil)
	}
	writeBenchmarkApmYML(b, root, deps)
	return root
}

func setupBenchmarkLinearChain(b *testing.B, length int) string {
	b.Helper()
	root := b.TempDir()
	modules := filepath.Join(root, "apm_modules")
	for i := 0; i < length; i++ {
		deps := []string(nil)
		if i < length-1 {
			deps = []string{fmt.Sprintf("org/pkg-%d", i+1)}
		}
		writeBenchmarkApmYML(b, filepath.Join(modules, fmt.Sprintf("pkg-%d", i)), deps)
	}
	writeBenchmarkApmYML(b, root, []string{"org/pkg-0"})
	return root
}

func BenchmarkBuildDependencyTreeWideFan50(b *testing.B) {
	root := setupBenchmarkWideFan(b, 50)
	apmModules := filepath.Join(root, "apm_modules")
	rootYML := filepath.Join(root, "apm.yml")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolver := New(Options{MaxDepth: 50, ApmModulesDir: apmModules, MaxParallel: 1})
		tree := resolver.buildDependencyTree(rootYML)
		resolverBenchSink = tree.MaxDepth + len(tree.GetNodesAtDepth(1))
	}
}

func BenchmarkBuildDependencyTreeLinearChain50(b *testing.B) {
	root := setupBenchmarkLinearChain(b, 50)
	apmModules := filepath.Join(root, "apm_modules")
	rootYML := filepath.Join(root, "apm.yml")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolver := New(Options{MaxDepth: 50, ApmModulesDir: apmModules, MaxParallel: 1})
		tree := resolver.buildDependencyTree(rootYML)
		resolverBenchSink = tree.MaxDepth
	}
}
