// Package simd provides SIMD-optimized distance calculations
// +build !amd64

package simd

// hasAVX2Check is always false for non-amd64 architectures
var hasAVX2Check = false

// Fallback implementations for non-amd64 architectures
// These just call the scalar versions

func cosineSimilarityAVX2(a, b []float32) float32 {
	return cosineSimilarityScalar(a, b)
}

func euclideanDistanceAVX2(a, b []float32) float32 {
	return euclideanDistanceScalar(a, b)
}

func dotProductAVX2(a, b []float32) float32 {
	return dotProductScalar(a, b)
}

func l2NormAVX2(a []float32) float32 {
	return l2NormScalar(a)
}
