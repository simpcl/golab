package executor

import (
	"container/list"
	"hash/crc32"
	"reflect"
)

const (
	DispatcherCount        = 4
	WorkerCount            = 8
	MaxDispatcherTaskCount = 128
	MaxWorkerTaskCount     = 16
)

type Worker struct {
	ProcesserCtx
	consumer *Consumer
	keyCh    chan string
}

func newWorker(taskCount int) *Worker {
	pc := newProcesserCtx(taskCount)
	w := &Worker{ProcesserCtx: *pc, keyCh: make(chan string, taskCount)}
	w.consumer = newConsumer(w)
	return w
}

func (w *Worker) processTask(task *Task) {
	if task == nil {
		return
	}
	res := task.Run()
	task.resCh <- res
	w.keyCh <- task.GetKey()
}

func (w *Worker) process() {
	for {
		select {
		case task := <-w.taskCh:
			w.processTask(task)
		}
	}
}

func (w *Worker) getResultChannel() interface{} {
	return w.keyCh
}

type Dispatcher struct {
	ProcesserCtx
	consumer       *Consumer
	maxWorkerCount int
	freeWorkers    *list.List
	busyWorkers    map[string]*Worker
}

func newDispatcher(workerCount int) *Dispatcher {
	pc := newProcesserCtx(MaxDispatcherTaskCount)

	fw := list.New()
	for i := 0; i < workerCount; i += 1 {
		fw.PushBack(newWorker(MaxWorkerTaskCount))
	}

	bw := make(map[string]*Worker)

	dispatcher := &Dispatcher{ProcesserCtx: *pc, maxWorkerCount: workerCount, freeWorkers: fw, busyWorkers: bw}

	dispatcher.consumer = newConsumer(dispatcher)

	return dispatcher
}

func (d *Dispatcher) getResultChannel() interface{} {
	return nil
}

func (d *Dispatcher) process() {
	for {
	Loop:
		cases := d.updateSelectCases()
		chose, value, _ := reflect.Select(cases)
		switch chose {
		case 0:
			task := value.Interface().(*Task)
			key := task.GetKey()
			if worker, found := d.busyWorkers[key]; found {
				worker.submit(task)
				goto Loop
			}
			if e := d.freeWorkers.Front(); e != nil {
				worker := d.freeWorkers.Remove(e).(*Worker)
				d.busyWorkers[key] = worker
				worker.submit(task)
				goto Loop
			}
		default:
			key := value.Interface().(string)
			if worker, found := d.busyWorkers[key]; found {
				refCount := worker.finish()
				if refCount <= 0 {
					delete(d.busyWorkers, key)
					d.freeWorkers.PushBack(worker)
				}
			}
		}
	}
}

func (d *Dispatcher) updateSelectCases() (cases []reflect.SelectCase) {
	var sc reflect.SelectCase
	if d.freeWorkers.Len() > 0 {
		sc = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(d.taskCh),
		}
		cases = append(cases, sc)
	}
	for _, consumer := range d.busyWorkers {
		sc = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(consumer.keyCh),
		}
		cases = append(cases, sc)
	}
	return
}

type ConsistencyExecutor struct {
	ds [DispatcherCount]*Dispatcher
}

func NewConsistencyExecutor() *ConsistencyExecutor {
	ce := &ConsistencyExecutor{}
	for i := 0; i < len(ce.ds); i += 1 {
		ce.ds[i] = newDispatcher(WorkerCount)
	}
	return ce
}

func (ce *ConsistencyExecutor) Submit(r Runner) *Future {
	task := NewTask(r)
	future := task.GetFuture()
	key := []byte(r.GetKey())
	index := int(crc32.ChecksumIEEE(key)) % len(ce.ds)
	ce.ds[index].getTaskChannel() <- task
	return future
}
