package executor

import "sync"

type Processer interface {
	submit(task *Task)
	process()
	getTaskChannel() chan *Task
	getResultChannel() interface{}
}

type Consumer struct {
	wg sync.WaitGroup
}

func newConsumer(p Processer) *Consumer {
	c := &Consumer{}
	c.wg.Add(1)
	go func() {
		c.wg.Done()
		p.process()
	}()
	c.wg.Wait()
	return c
}

type ProcesserCtx struct {
	maxCount        int
	processingCount int
	taskCh          chan *Task
}

func newProcesserCtx(maxCount int) *ProcesserCtx {
	return &ProcesserCtx{maxCount: maxCount, processingCount: 0, taskCh: make(chan *Task, maxCount)}
}

func (pc *ProcesserCtx) submit(task *Task) {
	pc.taskCh <- task
}

func (pc *ProcesserCtx) addCount(n int) int {
	pc.processingCount += n
	return pc.processingCount
}

func (pc *ProcesserCtx) removeCount(n int) int {
	pc.processingCount -= n
	return pc.processingCount
}

func (pc *ProcesserCtx) isFull() bool {
	return pc.processingCount >= pc.maxCount
}

func (pc *ProcesserCtx) getTaskChannel() chan *Task {
	return pc.taskCh
}
