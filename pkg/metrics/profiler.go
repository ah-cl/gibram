// Package metrics provides metrics collection for GibRAM
package metrics

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Profiler provides performance profiling
type Profiler struct {
	collector *Collector
	stopCh    chan struct{}
	wg        sync.WaitGroup

	// Configuration
	sampleInterval time.Duration
}

// NewProfiler creates a new profiler
func NewProfiler(collector *Collector) *Profiler {
	return &Profiler{
		collector:      collector,
		stopCh:         make(chan struct{}),
		sampleInterval: 10 * time.Second,
	}
}

// Start starts the profiler
func (p *Profiler) Start() {
	p.wg.Add(1)
	go p.sampleLoop()
}

// Stop stops the profiler
func (p *Profiler) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *Profiler) sampleLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.sampleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.sample()
		}
	}
}

func (p *Profiler) sample() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memory metrics
	p.collector.Gauge("memory.alloc_bytes", int64(m.Alloc))
	p.collector.Gauge("memory.total_alloc_bytes", int64(m.TotalAlloc))
	p.collector.Gauge("memory.sys_bytes", int64(m.Sys))
	p.collector.Gauge("memory.heap_objects", int64(m.HeapObjects))
	p.collector.Gauge("memory.num_gc", int64(m.NumGC))

	// Goroutine count
	p.collector.Gauge("goroutines", int64(runtime.NumGoroutine()))
}

// Timer is a helper for timing operations
type Timer struct {
	collector *Collector
	name      string
	start     time.Time
}

// NewTimer creates a new timer
func (p *Profiler) NewTimer(name string) *Timer {
	return &Timer{
		collector: p.collector,
		name:      name,
		start:     time.Now(),
	}
}

// Stop stops the timer and records the duration
func (t *Timer) Stop() time.Duration {
	duration := time.Since(t.start)
	t.collector.Histogram(t.name, float64(duration.Microseconds()))
	return duration
}

// FormatStats formats stats for display
func FormatStats(snap *Snapshot) string {
	s := fmt.Sprintf("Uptime: %s\n", snap.Uptime.Round(time.Second))
	s += "\nCounters:\n"
	for k, v := range snap.Counters {
		s += fmt.Sprintf("  %s: %d\n", k, v)
	}
	s += "\nGauges:\n"
	for k, v := range snap.Gauges {
		s += fmt.Sprintf("  %s: %d\n", k, v)
	}
	s += "\nHistograms:\n"
	for k, h := range snap.Histograms {
		s += fmt.Sprintf("  %s: count=%d avg=%.2f p99=%.2f\n", k, h.Count, h.Avg, h.P99)
	}
	return s
}
