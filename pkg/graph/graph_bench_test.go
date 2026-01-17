// Package graph provides graph algorithm benchmarks
package graph

import (
	"testing"

	"github.com/gibram-io/gibram/pkg/types"
)

// =============================================================================
// Benchmark Helpers
// =============================================================================

// createBenchGraph creates a graph with n entities and approximately n*avgDegree/2 relationships
func createBenchGraph(n int, avgDegree int) (*mockEntityStore, *mockRelationshipStore, []uint64) {
	entityStore := newMockEntityStore()
	relStore := newMockRelationshipStore()

	entityIDs := make([]uint64, n)
	for i := 0; i < n; i++ {
		entityIDs[i] = uint64(i + 1)
		entityStore.Add(&types.Entity{
			ID:    entityIDs[i],
			Title: "Entity" + benchItoa(i),
			Type:  "benchmark",
		})
	}

	// Create relationships with pseudo-random pattern
	relID := uint64(1)
	edgesPerNode := avgDegree / 2
	if edgesPerNode < 1 {
		edgesPerNode = 1
	}

	for i := 0; i < n; i++ {
		for j := 0; j < edgesPerNode; j++ {
			target := (i + j + 1) % n
			if target != i {
				relStore.Add(&types.Relationship{
					ID:       relID,
					SourceID: entityIDs[i],
					TargetID: entityIDs[target],
					Type:     "BENCH_REL",
					Weight:   1.0,
				})
				relID++
			}
		}
	}

	return entityStore, relStore, entityIDs
}

func benchItoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte(i%10) + '0'
		i /= 10
	}
	return string(buf[pos:])
}

// =============================================================================
// Leiden Benchmarks
// =============================================================================

func BenchmarkLeiden_100(b *testing.B) {
	entityStore, relStore, _ := createBenchGraph(100, 4)
	config := DefaultLeidenConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leiden := NewLeiden(entityStore, relStore, config)
		leiden.ComputeHierarchicalCommunities()
	}
}

func BenchmarkLeiden_500(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	entityStore, relStore, _ := createBenchGraph(500, 4)
	config := DefaultLeidenConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leiden := NewLeiden(entityStore, relStore, config)
		leiden.ComputeHierarchicalCommunities()
	}
}

func BenchmarkLeiden_1K(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	entityStore, relStore, _ := createBenchGraph(1000, 4)
	config := DefaultLeidenConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leiden := NewLeiden(entityStore, relStore, config)
		leiden.ComputeHierarchicalCommunities()
	}
}

// =============================================================================
// BFS Traversal Benchmarks
// =============================================================================

func BenchmarkBFSTraversal_100_Hop1(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)
	seeds := []uint64{entityIDs[0]}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFSTraversal(seeds, relStore, 1, 100)
	}
}

func BenchmarkBFSTraversal_100_Hop2(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)
	seeds := []uint64{entityIDs[0]}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFSTraversal(seeds, relStore, 2, 100)
	}
}

func BenchmarkBFSTraversal_100_Hop3(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)
	seeds := []uint64{entityIDs[0]}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFSTraversal(seeds, relStore, 3, 100)
	}
}

func BenchmarkBFSTraversal_1K_Hop2(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(1000, 4)
	seeds := []uint64{entityIDs[0]}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFSTraversal(seeds, relStore, 2, 200)
	}
}

func BenchmarkBFSTraversal_MultiSeed(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(500, 4)
	seeds := entityIDs[:5] // 5 seeds

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFSTraversal(seeds, relStore, 2, 100)
	}
}

// =============================================================================
// PageRank Benchmarks
// =============================================================================

func BenchmarkPageRank_100_10Iter(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PageRank(entityIDs, relStore, 0.85, 10)
	}
}

func BenchmarkPageRank_100_20Iter(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PageRank(entityIDs, relStore, 0.85, 20)
	}
}

func BenchmarkPageRank_500_10Iter(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(500, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PageRank(entityIDs, relStore, 0.85, 10)
	}
}

func BenchmarkPageRank_1K_10Iter(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(1000, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PageRank(entityIDs, relStore, 0.85, 10)
	}
}

// =============================================================================
// Connected Components Benchmarks
// =============================================================================

func BenchmarkConnectedComponents_100(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(100, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConnectedComponents(entityIDs, relStore)
	}
}

func BenchmarkConnectedComponents_500(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(500, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConnectedComponents(entityIDs, relStore)
	}
}

func BenchmarkConnectedComponents_1K(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(1000, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConnectedComponents(entityIDs, relStore)
	}
}

// =============================================================================
// Betweenness Centrality Benchmarks
// =============================================================================

func BenchmarkBetweenness_50(b *testing.B) {
	_, relStore, entityIDs := createBenchGraph(50, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Betweenness(entityIDs, relStore, 0) // 0 = use all nodes
	}
}

func BenchmarkBetweenness_100(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping in short mode")
	}

	_, relStore, entityIDs := createBenchGraph(100, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Betweenness(entityIDs, relStore, 0)
	}
}

// Note: Betweenness centrality is O(V*E) so gets expensive quickly

// =============================================================================
// Graph Density Benchmarks
// =============================================================================

func BenchmarkLeiden_Sparse(b *testing.B) {
	// avgDegree = 2 (sparse)
	entityStore, relStore, _ := createBenchGraph(200, 2)
	config := DefaultLeidenConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leiden := NewLeiden(entityStore, relStore, config)
		leiden.ComputeHierarchicalCommunities()
	}
}

func BenchmarkLeiden_Dense(b *testing.B) {
	// avgDegree = 8 (denser)
	entityStore, relStore, _ := createBenchGraph(200, 8)
	config := DefaultLeidenConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leiden := NewLeiden(entityStore, relStore, config)
		leiden.ComputeHierarchicalCommunities()
	}
}
