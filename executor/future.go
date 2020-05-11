package executor

type Future struct {
	resCh chan interface{}
}

func newFuture(c chan interface{}) *Future {
	return &Future{resCh: c}
}

func (f *Future) GetResult() interface{} {
	return <-f.resCh
}
