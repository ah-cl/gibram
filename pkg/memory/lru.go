// Package memory provides memory management for GibRAM
package memory

import (
	"container/list"
	"sync"
)

// LRUCache is a thread-safe LRU cache
type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
	mu       sync.RWMutex

	// Stats
	hits   int64
	misses int64
}

type lruEntry struct {
	key   string
	value interface{}
	size  int64
}

// NewLRUCache creates a new LRU cache with given capacity
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves an item from cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		c.hits++
		return elem.Value.(*lruEntry).value, true
	}

	c.misses++
	return nil, false
}

// Put adds an item to cache
func (c *LRUCache) Put(key string, value interface{}, size int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		entry := elem.Value.(*lruEntry)
		entry.value = value
		entry.size = size
		return
	}

	// Evict if at capacity
	if c.capacity > 0 && c.order.Len() >= c.capacity {
		c.evictOldest()
	}

	// Add new entry
	entry := &lruEntry{key: key, value: value, size: size}
	elem := c.order.PushFront(entry)
	c.items[key] = elem
}

// Remove removes an item from cache
func (c *LRUCache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

// Len returns number of items in cache
func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.order.Len()
}

// Clear removes all items from cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.order.Init()
}

// Stats returns cache hit/miss statistics
func (c *LRUCache) Stats() (hits, misses int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}

// EvictOldest evicts up to count oldest items from the cache
// Returns the actual number of items evicted
func (c *LRUCache) EvictOldest(count int) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	evicted := 0
	for i := 0; i < count && c.order.Len() > 0; i++ {
		c.evictOldest()
		evicted++
	}
	return evicted
}

func (c *LRUCache) evictOldest() {
	oldest := c.order.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

func (c *LRUCache) removeElement(elem *list.Element) {
	c.order.Remove(elem)
	entry := elem.Value.(*lruEntry)
	delete(c.items, entry.key)
}
