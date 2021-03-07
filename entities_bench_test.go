package wcrawler_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler"
)

func BenchmarkMarshalEdgesSet(b *testing.B) {
	benchmarks := map[string]struct {
		n int
	}{
		"1000 entries":   {n: 1000},
		"10000 entries":  {n: 10000},
		"100000 entries": {n: 100000},
	}

	for name, bm := range benchmarks {
		b.Run(name, func(b *testing.B) {

			// Initialize and load EdgesSet
			set := wcrawler.NewEdgesSet()
			for i := 0; i < bm.n; i++ {
				set.Add(i)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = set.MarshalJSON()
			}
		})
	}
}

func BenchmarkDumpEdgesSet(b *testing.B) {
	benchmarks := map[string]struct {
		n int
	}{
		"1000 entries":   {n: 1000},
		"10000 entries":  {n: 10000},
		"100000 entries": {n: 100000},
	}

	for name, bm := range benchmarks {
		b.Run(name, func(b *testing.B) {

			// Initialize and load EdgesSet
			set := wcrawler.NewEdgesSet()
			for i := 0; i < bm.n; i++ {
				set.Add(i)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = set.Dump()
			}
		})
	}
}
