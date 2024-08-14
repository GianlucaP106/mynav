package tasks

var executor *TaskExecutor = NewTaskExecutor()

func QueueTask(task func()) *Task {
	return executor.QueueTask(task)
}

func StartExecutor() {
	executor.Start(10)
}
