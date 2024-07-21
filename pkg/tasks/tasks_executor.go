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

func (te *TaskExecutor) Add(task func()) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.q = append(te.q, &Task{task: task})
	te.qNotEmpty.Signal()
}

func (te *TaskExecutor) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				task := te.pop()
				if task != nil {
					task.task()
				}

			}
		}()
	}
}

type Task struct {
	task func()
}

func newTask(task func()) *Task {
	return &Task{
		task: task,
	}
}
