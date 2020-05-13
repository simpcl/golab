package executor

import "sync"

type Processer interface {
	submit(task *Task) int
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

func (pc *ProcesserCtx) submit(task *Task) int {
	if pc.processingCount >= pc.maxCount {
		return -1
	}
	pc.taskCh <- task
	pc.processingCount += 1

	return pc.maxCount - pc.processingCount
}

func (pc *ProcesserCtx) finish() int {
	if pc.processingCount <= 0 {
		pc.processingCount = 0
		return 0
	}
	pc.processingCount -= 1
	return pc.processingCount
}

func (pc *ProcesserCtx) getTaskChannel() chan *Task {
	return pc.taskCh
}
