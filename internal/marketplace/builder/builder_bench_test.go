package builder

import (
	"fmt"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/ymlschema"
)

var builderBenchSink map[string]interface{}

func benchmarkMarketplaceConfig(count int) (*ymlschema.MarketplaceConfig, []ResolvedPackage) {
	cfg := &ymlschema.MarketplaceConfig{
		Name:        "bench-marketplace",
		Description: "benchmark marketplace",
		Version:     "1.0.0",
		Owner:       ymlschema.MarketplaceOwner{Name: "Benchmark"},
		Output:      "marketplace.json",
		Metadata:    map[string]interface{}{"category": "benchmark"},
	}
	resolved := make([]ResolvedPackage, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("plugin-%03d", i)
		cfg.Packages = append(cfg.Packages, ymlschema.PackageEntry{
			Name:        name,
			Source:      fmt.Sprintf("bench-org/%s", name),
			Version:     "1.0.0",
			Description: "Synthetic benchmark plugin",
			Tags:        []string{"bench", "migration"},
		})
		resolved = append(resolved, ResolvedPackage{
			Name:             name,
			SourceRepo:       fmt.Sprintf("bench-org/%s", name),
			Ref:              "v1.0.0",
			SHA:              fmt.Sprintf("%040d", i),
			RequestedVersion: "1.0.0",
			Tags:             []string{"bench", "migration"},
		})
	}
	return cfg, resolved
}

func BenchmarkComposeMarketplaceJSON50Packages(b *testing.B) {
	cfg, resolved := benchmarkMarketplaceConfig(50)
	builder := FromConfig(cfg, b.TempDir(), DefaultBuildOptions())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := builder.ComposeMarketplaceJSON(resolved)
		if err != nil {
			b.Fatal(err)
		}
		builderBenchSink = doc
	}
}
