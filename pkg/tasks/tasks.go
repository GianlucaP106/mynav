package tasks

var executor *TaskExecutor = NewTaskExecutor()

func QueueTask(task func()) {
	executor.QueueTask(task)
}

func StartExecutor() {
	executor.Start(10)
}
