// Package metrics - comprehensive tests for metrics collection
package metrics

import (
	"math"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Histogram Tests
// =============================================================================

func TestNewHistogram(t *testing.T) {
	h := NewHistogram()
	if h == nil {
		t.Fatal("NewHistogram() returned nil")
	}
}

func TestHistogram_Record(t *testing.T) {
	h := NewHistogram()

	h.Record(10.0)
	h.Record(20.0)
	h.Record(30.0)

	stats := h.Stats()
	if stats.Count != 3 {
		t.Errorf("Count = %d, want 3", stats.Count)
	}
}

func TestHistogram_Stats_Empty(t *testing.T) {
	h := NewHistogram()

	stats := h.Stats()
	if stats.Count != 0 {
		t.Errorf("Empty histogram Count = %d, want 0", stats.Count)
	}
}

func TestHistogram_Stats_Sum(t *testing.T) {
	h := NewHistogram()

	h.Record(10.0)
	h.Record(20.0)
	h.Record(30.0)

	stats := h.Stats()
	if stats.Sum != 60.0 {
		t.Errorf("Sum = %f, want 60.0", stats.Sum)
	}
}

func TestHistogram_Stats_MinMax(t *testing.T) {
	h := NewHistogram()

	h.Record(5.0)
	h.Record(100.0)
	h.Record(50.0)

	stats := h.Stats()
	if stats.Min != 5.0 {
		t.Errorf("Min = %f, want 5.0", stats.Min)
	}
	if stats.Max != 100.0 {
		t.Errorf("Max = %f, want 100.0", stats.Max)
	}
}

func TestHistogram_Stats_Avg(t *testing.T) {
	h := NewHistogram()

	h.Record(10.0)
	h.Record(20.0)
	h.Record(30.0)

	stats := h.Stats()
	if stats.Avg != 20.0 {
		t.Errorf("Avg = %f, want 20.0", stats.Avg)
	}
}

func TestHistogram_Stats_Percentiles(t *testing.T) {
	h := NewHistogram()

	// Add values 1-100
	for i := 1; i <= 100; i++ {
		h.Record(float64(i))
	}

	stats := h.Stats()

	// P50 should be around 50
	if stats.P50 < 45 || stats.P50 > 55 {
		t.Errorf("P50 = %f, expected around 50", stats.P50)
	}

	// P90 should be around 90
	if stats.P90 < 85 || stats.P90 > 95 {
		t.Errorf("P90 = %f, expected around 90", stats.P90)
	}

	// P99 should be around 99
	if stats.P99 < 95 || stats.P99 > 100 {
		t.Errorf("P99 = %f, expected around 99", stats.P99)
	}
}

func TestHistogram_LargeDataset(t *testing.T) {
	h := NewHistogram()

	// Record more than 10000 values (tests truncation)
	for i := 0; i < 15000; i++ {
		h.Record(float64(i))
	}

	stats := h.Stats()

	// Should have recorded 15000 values but only keep last 10000
	if stats.Count != 15000 {
		t.Errorf("Count = %d, want 15000", stats.Count)
	}

	// Percentiles should still work on truncated data
	if stats.P50 <= 0 {
		t.Error("P50 should be positive")
	}
}

func TestHistogram_Concurrent(t *testing.T) {
	h := NewHistogram()

	var wg sync.WaitGroup
	const n = 1000

	// Concurrent recording
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			h.Record(val)
		}(float64(i))
	}

	// Concurrent stats reading
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Stats()
		}()
	}

	wg.Wait()

	stats := h.Stats()
	if stats.Count != n {
		t.Errorf("Count = %d, want %d", stats.Count, n)
	}
}

func TestHistogram_NegativeValues(t *testing.T) {
	h := NewHistogram()

	h.Record(-10.0)
	h.Record(0.0)
	h.Record(10.0)

	stats := h.Stats()
	if stats.Min != -10.0 {
		t.Errorf("Min = %f, want -10.0", stats.Min)
	}
	if stats.Max != 10.0 {
		t.Errorf("Max = %f, want 10.0", stats.Max)
	}
	if stats.Sum != 0.0 {
		t.Errorf("Sum = %f, want 0.0", stats.Sum)
	}
}

func TestHistogram_SingleValue(t *testing.T) {
	h := NewHistogram()

	h.Record(42.0)

	stats := h.Stats()
	if stats.Count != 1 {
		t.Errorf("Count = %d, want 1", stats.Count)
	}
	if stats.Min != 42.0 {
		t.Errorf("Min = %f, want 42.0", stats.Min)
	}
	if stats.Max != 42.0 {
		t.Errorf("Max = %f, want 42.0", stats.Max)
	}
	if stats.Avg != 42.0 {
		t.Errorf("Avg = %f, want 42.0", stats.Avg)
	}
	// All percentiles should be the same value
	if stats.P50 != 42.0 || stats.P90 != 42.0 || stats.P99 != 42.0 {
		t.Error("Percentiles should all be 42.0 for single value")
	}
}

// =============================================================================
// HistogramStats Tests
// =============================================================================

func TestHistogramStats_Structure(t *testing.T) {
	stats := &HistogramStats{
		Count: 1000,
		Sum:   50000.0,
		Min:   1.0,
		Max:   100.0,
		Avg:   50.0,
		P50:   50.0,
		P90:   90.0,
		P95:   95.0,
		P99:   99.0,
	}

	if stats.Count != 1000 {
		t.Error("Count incorrect")
	}
	if stats.Avg != 50.0 {
		t.Error("Avg incorrect")
	}
	if stats.P95 != 95.0 {
		t.Error("P95 incorrect")
	}
}

// =============================================================================
// Percentile Function Tests
// =============================================================================

func TestPercentile_Empty(t *testing.T) {
	result := percentile([]float64{}, 0.5)
	if result != 0 {
		t.Errorf("percentile(empty) = %f, want 0", result)
	}
}

func TestPercentile_Single(t *testing.T) {
	result := percentile([]float64{42.0}, 0.5)
	if result != 42.0 {
		t.Errorf("percentile([42], 0.5) = %f, want 42.0", result)
	}
}

func TestPercentile_Various(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// P0 should be first element
	p0 := percentile(data, 0.0)
	if p0 != 1.0 {
		t.Errorf("P0 = %f, want 1.0", p0)
	}

	// P100 should be last element (approximately)
	p100 := percentile(data, 1.0)
	if p100 != 10.0 {
		t.Errorf("P100 = %f, want 10.0", p100)
	}
}

// =============================================================================
// Collector Tests
// =============================================================================

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	if c == nil {
		t.Fatal("NewCollector() returned nil")
	}
}

func TestCollector_Counter(t *testing.T) {
	c := NewCollector()

	c.Counter("requests", 1)
	c.Counter("requests", 1)
	c.Counter("errors", 1)

	requests := c.GetCounter("requests")
	if requests != 2 {
		t.Errorf("requests counter = %d, want 2", requests)
	}

	errors := c.GetCounter("errors")
	if errors != 1 {
		t.Errorf("errors counter = %d, want 1", errors)
	}
}

func TestCollector_CounterDelta(t *testing.T) {
	c := NewCollector()

	c.Counter("requests", 5)
	c.Counter("requests", 10)

	requests := c.GetCounter("requests")
	if requests != 15 {
		t.Errorf("requests counter = %d, want 15", requests)
	}
}

func TestCollector_Gauge(t *testing.T) {
	c := NewCollector()

	c.Gauge("memory_bytes", 1024*1024)
	c.Gauge("connections", 50)

	mem := c.GetGauge("memory_bytes")
	if mem != 1024*1024 {
		t.Errorf("memory_bytes gauge = %d, want %d", mem, 1024*1024)
	}

	conn := c.GetGauge("connections")
	if conn != 50 {
		t.Errorf("connections gauge = %d, want 50", conn)
	}
}

func TestCollector_GaugeOverwrite(t *testing.T) {
	c := NewCollector()

	c.Gauge("memory", 100)
	c.Gauge("memory", 200)
	c.Gauge("memory", 50)

	mem := c.GetGauge("memory")
	if mem != 50 {
		t.Errorf("memory gauge = %d, want 50 (last set value)", mem)
	}
}

func TestCollector_Histogram(t *testing.T) {
	c := NewCollector()

	c.Histogram("query_latency", 10.5)
	c.Histogram("query_latency", 20.5)
	c.Histogram("insert_latency", 5.0)

	// Get query histogram
	queryStats := c.GetHistogram("query_latency")
	if queryStats == nil {
		t.Fatal("GetHistogram('query_latency') returned nil")
	}

	if queryStats.Count != 2 {
		t.Errorf("Query count = %d, want 2", queryStats.Count)
	}
}

func TestCollector_GetCounterNotExists(t *testing.T) {
	c := NewCollector()

	val := c.GetCounter("nonexistent")
	if val != 0 {
		t.Errorf("GetCounter for nonexistent key = %d, want 0", val)
	}
}

func TestCollector_GetGaugeNotExists(t *testing.T) {
	c := NewCollector()

	val := c.GetGauge("nonexistent")
	if val != 0 {
		t.Errorf("GetGauge for nonexistent key = %d, want 0", val)
	}
}

func TestCollector_GetHistogramNotExists(t *testing.T) {
	c := NewCollector()

	stats := c.GetHistogram("nonexistent")
	if stats != nil {
		t.Error("GetHistogram for nonexistent key should be nil")
	}
}

func TestCollector_Snapshot(t *testing.T) {
	c := NewCollector()

	c.Counter("requests", 100)
	c.Gauge("connections", 10)
	c.Histogram("latency", 25.0)

	snap := c.Snapshot()
	if snap == nil {
		t.Fatal("Snapshot() returned nil")
	}

	if snap.Counters["requests"] != 100 {
		t.Errorf("Snapshot counters[requests] = %d, want 100", snap.Counters["requests"])
	}
	if snap.Gauges["connections"] != 10 {
		t.Errorf("Snapshot gauges[connections] = %d, want 10", snap.Gauges["connections"])
	}
	if snap.Histograms["latency"] == nil {
		t.Error("Snapshot histograms[latency] should not be nil")
	}
	if snap.Uptime <= 0 {
		t.Error("Snapshot uptime should be positive")
	}
}

func TestCollector_Concurrent(t *testing.T) {
	c := NewCollector()

	var wg sync.WaitGroup
	const n = 1000

	// Concurrent counter increments
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Counter("requests", 1)
		}()
	}

	// Concurrent latency recording
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			c.Histogram("query", val)
		}(float64(i))
	}

	// Concurrent gauge updates
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int64) {
			defer wg.Done()
			c.Gauge("active", val)
		}(int64(i))
	}

	wg.Wait()

	requests := c.GetCounter("requests")
	if requests != n {
		t.Errorf("requests counter = %d, want %d", requests, n)
	}
}

// =============================================================================
// Profiler Tests
// =============================================================================

func TestNewProfiler(t *testing.T) {
	c := NewCollector()
	p := NewProfiler(c)
	if p == nil {
		t.Fatal("NewProfiler() returned nil")
	}
}

func TestProfiler_StartStop(t *testing.T) {
	c := NewCollector()
	p := NewProfiler(c)

	// Start should not panic
	p.Start()

	// Give it time to sample
	time.Sleep(50 * time.Millisecond)

	// Stop should not panic
	p.Stop()
}

func TestProfiler_Timer(t *testing.T) {
	c := NewCollector()
	p := NewProfiler(c)

	timer := p.NewTimer("operation")
	if timer == nil {
		t.Fatal("NewTimer() returned nil")
	}

	// Do some work
	time.Sleep(10 * time.Millisecond)

	duration := timer.Stop()
	if duration < 10*time.Millisecond {
		t.Errorf("Duration = %v, expected at least 10ms", duration)
	}
}

func TestProfiler_TimerRecordsToHistogram(t *testing.T) {
	c := NewCollector()
	p := NewProfiler(c)

	timer := p.NewTimer("test_op")
	time.Sleep(1 * time.Millisecond)
	timer.Stop()

	// Check that the histogram was recorded
	stats := c.GetHistogram("test_op")
	if stats == nil {
		t.Fatal("Timer did not record to histogram")
	}
	if stats.Count != 1 {
		t.Errorf("Histogram count = %d, want 1", stats.Count)
	}
}

func TestProfiler_MultipleTimers(t *testing.T) {
	c := NewCollector()
	p := NewProfiler(c)

	for i := 0; i < 5; i++ {
		timer := p.NewTimer("batch_op")
		time.Sleep(1 * time.Millisecond)
		timer.Stop()
	}

	stats := c.GetHistogram("batch_op")
	if stats == nil {
		t.Fatal("Histogram not found")
	}
	if stats.Count != 5 {
		t.Errorf("Histogram count = %d, want 5", stats.Count)
	}
}

// =============================================================================
// Snapshot Tests
// =============================================================================

func TestSnapshot_Structure(t *testing.T) {
	snap := &Snapshot{
		Timestamp:  time.Now(),
		Uptime:     5 * time.Minute,
		Counters:   make(map[string]int64),
		Gauges:     make(map[string]int64),
		Histograms: make(map[string]*HistogramStats),
	}

	snap.Counters["test"] = 100
	snap.Gauges["test"] = 50
	snap.Histograms["test"] = &HistogramStats{Count: 10}

	if snap.Counters["test"] != 100 {
		t.Error("Snapshot counters not set correctly")
	}
	if snap.Uptime != 5*time.Minute {
		t.Error("Snapshot uptime not set correctly")
	}
}

// =============================================================================
// Collector Reset Tests
// =============================================================================

func TestCollector_Reset(t *testing.T) {
	c := NewCollector()

	c.Counter("requests", 100)
	c.Gauge("connections", 50)
	c.Histogram("latency", 10.0)

	// Verify data exists
	if c.GetCounter("requests") != 100 {
		t.Error("Counter not set before reset")
	}

	// Reset
	c.Reset()

	// Verify data is cleared
	if c.GetCounter("requests") != 0 {
		t.Error("Counter not cleared after reset")
	}
	if c.GetGauge("connections") != 0 {
		t.Error("Gauge not cleared after reset")
	}
	if c.GetHistogram("latency") != nil {
		t.Error("Histogram not cleared after reset")
	}
}

// =============================================================================
// FormatStats Tests
// =============================================================================

func TestFormatStats(t *testing.T) {
	c := NewCollector()
	c.Counter("requests", 100)
	c.Gauge("memory", 1024)
	c.Histogram("latency", 10.0)

	snap := c.Snapshot()
	output := FormatStats(snap)

	if output == "" {
		t.Error("FormatStats returned empty string")
	}

	// Should contain uptime info
	if len(output) < 10 {
		t.Error("FormatStats output too short")
	}
}

func TestFormatStats_Empty(t *testing.T) {
	c := NewCollector()
	snap := c.Snapshot()
	output := FormatStats(snap)

	if output == "" {
		t.Error("FormatStats with empty snapshot returned empty string")
	}
}

func TestFormatStats_Content(t *testing.T) {
	c := NewCollector()
	c.Counter("test_counter", 42)
	c.Gauge("test_gauge", 100)
	c.Histogram("test_hist", 5.0)

	snap := c.Snapshot()
	output := FormatStats(snap)

	// Should contain section headers
	if !containsString(output, "Counters:") {
		t.Error("FormatStats missing Counters section")
	}
	if !containsString(output, "Gauges:") {
		t.Error("FormatStats missing Gauges section")
	}
	if !containsString(output, "Histograms:") {
		t.Error("FormatStats missing Histograms section")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestHistogram_VerySmallValues(t *testing.T) {
	h := NewHistogram()

	h.Record(0.000001)
	h.Record(0.000002)
	h.Record(0.000003)

	stats := h.Stats()
	if stats.Min > 0.00001 {
		t.Errorf("Min not handling small values: %f", stats.Min)
	}
}

func TestHistogram_VeryLargeValues(t *testing.T) {
	h := NewHistogram()

	h.Record(1e15)
	h.Record(1e16)

	stats := h.Stats()
	if stats.Max < 1e15 {
		t.Error("Max not handling large values")
	}
}

func TestHistogram_MixedScale(t *testing.T) {
	h := NewHistogram()

	h.Record(0.001)
	h.Record(1000000.0)

	stats := h.Stats()
	if stats.Min != 0.001 {
		t.Errorf("Min = %f, want 0.001", stats.Min)
	}
	if stats.Max != 1000000.0 {
		t.Errorf("Max = %f, want 1000000.0", stats.Max)
	}
}

func TestHistogram_InfNaN(t *testing.T) {
	h := NewHistogram()

	// These should be handled gracefully (not panic)
	h.Record(math.Inf(1))
	h.Record(math.Inf(-1))
	h.Record(math.NaN())

	stats := h.Stats()
	// Just verify it doesn't panic and returns something
	_ = stats
}
