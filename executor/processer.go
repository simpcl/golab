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

func (c *ProcesserCtx) submit(task *Task) int {
	if c.processingCount >= c.maxCount {
		return -1
	}
	c.taskCh <- task
	c.processingCount += 1

	return c.maxCount - c.processingCount
}

func (c *ProcesserCtx) finish() int {
	if c.processingCount <= 0 {
		c.processingCount = 0
		return 0
	}
	c.processingCount -= 1
	return c.processingCount
}

func (c *ProcesserCtx) getTaskChannel() chan *Task {
	return c.taskCh
}
