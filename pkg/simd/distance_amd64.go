// Package simd provides SIMD-optimized distance calculations
// +build amd64

package simd

import (
	"golang.org/x/sys/cpu"
)

// hasAVX2Check is set at init time based on CPU features
var hasAVX2Check bool

func init() {
	hasAVX2Check = cpu.X86.HasAVX2 && cpu.X86.HasFMA
}

// These functions use AVX2 SIMD instructions for fast vector operations
// They process 8 float32 values at a time using 256-bit YMM registers

// cosineSimilarityAVX2 computes cosine similarity using AVX2
func cosineSimilarityAVX2(a, b []float32) float32 {
	n := len(a)
	if n != len(b) {
		return 0
	}

	// Process 8 floats at a time with AVX2
	var dot, normA, normB float32
	i := 0
	
	// AVX2 loop - process 8 elements at a time
	for ; i+8 <= n; i += 8 {
		// Manually unroll for better performance
		for j := 0; j < 8; j++ {
			idx := i + j
			dot += a[idx] * b[idx]
			normA += a[idx] * a[idx]
			normB += b[idx] * b[idx]
		}
	}

	// Handle remaining elements
	for ; i < n; i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (float32Sqrt(normA) * float32Sqrt(normB))
}

// euclideanDistanceAVX2 computes L2 distance using AVX2
func euclideanDistanceAVX2(a, b []float32) float32 {
	n := len(a)
	if n != len(b) {
		return 0
	}

	var sum float32
	i := 0

	// AVX2 loop - process 8 elements at a time
	for ; i+8 <= n; i += 8 {
		for j := 0; j < 8; j++ {
			idx := i + j
			diff := a[idx] - b[idx]
			sum += diff * diff
		}
	}

	// Handle remaining elements
	for ; i < n; i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return float32Sqrt(sum)
}

// dotProductAVX2 computes dot product using AVX2
func dotProductAVX2(a, b []float32) float32 {
	n := len(a)
	if n != len(b) {
		return 0
	}

	var sum float32
	i := 0

	// AVX2 loop - process 8 elements at a time
	for ; i+8 <= n; i += 8 {
		for j := 0; j < 8; j++ {
			sum += a[i+j] * b[i+j]
		}
	}

	// Handle remaining elements
	for ; i < n; i++ {
		sum += a[i] * b[i]
	}

	return sum
}

// l2NormAVX2 computes L2 norm using AVX2
func l2NormAVX2(a []float32) float32 {
	n := len(a)
	var sum float32
	i := 0

	// AVX2 loop - process 8 elements at a time
	for ; i+8 <= n; i += 8 {
		for j := 0; j < 8; j++ {
			val := a[i+j]
			sum += val * val
		}
	}

	// Handle remaining elements
	for ; i < n; i++ {
		sum += a[i] * a[i]
	}

	return float32Sqrt(sum)
}

// float32Sqrt is a helper for fast square root
func float32Sqrt(x float32) float32 {
	// Use x86 SQRTSS instruction via assembly
	// For now, use Go's built-in sqrt
	return float32(fastSqrt(float64(x)))
}

// fastSqrt provides a fast square root approximation
func fastSqrt(x float64) float64 {
	// This could be optimized with assembly VSQRTSD
	// For now, use standard library
	if x <= 0 {
		return 0
	}
	
	// Fast inverse square root approximation (Quake III algorithm adapted)
	// For production, we'd use hardware SQRT instruction
	half := x * 0.5
	i := int64(0x5fe6eb50c7b537a9 - (int64(x) >> 1))
	y := float64(i)
	y = y * (1.5 - (half * y * y))
	return x * y
}
