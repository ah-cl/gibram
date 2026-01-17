// Package memory provides memory management for GibRAM
package memory

import (
	"runtime"
	"sync"
	"time"
)

// Tracker tracks memory usage and provides alerts
type Tracker struct {
	maxBytes     int64
	warningBytes int64

	mu            sync.RWMutex
	lastCheck     time.Time
	lastStats     runtime.MemStats
	alertCallback func(level string, usedBytes, maxBytes int64)
}

// NewTracker creates a new memory tracker
func NewTracker(maxBytes int64) *Tracker {
	return &Tracker{
		maxBytes:     maxBytes,
		warningBytes: int64(float64(maxBytes) * 0.8), // 80% warning threshold
	}
}

// SetAlertCallback sets the callback for memory alerts
func (t *Tracker) SetAlertCallback(cb func(level string, usedBytes, maxBytes int64)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.alertCallback = cb
}

// Check checks current memory usage
func (t *Tracker) Check() (usedBytes int64, level string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	t.mu.Lock()
	t.lastCheck = time.Now()
	t.lastStats = stats
	cb := t.alertCallback
	t.mu.Unlock()

	usedBytes = int64(stats.Alloc)

	if t.maxBytes > 0 {
		if usedBytes >= t.maxBytes {
			level = "critical"
		} else if usedBytes >= t.warningBytes {
			level = "warning"
		} else {
			level = "ok"
		}

		if cb != nil && level != "ok" {
			cb(level, usedBytes, t.maxBytes)
		}
	} else {
		level = "ok"
	}

	return usedBytes, level
}

// GetStats returns last memory stats
func (t *Tracker) GetStats() (stats runtime.MemStats, lastCheck time.Time) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lastStats, t.lastCheck
}

// ForceGC forces garbage collection
func (t *Tracker) ForceGC() {
	runtime.GC()
}
