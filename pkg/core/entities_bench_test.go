package core_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
)

func BenchmarkMarshalEdgesSet(b *testing.B) {
	benchmarks := map[string]struct {
		n int
	}{
		"1000 entries": {n: 1000},
	}

	for name, bm := range benchmarks {
		b.Run(name, func(b *testing.B) {
			// b.ReportAllocs()

			// Initialize and load EdgesSet
			set := core.NewEdgesSet()
			for i := 0; i < bm.n; i++ {
				set.Add(i)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = set.MarshalJSON()
			}
		})
	}
}
