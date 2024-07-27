package tasks

var executor *TaskExecutor = NewTaskExecutor()

func AddTask(task func()) {
	executor.Add(task)
}

func StartExecutor() {
	executor.Start(10)
}
