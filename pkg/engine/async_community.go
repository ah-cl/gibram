// Package engine - Async community detection support
package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gibram-io/gibram/pkg/graph"
	"github.com/gibram-io/gibram/pkg/types"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrTaskNotComplete = errors.New("task not complete")
	ErrTaskFailed      = errors.New("task failed")
)

// TaskStatus represents the state of an async task
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusComplete
	TaskStatusFailed
)

// CommunityTask represents an async community detection task
type CommunityTask struct {
	ID          string
	SessionID   string
	Status      TaskStatus
	Config      graph.LeidenConfig
	Hierarchical bool
	StartTime   time.Time
	EndTime     time.Time
	Result      []*types.Community
	Error       error
	Progress    float64 // 0.0 to 1.0
}

// CommunityTaskManager manages async community detection tasks
type CommunityTaskManager struct {
	mu      sync.RWMutex
	tasks   map[string]*CommunityTask
	engine  *Engine
	workers int // number of concurrent workers
	queue   chan *CommunityTask
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewCommunityTaskManager creates a new task manager
func NewCommunityTaskManager(engine *Engine, workers int) *CommunityTaskManager {
	if workers <= 0 {
		workers = 2
	}

	ctx, cancel := context.WithCancel(context.Background())

	tm := &CommunityTaskManager{
		tasks:   make(map[string]*CommunityTask),
		engine:  engine,
		workers: workers,
		queue:   make(chan *CommunityTask, 100),
		ctx:     ctx,
		cancel:  cancel,
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go tm.worker(i)
	}

	return tm
}

// worker processes community detection tasks
func (tm *CommunityTaskManager) worker(id int) {
	for {
		select {
		case <-tm.ctx.Done():
			return
		case task := <-tm.queue:
			tm.processTask(task)
		}
	}
}

// processTask executes a community detection task
func (tm *CommunityTaskManager) processTask(task *CommunityTask) {
	// Mark as running
	tm.mu.Lock()
	task.Status = TaskStatusRunning
	task.StartTime = time.Now()
	tm.mu.Unlock()

	var communities []*types.Community
	var err error

	// Execute community detection
	if task.Hierarchical {
		communities, err = tm.engine.ComputeHierarchicalCommunities(task.SessionID, task.Config)
	} else {
		communities, err = tm.engine.ComputeCommunities(task.SessionID, task.Config)
	}

	// Update task with result
	tm.mu.Lock()
	task.EndTime = time.Now()
	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err
	} else {
		task.Status = TaskStatusComplete
		task.Result = communities
		task.Progress = 1.0
	}
	tm.mu.Unlock()
}

// SubmitCommunityTask submits a new community detection task
func (tm *CommunityTaskManager) SubmitCommunityTask(sessionID string, config graph.LeidenConfig, hierarchical bool) (string, error) {
	taskID := fmt.Sprintf("comm_%s_%d", sessionID, time.Now().UnixNano())

	task := &CommunityTask{
		ID:           taskID,
		SessionID:    sessionID,
		Status:       TaskStatusPending,
		Config:       config,
		Hierarchical: hierarchical,
		Progress:     0.0,
	}

	tm.mu.Lock()
	tm.tasks[taskID] = task
	tm.mu.Unlock()

	// Queue task for processing
	select {
	case tm.queue <- task:
		return taskID, nil
	case <-tm.ctx.Done():
		return "", errors.New("task manager shutting down")
	}
}

// GetTaskStatus returns the status of a task
func (tm *CommunityTaskManager) GetTaskStatus(taskID string) (*CommunityTask, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	task, ok := tm.tasks[taskID]
	if !ok {
		return nil, ErrTaskNotFound
	}

	// Return a copy to avoid race conditions
	return &CommunityTask{
		ID:           task.ID,
		SessionID:    task.SessionID,
		Status:       task.Status,
		Config:       task.Config,
		Hierarchical: task.Hierarchical,
		StartTime:    task.StartTime,
		EndTime:      task.EndTime,
		Result:       task.Result,
		Error:        task.Error,
		Progress:     task.Progress,
	}, nil
}

// GetTaskResult waits for task completion and returns result
func (tm *CommunityTaskManager) GetTaskResult(taskID string, timeout time.Duration) ([]*types.Community, error) {
	deadline := time.Now().Add(timeout)

	for {
		task, err := tm.GetTaskStatus(taskID)
		if err != nil {
			return nil, err
		}

		switch task.Status {
		case TaskStatusComplete:
			return task.Result, nil
		case TaskStatusFailed:
			return nil, fmt.Errorf("%w: %v", ErrTaskFailed, task.Error)
		case TaskStatusPending, TaskStatusRunning:
			// Wait and retry
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("timeout waiting for task completion")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// CancelTask cancels a pending or running task
func (tm *CommunityTaskManager) CancelTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, ok := tm.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}

	// Can only cancel pending tasks
	if task.Status != TaskStatusPending {
		return errors.New("can only cancel pending tasks")
	}

	task.Status = TaskStatusFailed
	task.Error = errors.New("task cancelled")
	return nil
}

// CleanupOldTasks removes completed tasks older than the specified duration
func (tm *CommunityTaskManager) CleanupOldTasks(maxAge time.Duration) int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	removed := 0

	for id, task := range tm.tasks {
		if (task.Status == TaskStatusComplete || task.Status == TaskStatusFailed) && task.EndTime.Before(cutoff) {
			delete(tm.tasks, id)
			removed++
		}
	}

	return removed
}

// GetAllTasks returns all tasks for a session
func (tm *CommunityTaskManager) GetAllTasks(sessionID string) []*CommunityTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*CommunityTask, 0)
	for _, task := range tm.tasks {
		if task.SessionID == sessionID {
			tasks = append(tasks, &CommunityTask{
				ID:           task.ID,
				SessionID:    task.SessionID,
				Status:       task.Status,
				Config:       task.Config,
				Hierarchical: task.Hierarchical,
				StartTime:    task.StartTime,
				EndTime:      task.EndTime,
				Error:        task.Error,
				Progress:     task.Progress,
			})
		}
	}

	return tasks
}

// Shutdown gracefully shuts down the task manager
func (tm *CommunityTaskManager) Shutdown() {
	tm.cancel()
	close(tm.queue)
}

// GetStats returns task manager statistics
func (tm *CommunityTaskManager) GetStats() TaskManagerStats {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	stats := TaskManagerStats{
		TotalTasks: len(tm.tasks),
		QueueSize:  len(tm.queue),
		Workers:    tm.workers,
	}

	for _, task := range tm.tasks {
		switch task.Status {
		case TaskStatusPending:
			stats.PendingTasks++
		case TaskStatusRunning:
			stats.RunningTasks++
		case TaskStatusComplete:
			stats.CompletedTasks++
		case TaskStatusFailed:
			stats.FailedTasks++
		}
	}

	return stats
}

// TaskManagerStats holds task manager statistics
type TaskManagerStats struct {
	TotalTasks     int
	PendingTasks   int
	RunningTasks   int
	CompletedTasks int
	FailedTasks    int
	QueueSize      int
	Workers        int
}
