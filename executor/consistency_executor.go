package executor

import "reflect"

type ConsistencyExecutor struct {
	taskCh chan *HashTask
}

func NewConsistencyExecutor() *ConsistencyExecutor {
	ce := &ConsistencyExecutor{taskCh: make(chan *HashTask, 1)}
	go ce.run()
	return ce
}

func (ce *ConsistencyExecutor) run() {
	select {
	case t := <-ce.taskCh:
		go func() {
			res := t.hashRunner.run()
			t.resCh <- res
		}()
	}
}

func (ce *ConsistencyExecutor) Submit(hr HashRunner) *Future {
	t := NewHashTask(hr)
	future := newFuture(t.resCh)
	ce.taskCh <- t
	return future
}

func (ce *ConsistencyExecutor) getSelectCases() (cases []reflect.SelectCase) {
	return nil
}
