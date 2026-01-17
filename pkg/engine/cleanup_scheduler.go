// Package engine - Heap-based session cleanup scheduler
package engine

import (
	"container/heap"
	"sync"
	"time"
)

// sessionExpiry represents a session's expiration time
type sessionExpiry struct {
	sessionID string
	expireAt  int64 // nanoseconds
	index     int   // index in heap
}

// expiryHeap implements heap.Interface for session expirations
type expiryHeap []*sessionExpiry

func (h expiryHeap) Len() int           { return len(h) }
func (h expiryHeap) Less(i, j int) bool { return h[i].expireAt < h[j].expireAt }
func (h expiryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *expiryHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*sessionExpiry)
	item.index = n
	*h = append(*h, item)
}

func (h *expiryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

// SessionCleanupScheduler uses a min-heap to efficiently track and cleanup expired sessions
type SessionCleanupScheduler struct {
	mu         sync.Mutex
	heap       expiryHeap
	heapIndex  map[string]*sessionExpiry // fast lookup by sessionID
	engine     *Engine
	stopChan   chan struct{}
	wg         sync.WaitGroup
	minDelay   time.Duration // minimum delay between checks (prevents tight loops)
}

// NewSessionCleanupScheduler creates a new heap-based cleanup scheduler
func NewSessionCleanupScheduler(engine *Engine) *SessionCleanupScheduler {
	s := &SessionCleanupScheduler{
		heap:      make(expiryHeap, 0),
		heapIndex: make(map[string]*sessionExpiry),
		engine:    engine,
		stopChan:  make(chan struct{}),
		minDelay:  100 * time.Millisecond, // avoid checking too frequently
	}
	heap.Init(&s.heap)
	return s
}

// Start starts the cleanup scheduler
func (s *SessionCleanupScheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

// run is the main scheduler loop
func (s *SessionCleanupScheduler) run() {
	defer s.wg.Done()

	timer := time.NewTimer(time.Hour) // will be reset
	defer timer.Stop()

	for {
		// Check next expiration time
		delay := s.getNextDelay()
		timer.Reset(delay)

		select {
		case <-timer.C:
			s.cleanup()
		case <-s.stopChan:
			return
		}
	}
}

// getNextDelay returns the time until the next session expires
func (s *SessionCleanupScheduler) getNextDelay() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.heap.Len() == 0 {
		return time.Hour // check back in an hour if no sessions
	}

	// Get time until next expiration
	next := s.heap[0].expireAt
	now := time.Now().UnixNano()
	delay := time.Duration(next - now)

	if delay < s.minDelay {
		delay = s.minDelay
	}

	return delay
}

// cleanup removes expired sessions
func (s *SessionCleanupScheduler) cleanup() {
	now := time.Now().UnixNano()

	s.mu.Lock()
	var toRemove []string

	// Collect all expired sessions
	for s.heap.Len() > 0 && s.heap[0].expireAt <= now {
		expiry := heap.Pop(&s.heap).(*sessionExpiry)
		toRemove = append(toRemove, expiry.sessionID)
		delete(s.heapIndex, expiry.sessionID)
	}
	s.mu.Unlock()

	// Remove from engine
	if len(toRemove) > 0 {
		s.engine.mu.Lock()
		for _, sessionID := range toRemove {
			// Re-check expiry in case session was touched
			if sess, ok := s.engine.sessions[sessionID]; ok && sess.IsExpired() {
				delete(s.engine.sessions, sessionID)
			}
		}
		s.engine.mu.Unlock()
	}
}

// UpdateSession updates or adds a session's expiration time
func (s *SessionCleanupScheduler) UpdateSession(sessionID string, expireAt int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If session already in heap, update it
	if existing, ok := s.heapIndex[sessionID]; ok {
		existing.expireAt = expireAt
		heap.Fix(&s.heap, existing.index)
		return
	}

	// Add new session
	expiry := &sessionExpiry{
		sessionID: sessionID,
		expireAt:  expireAt,
	}
	heap.Push(&s.heap, expiry)
	s.heapIndex[sessionID] = expiry
}

// RemoveSession removes a session from the scheduler (e.g., when manually deleted)
func (s *SessionCleanupScheduler) RemoveSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if expiry, ok := s.heapIndex[sessionID]; ok {
		heap.Remove(&s.heap, expiry.index)
		delete(s.heapIndex, sessionID)
	}
}

// Stop stops the cleanup scheduler
func (s *SessionCleanupScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// GetStats returns scheduler statistics
func (s *SessionCleanupScheduler) GetStats() SchedulerStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	return SchedulerStats{
		TrackedSessions: s.heap.Len(),
		NextExpiration:  s.getNextExpirationUnlocked(),
	}
}

// getNextExpirationUnlocked returns the next expiration time (must be called with lock held)
func (s *SessionCleanupScheduler) getNextExpirationUnlocked() int64 {
	if s.heap.Len() == 0 {
		return 0
	}
	return s.heap[0].expireAt
}

// SchedulerStats holds cleanup scheduler statistics
type SchedulerStats struct {
	TrackedSessions int
	NextExpiration  int64 // nanoseconds, 0 if no sessions
}
