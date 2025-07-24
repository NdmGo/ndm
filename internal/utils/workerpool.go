package utils

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Task task interface
type Task interface {
	Execute(ctx context.Context) error
}

// TaskFunc function type task
type TaskFunc func(ctx context.Context) error

// Execute implement Task interface
func (f TaskFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

// WorkerPool worker pool
type WorkerPool struct {
	workerCount int
	taskQueue   chan Task
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	started     bool
	mu          sync.RWMutex
}

// NewWorkerPool create new worker pool
func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU() * 2 // default to 2 times the number of CPU cores
	}
	if queueSize <= 0 {
		queueSize = 1000 // default queue size
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan Task, queueSize),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start start worker pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return
	}

	wp.started = true
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop stop worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.started {
		return
	}

	wp.cancel()
	close(wp.taskQueue)
	wp.wg.Wait()
	wp.started = false
}

// Submit submit task
func (wp *WorkerPool) Submit(task Task) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.started {
		return ErrWorkerPoolNotStarted
	}

	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return ErrTaskQueueFull
	}
}

// SubmitFunc submit function task
func (wp *WorkerPool) SubmitFunc(fn func(ctx context.Context) error) error {
	return wp.Submit(TaskFunc(fn))
}

// SubmitWithTimeout submit task with timeout
func (wp *WorkerPool) SubmitWithTimeout(task Task, timeout time.Duration) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.started {
		return ErrWorkerPoolNotStarted
	}

	ctx, cancel := context.WithTimeout(wp.ctx, timeout)
	defer cancel()

	select {
	case wp.taskQueue <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// worker worker goroutine
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}

			// execute task, catch panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						// log panic
						// log.Errorf("Worker panic: %v", r)
					}
				}()

				if err := task.Execute(wp.ctx); err != nil {
					// log task execution error
					// log.Errorf("Task execution error: %v", err)
				}
			}()

		case <-wp.ctx.Done():
			return
		}
	}
}

// GetStats get worker pool statistics
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	return WorkerPoolStats{
		WorkerCount:  wp.workerCount,
		QueueSize:    cap(wp.taskQueue),
		PendingTasks: len(wp.taskQueue),
		IsRunning:    wp.started,
	}
}

// WorkerPoolStats worker pool statistics
type WorkerPoolStats struct {
	WorkerCount  int  `json:"worker_count"`
	QueueSize    int  `json:"queue_size"`
	PendingTasks int  `json:"pending_tasks"`
	IsRunning    bool `json:"is_running"`
}

// global worker pool instance
var (
	DefaultWorkerPool *WorkerPool
	once              sync.Once
)

// GetDefaultWorkerPool get default worker pool
func GetDefaultWorkerPool() *WorkerPool {
	once.Do(func() {
		DefaultWorkerPool = NewWorkerPool(0, 0)
		DefaultWorkerPool.Start()
	})
	return DefaultWorkerPool
}

// error definitions
var (
	ErrWorkerPoolNotStarted = NewErr(nil, "worker pool not started")
	ErrTaskQueueFull        = NewErr(nil, "task queue is full")
)
