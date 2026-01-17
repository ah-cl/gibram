// Package memory provides memory management for GibRAM
package memory

import "time"

// Config holds memory manager configuration
type Config struct {
	// MaxMemoryBytes is the maximum memory to use (0 = no limit)
	MaxMemoryBytes int64

	// MaxItems is the maximum number of items to cache (0 = no limit)
	MaxItems int

	// EvictionPolicy determines how items are evicted
	EvictionPolicy EvictionPolicy

	// TTLCheckInterval is how often to check for expired items
	TTLCheckInterval time.Duration

	// EnableMetrics enables memory usage metrics
	EnableMetrics bool
}

// EvictionPolicy defines the cache eviction strategy
type EvictionPolicy int

const (
	// EvictionLRU evicts least recently used items
	EvictionLRU EvictionPolicy = iota

	// EvictionLFU evicts least frequently used items
	EvictionLFU

	// EvictionFIFO evicts oldest items first
	EvictionFIFO
)

// DefaultConfig returns default memory configuration
func DefaultConfig() *Config {
	return &Config{
		MaxMemoryBytes:   0, // No limit
		MaxItems:         0, // No limit
		EvictionPolicy:   EvictionLRU,
		TTLCheckInterval: 60 * time.Second,
		EnableMetrics:    true,
	}
}
