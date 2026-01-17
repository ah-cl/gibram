// Package metrics provides metrics collection for GibRAM
package metrics

import (
	"math"
	"sort"
	"sync"
)

// Histogram tracks distribution of values
type Histogram struct {
	mu     sync.Mutex
	values []float64
	count  int64
	sum    float64
	min    float64
	max    float64
}

// NewHistogram creates a new histogram
func NewHistogram() *Histogram {
	return &Histogram{
		values: make([]float64, 0, 1000),
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
	}
}

// Record records a value
func (h *Histogram) Record(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.values = append(h.values, value)
	h.count++
	h.sum += value

	if value < h.min {
		h.min = value
	}
	if value > h.max {
		h.max = value
	}

	// Keep only last 10000 values for percentile calculations
	if len(h.values) > 10000 {
		h.values = h.values[len(h.values)-10000:]
	}
}

// Stats returns histogram statistics
func (h *Histogram) Stats() *HistogramStats {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.count == 0 {
		return &HistogramStats{}
	}

	// Sort for percentiles
	sorted := make([]float64, len(h.values))
	copy(sorted, h.values)
	sort.Float64s(sorted)

	return &HistogramStats{
		Count: h.count,
		Sum:   h.sum,
		Min:   h.min,
		Max:   h.max,
		Avg:   h.sum / float64(h.count),
		P50:   percentile(sorted, 0.50),
		P90:   percentile(sorted, 0.90),
		P95:   percentile(sorted, 0.95),
		P99:   percentile(sorted, 0.99),
	}
}

// HistogramStats holds computed histogram statistics
type HistogramStats struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
	Avg   float64
	P50   float64
	P90   float64
	P95   float64
	P99   float64
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}
