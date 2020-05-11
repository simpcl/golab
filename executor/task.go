package executor

type Runner interface {
	run() interface{}
}

type HashRunner interface {
	Runner
	getHashCode() int32
}

type Task struct {
	resCh  chan interface{}
	runner Runner
}

type HashTask struct {
	resCh      chan interface{}
	hashRunner HashRunner
}

func NewTask(r Runner) *Task {
	return &Task{resCh: make(chan interface{}, 1), runner: r}
}

func NewHashTask(hr HashRunner) *HashTask {
	return &HashTask{resCh: make(chan interface{}, 1), hashRunner: hr}
}
