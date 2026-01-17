// Package resilience provides reliability patterns for GibRAM
package resilience

import (
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// StateClosed allows requests to pass through
	StateClosed CircuitState = iota

	// StateOpen blocks requests and returns errors immediately
	StateOpen

	// StateHalfOpen allows limited requests to test recovery
	StateHalfOpen
)

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu sync.RWMutex

	// Configuration
	maxFailures     uint32        // failures before opening
	timeout         time.Duration // how long to stay open
	halfOpenMaxReqs uint32        // max concurrent requests in half-open

	// State
	state            CircuitState
	failures         uint32
	lastFailureTime  time.Time
	lastStateChange  time.Time
	halfOpenRequests uint32

	// Statistics
	totalRequests uint64
	totalFailures uint64
	totalSuccesses uint64
}

// Config holds circuit breaker configuration
type Config struct {
	MaxFailures     uint32        // failures before opening (default: 5)
	Timeout         time.Duration // open state duration (default: 60s)
	HalfOpenMaxReqs uint32        // max requests in half-open (default: 1)
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	if cfg.MaxFailures == 0 {
		cfg.MaxFailures = 5
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60 * time.Second
	}
	if cfg.HalfOpenMaxReqs == 0 {
		cfg.HalfOpenMaxReqs = 1
	}

	return &CircuitBreaker{
		maxFailures:     cfg.MaxFailures,
		timeout:         cfg.Timeout,
		halfOpenMaxReqs: cfg.HalfOpenMaxReqs,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// Check if we can execute
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// Execute function
	err := fn()

	// Record result
	cb.afterRequest(err)

	return err
}

// beforeRequest checks if the request can proceed
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalRequests++

	switch cb.state {
	case StateClosed:
		// Allow request
		return nil

	case StateOpen:
		// Check if timeout has elapsed
		if time.Since(cb.lastStateChange) > cb.timeout {
			// Transition to half-open
			cb.state = StateHalfOpen
			cb.lastStateChange = time.Now()
			cb.halfOpenRequests = 1
			return nil
		}
		// Still open, reject request
		return ErrCircuitOpen

	case StateHalfOpen:
		// Check if we can allow more requests
		if cb.halfOpenRequests >= cb.halfOpenMaxReqs {
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
		return nil

	default:
		return ErrCircuitOpen
	}
}

// afterRequest records the result and updates state
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.totalFailures++
		cb.onFailure()
	} else {
		cb.totalSuccesses++
		cb.onSuccess()
	}
}

// onFailure handles a failed request
func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.maxFailures {
			// Open the circuit
			cb.state = StateOpen
			cb.lastStateChange = time.Now()
		}

	case StateHalfOpen:
		// Failure in half-open means we're not recovered
		// Go back to open state
		cb.state = StateOpen
		cb.lastStateChange = time.Now()
		cb.halfOpenRequests = 0
	}
}

// onSuccess handles a successful request
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		// Reset failure counter on success
		cb.failures = 0

	case StateHalfOpen:
		// Success in half-open means we can close the circuit
		cb.state = StateClosed
		cb.failures = 0
		cb.halfOpenRequests = 0
		cb.lastStateChange = time.Now()
	}
}

// GetState returns the current state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats returns circuit breaker statistics
func (cb *CircuitBreaker) GetStats() Stats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return Stats{
		State:           cb.state,
		Failures:        cb.failures,
		TotalRequests:   cb.totalRequests,
		TotalFailures:   cb.totalFailures,
		TotalSuccesses:  cb.totalSuccesses,
		LastFailureTime: cb.lastFailureTime,
		LastStateChange: cb.lastStateChange,
	}
}

// Stats holds circuit breaker statistics
type Stats struct {
	State           CircuitState
	Failures        uint32
	TotalRequests   uint64
	TotalFailures   uint64
	TotalSuccesses  uint64
	LastFailureTime time.Time
	LastStateChange time.Time
}

// Reset manually resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenRequests = 0
	cb.lastStateChange = time.Now()
}

// String returns a string representation of the state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}
