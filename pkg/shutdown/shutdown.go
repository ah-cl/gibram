// Package shutdown provides graceful shutdown handling for GibRAM
package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Handler manages graceful shutdown
type Handler struct {
	hooks    []ShutdownHook
	mu       sync.Mutex
	timeout  time.Duration
	signals  []os.Signal
	done     chan struct{}
	started  bool
}

// ShutdownHook is a function called during shutdown
type ShutdownHook struct {
	Name     string
	Priority int // Lower priority runs first
	Fn       func(ctx context.Context) error
}

// NewHandler creates a new shutdown handler
func NewHandler() *Handler {
	return &Handler{
		hooks:   make([]ShutdownHook, 0),
		timeout: 30 * time.Second,
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		done:    make(chan struct{}),
	}
}

// SetTimeout sets the shutdown timeout
func (h *Handler) SetTimeout(d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = d
}

// SetSignals sets which signals trigger shutdown
func (h *Handler) SetSignals(signals ...os.Signal) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.signals = signals
}

// Register registers a shutdown hook
func (h *Handler) Register(name string, priority int, fn func(ctx context.Context) error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.hooks = append(h.hooks, ShutdownHook{
		Name:     name,
		Priority: priority,
		Fn:       fn,
	})

	// Sort by priority
	for i := len(h.hooks) - 1; i > 0; i-- {
		if h.hooks[i].Priority < h.hooks[i-1].Priority {
			h.hooks[i], h.hooks[i-1] = h.hooks[i-1], h.hooks[i]
		}
	}
}

// Start starts listening for shutdown signals
func (h *Handler) Start() {
	h.mu.Lock()
	if h.started {
		h.mu.Unlock()
		return
	}
	h.started = true
	h.mu.Unlock()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, h.signals...)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v, starting graceful shutdown...", sig)
		h.Shutdown()
	}()
}

// Shutdown executes all shutdown hooks
func (h *Handler) Shutdown() {
	h.mu.Lock()
	hooks := make([]ShutdownHook, len(h.hooks))
	copy(hooks, h.hooks)
	timeout := h.timeout
	h.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	errors := make(chan error, len(hooks))

	// Group hooks by priority
	priorityGroups := make(map[int][]ShutdownHook)
	for _, hook := range hooks {
		priorityGroups[hook.Priority] = append(priorityGroups[hook.Priority], hook)
	}

	// Get sorted priorities
	priorities := make([]int, 0, len(priorityGroups))
	for p := range priorityGroups {
		priorities = append(priorities, p)
	}
	for i := 0; i < len(priorities)-1; i++ {
		for j := i + 1; j < len(priorities); j++ {
			if priorities[i] > priorities[j] {
				priorities[i], priorities[j] = priorities[j], priorities[i]
			}
		}
	}

	// Execute hooks by priority group
	for _, priority := range priorities {
		group := priorityGroups[priority]

		// Run hooks in this priority group concurrently
		for _, hook := range group {
			wg.Add(1)
			go func(h ShutdownHook) {
				defer wg.Done()
				log.Printf("Shutdown: running hook '%s' (priority %d)", h.Name, h.Priority)
				if err := h.Fn(ctx); err != nil {
					log.Printf("Shutdown: hook '%s' error: %v", h.Name, err)
					errors <- err
				} else {
					log.Printf("Shutdown: hook '%s' completed", h.Name)
				}
			}(hook)
		}

		// Wait for this priority group to complete before next
		wg.Wait()
	}

	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Printf("Shutdown completed with %d errors", len(errs))
	} else {
		log.Printf("Shutdown completed successfully")
	}

	close(h.done)
}

// Wait waits for shutdown to complete
func (h *Handler) Wait() {
	<-h.done
}

// Done returns a channel that's closed when shutdown is complete
func (h *Handler) Done() <-chan struct{} {
	return h.done
}

// Default creates a default shutdown handler with common hooks
func Default() *Handler {
	h := NewHandler()
	h.SetTimeout(30 * time.Second)
	return h
}

// GracefulShutdown is a helper that creates and starts a shutdown handler
func GracefulShutdown(timeout time.Duration, hooks ...ShutdownHook) *Handler {
	h := NewHandler()
	h.SetTimeout(timeout)

	for _, hook := range hooks {
		h.Register(hook.Name, hook.Priority, hook.Fn)
	}

	h.Start()
	return h
}
