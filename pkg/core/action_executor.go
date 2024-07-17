package core

type ActionExecutor struct {
	q []*Action
}

func NewActionExecutor() *ActionExecutor {
	return &ActionExecutor{
		q: make([]*Action, 0),
	}
}

func (tq *ActionExecutor) Pop() *Action {
	if len(tq.q) == 0 {
		return nil
	}

	first := tq.q[0]
	if len(tq.q) == 1 {
		tq.q = []*Action{}
		return first
	}

	rest := tq.q[1:]
	tq.q = rest
	return first
}

func (tq *ActionExecutor) Add(task *Action) {
	tq.q = append(tq.q, task)
}

func (tq *ActionExecutor) RunAll() {
	for {
		task := tq.Pop()
		if task == nil {
			break
		}

		task.task()
	}
}

func (tq *ActionExecutor) RunOne() {
	task := tq.Pop()
	if task != nil {
		task.task()
	}
}

type Action struct {
	task func()
}

func NewAction(action func()) *Action {
	return &Action{
		task: action,
	}
}
