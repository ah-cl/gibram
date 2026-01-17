// Package simd provides SIMD-optimized distance calculations
package simd

import (
	"math"
)

// CosineSimilarity calculates cosine similarity between two vectors
// This function automatically selects the best implementation based on CPU features
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	// Use SIMD implementation if available and vector is large enough
	if hasAVX2() && len(a) >= 8 {
		return cosineSimilarityAVX2(a, b)
	}

	// Fallback to scalar implementation
	return cosineSimilarityScalar(a, b)
}

// cosineSimilarityScalar is the baseline scalar implementation
func cosineSimilarityScalar(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// EuclideanDistance calculates L2 distance between two vectors
func EuclideanDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	// Use SIMD implementation if available and vector is large enough
	if hasAVX2() && len(a) >= 8 {
		return euclideanDistanceAVX2(a, b)
	}

	// Fallback to scalar implementation
	return euclideanDistanceScalar(a, b)
}

// euclideanDistanceScalar is the baseline scalar implementation
func euclideanDistanceScalar(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return float32(math.Sqrt(float64(sum)))
}

// DotProduct calculates dot product between two vectors
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	// Use SIMD implementation if available and vector is large enough
	if hasAVX2() && len(a) >= 8 {
		return dotProductAVX2(a, b)
	}

	// Fallback to scalar implementation
	return dotProductScalar(a, b)
}

// dotProductScalar is the baseline scalar implementation
func dotProductScalar(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}

	return sum
}

// L2Norm calculates the L2 norm (magnitude) of a vector
func L2Norm(a []float32) float32 {
	if hasAVX2() && len(a) >= 8 {
		return l2NormAVX2(a)
	}
	return l2NormScalar(a)
}

// l2NormScalar is the baseline scalar implementation
func l2NormScalar(a []float32) float32 {
	var sum float32
	for i := range a {
		sum += a[i] * a[i]
	}
	return float32(math.Sqrt(float64(sum)))
}

// hasAVX2 checks if the CPU supports AVX2 instructions
// This is implemented in simd_amd64.go for amd64 and returns false for other architectures
func hasAVX2() bool {
	return hasAVX2Check
}
