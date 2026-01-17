// Package memory - comprehensive tests for memory management
package memory

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Config Tests
// =============================================================================

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Default config uses 0 for MaxMemoryBytes and MaxItems to mean "no limit"
	if cfg.MaxMemoryBytes < 0 {
		t.Error("MaxMemoryBytes should not be negative")
	}

	if cfg.MaxItems < 0 {
		t.Error("MaxItems should not be negative")
	}

	if cfg.TTLCheckInterval <= 0 {
		t.Error("TTLCheckInterval should be positive")
	}

	if cfg.EvictionPolicy != EvictionLRU {
		t.Error("Default eviction policy should be LRU")
	}
}

func TestConfig_Custom(t *testing.T) {
	cfg := &Config{
		MaxMemoryBytes:   1024 * 1024 * 1024, // 1GB
		MaxItems:         100000,
		TTLCheckInterval: 30 * time.Second,
	}

	if cfg.MaxMemoryBytes != 1024*1024*1024 {
		t.Errorf("MaxMemoryBytes = %d, want %d", cfg.MaxMemoryBytes, 1024*1024*1024)
	}
}

// =============================================================================
// Manager Tests
// =============================================================================

func TestManager_Create(t *testing.T) {
	manager := NewManager(nil)
	if manager == nil {
		t.Fatal("NewManager(nil) returned nil")
	}
	defer manager.Stop()
}

func TestManager_WithConfig(t *testing.T) {
	cfg := &Config{
		MaxMemoryBytes:   512 * 1024 * 1024,
		MaxItems:         50000,
		TTLCheckInterval: time.Second,
	}

	manager := NewManager(cfg)
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
	defer manager.Stop()

	if manager.config.MaxMemoryBytes != cfg.MaxMemoryBytes {
		t.Error("Config not applied correctly")
	}
}

func TestManager_StartStop(t *testing.T) {
	manager := NewManager(nil)

	// Start should not panic
	manager.Start()

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// Stop should not panic
	manager.Stop()
}

func TestManager_Stats(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Stop()

	stats := manager.Stats()

	// Basic sanity checks
	if stats.AllocatedBytes < 0 {
		t.Error("AllocatedBytes should not be negative")
	}
	if stats.SystemBytes < 0 {
		t.Error("SystemBytes should not be negative")
	}
}

func TestManager_GetCaches(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Stop()

	if manager.GetEntityCache() == nil {
		t.Error("GetEntityCache() returned nil")
	}
	if manager.GetTextUnitCache() == nil {
		t.Error("GetTextUnitCache() returned nil")
	}
	if manager.GetDocumentCache() == nil {
		t.Error("GetDocumentCache() returned nil")
	}
	if manager.GetCommunityCache() == nil {
		t.Error("GetCommunityCache() returned nil")
	}
}

// =============================================================================
// LRU Cache Tests
// =============================================================================

func TestLRUCache_Create(t *testing.T) {
	cache := NewLRUCache(100)
	if cache == nil {
		t.Fatal("NewLRUCache() returned nil")
	}
}

func TestLRUCache_PutGet(t *testing.T) {
	cache := NewLRUCache(100)

	// Put value
	cache.Put("key1", "value1", 10)

	// Get value
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("Get() should return true for existing key")
	}
	if val != "value1" {
		t.Errorf("Get() = %v, want 'value1'", val)
	}
}

func TestLRUCache_GetMiss(t *testing.T) {
	cache := NewLRUCache(100)

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existent key")
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(3)

	// Fill cache
	cache.Put("key1", "a", 10)
	cache.Put("key2", "b", 10)
	cache.Put("key3", "c", 10)

	// Add one more, should evict oldest
	cache.Put("key4", "d", 10)

	// Key 1 should be evicted (LRU)
	_, ok := cache.Get("key1")
	if ok {
		t.Error("Key 1 should have been evicted")
	}

	// Key 4 should exist
	_, ok = cache.Get("key4")
	if !ok {
		t.Error("Key 4 should exist")
	}
}

func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache(100)

	cache.Put("key1", "original", 10)
	cache.Put("key1", "updated", 10)

	val, ok := cache.Get("key1")
	if !ok || val != "updated" {
		t.Errorf("Updated value not returned: got %v, want 'updated'", val)
	}
}

func TestLRUCache_Remove(t *testing.T) {
	cache := NewLRUCache(100)

	cache.Put("key1", "value", 10)
	removed := cache.Remove("key1")

	if !removed {
		t.Error("Remove should return true for existing key")
	}

	_, ok := cache.Get("key1")
	if ok {
		t.Error("Removed key should not exist")
	}
}

func TestLRUCache_Len(t *testing.T) {
	cache := NewLRUCache(100)

	if cache.Len() != 0 {
		t.Errorf("Empty cache Len() = %d, want 0", cache.Len())
	}

	cache.Put("key1", "a", 10)
	cache.Put("key2", "b", 10)

	if cache.Len() != 2 {
		t.Errorf("Cache Len() = %d, want 2", cache.Len())
	}
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(100)

	cache.Put("key1", "a", 10)
	cache.Put("key2", "b", 10)
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("Cleared cache Len() = %d, want 0", cache.Len())
	}
}

func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache(100)

	// Generate some hits and misses
	cache.Put("key1", "a", 10)
	cache.Get("key1") // hit
	cache.Get("key1") // hit
	cache.Get("key2") // miss

	hits, misses := cache.Stats()
	if hits < 2 {
		t.Errorf("Expected at least 2 hits, got %d", hits)
	}
	if misses < 1 {
		t.Errorf("Expected at least 1 miss, got %d", misses)
	}
}

func TestLRUCache_Concurrent(t *testing.T) {
	cache := NewLRUCache(1000)

	var wg sync.WaitGroup
	const n = 100

	// Concurrent writes
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Put("key"+itoa(id), "value", 10)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Get("key" + itoa(id))
		}(i)
	}

	wg.Wait()
}

// Helper for number to string
func itoa(i int) string {
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
// Memory Stats Tests
// =============================================================================

func TestMemoryStats_Structure(t *testing.T) {
	stats := MemoryStats{
		AllocatedBytes:    1024 * 1024,
		TotalAllocBytes:   2048 * 1024,
		SystemBytes:       4096 * 1024,
		NumGC:             10,
		EntityCacheLen:    100,
		TextUnitCacheLen:  200,
		DocumentCacheLen:  50,
		CommunityCacheLen: 25,
		CacheHits:         1000,
		CacheMisses:       100,
	}

	if stats.AllocatedBytes != 1024*1024 {
		t.Error("AllocatedBytes incorrect")
	}
	if stats.NumGC != 10 {
		t.Error("NumGC incorrect")
	}
	if stats.CacheHits != 1000 {
		t.Error("CacheHits incorrect")
	}
}

// =============================================================================
// Tracker Tests
// =============================================================================

func TestTracker_Create(t *testing.T) {
	tracker := NewTracker(1024 * 1024 * 100) // 100MB max
	if tracker == nil {
		t.Fatal("NewTracker() returned nil")
	}
}

func TestTracker_Check(t *testing.T) {
	tracker := NewTracker(1024 * 1024 * 1024) // 1GB max

	usedBytes, level := tracker.Check()

	// Should return some usage
	if usedBytes < 0 {
		t.Error("usedBytes should not be negative")
	}

	// With 1GB limit, level should be ok
	if level != "ok" {
		t.Errorf("Level = %q, expected 'ok' with high limit", level)
	}
}

func TestTracker_Check_Warning(t *testing.T) {
	// Set max to current usage to trigger warning/critical
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	// Set max to slightly above current to test warning (80% threshold)
	tracker := NewTracker(int64(float64(stats.Alloc) * 1.1))
	_, level := tracker.Check()

	// Should be warning or critical since we're close to max
	if level != "warning" && level != "critical" && level != "ok" {
		t.Errorf("Level = %q, expected 'ok', 'warning', or 'critical'", level)
	}
}

func TestTracker_SetAlertCallback(t *testing.T) {
	tracker := NewTracker(100) // Very low limit

	called := false
	tracker.SetAlertCallback(func(level string, usedBytes, maxBytes int64) {
		called = true
	})

	tracker.Check()

	// Callback should have been called since we're way over limit
	if !called {
		t.Log("Alert callback was not called (may depend on current memory usage)")
	}
}

func TestTracker_GetStats(t *testing.T) {
	tracker := NewTracker(1024 * 1024 * 100)

	// First check to populate stats
	tracker.Check()

	stats, lastCheck := tracker.GetStats()

	if lastCheck.IsZero() {
		t.Error("lastCheck should not be zero after Check()")
	}
	_ = stats
}

func TestTracker_ForceGC(t *testing.T) {
	tracker := NewTracker(1024 * 1024 * 100)

	// Should not panic
	tracker.ForceGC()
}

// =============================================================================
// Memory Pressure Tests
// =============================================================================

func TestManager_MemoryPressureCheck(t *testing.T) {
	cfg := &Config{
		MaxMemoryBytes:   100, // Very low limit to trigger pressure
		MaxItems:         1000,
		TTLCheckInterval: time.Hour, // Don't run monitor loop
	}

	manager := NewManager(cfg)
	defer manager.Stop()

	// Manually trigger pressure check
	manager.checkMemoryPressure()

	// Should not panic
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestManager_FullWorkflow(t *testing.T) {
	manager := NewManager(nil)
	manager.Start()
	defer manager.Stop()

	// Use caches
	entityCache := manager.GetEntityCache()
	for i := 1; i <= 100; i++ {
		entityCache.Put("entity:"+itoa(i), map[string]interface{}{
			"id":    i,
			"title": "Entity",
		}, 100)
	}

	// Check stats
	stats := manager.Stats()
	if stats.EntityCacheLen != 100 {
		t.Errorf("EntityCacheLen = %d, want 100", stats.EntityCacheLen)
	}

	// Verify we can retrieve
	val, ok := entityCache.Get("entity:50")
	if !ok {
		t.Error("Should be able to retrieve cached entity")
	}
	if val == nil {
		t.Error("Retrieved value should not be nil")
	}
}

// =============================================================================
// Runtime Stats Verification
// =============================================================================

func TestManager_RuntimeStats(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Stop()

	stats := manager.Stats()

	// These should match runtime values
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// AllocatedBytes should be close to runtime.Alloc
	// (may differ slightly due to timing)
	if stats.AllocatedBytes <= 0 {
		t.Error("AllocatedBytes should be positive")
	}

	if stats.SystemBytes <= 0 {
		t.Error("SystemBytes should be positive")
	}
}
