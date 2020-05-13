package executor

type Runner interface {
	Run() interface{}
	GetKey() string
}

type Task struct {
	Runner
	resCh  chan interface{}
	future *Future
}

func NewTask(r Runner) *Task {
	return &Task{Runner: r, resCh: make(chan interface{}, 1), future: nil}
}

func (t *Task) GetFuture() *Future {
	if t.future != nil {
		return t.future
	}

	return newFuture(t.resCh)
}
