package tasks

import (
	"sync"
)

type TaskExecutor struct {
	mu        *sync.Mutex
	qNotEmpty *sync.Cond
	q         []*Task
}

func NewTaskExecutor() *TaskExecutor {
	te := &TaskExecutor{}
	te.q = make([]*Task, 0)
	te.mu = &sync.Mutex{}
	te.qNotEmpty = sync.NewCond(te.mu)
	return te
}

func (te *TaskExecutor) pop() *Task {
	te.mu.Lock()
	defer te.mu.Unlock()

	for len(te.q) == 0 {
		te.qNotEmpty.Wait()
	}

	first := te.q[0]
	if len(te.q) == 1 {
		te.q = []*Task{}
		return first
	}

	rest := te.q[1:]
	te.q = rest
	return first
}

func (te *TaskExecutor) QueueTask(task func()) *Task {
	te.mu.Lock()
	defer te.mu.Unlock()
	t := newTask(task)
	te.q = append(te.q, t)
	te.qNotEmpty.Signal()
	return t
}

func (te *TaskExecutor) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				task := te.pop()
				if task != nil {
					task.start()
					task.task()
					task.complete()
				}

			}
		}()
	}
}

type Task struct {
	task      func()
	mu        sync.RWMutex
	started   bool
	completed bool
}

func newTask(task func()) *Task {
	return &Task{
		task:      task,
		started:   false,
		completed: false,
	}
}

func (t *Task) IsCompleted() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.completed
}

func (t *Task) IsStarted() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.started
}

func (t *Task) complete() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.completed = true
}

func (t *Task) start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.started = true
}
