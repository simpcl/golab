package executor

import (
	"container/list"
	"hash/crc32"
	"log"
	"reflect"
)

const (
	MaxSelectCases         = 1024
	DispatcherCount        = 8
	WorkerCount            = 8
	MaxDispatcherTaskCount = 128
	MaxWorkerTaskCount     = 1
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
	consumer         *Consumer
	isTaskChSelected bool
	selectedCases    [MaxSelectCases]reflect.SelectCase
	maxWorkerCount   int
	freeWorkers      *list.List
	busyWorkers      map[string]*Worker
	pendingTasks     *list.List
}

func newDispatcher(workerCount int) *Dispatcher {
	pc := newProcesserCtx(MaxDispatcherTaskCount)

	fw := list.New()
	for i := 0; i < workerCount; i += 1 {
		fw.PushBack(newWorker(MaxWorkerTaskCount))
	}

	bw := make(map[string]*Worker)

	tasks := list.New()

	dispatcher := &Dispatcher{ProcesserCtx: *pc, maxWorkerCount: workerCount, freeWorkers: fw, busyWorkers: bw, pendingTasks: tasks}

	dispatcher.consumer = newConsumer(dispatcher)

	return dispatcher
}

func (d *Dispatcher) getResultChannel() interface{} {
	return nil
}

func (d *Dispatcher) submitToWorker(task *Task) bool {
	key := task.GetKey()
	if _, found := d.busyWorkers[key]; found {
		return false
	}

	if e := d.freeWorkers.Front(); e != nil {
		worker := d.freeWorkers.Remove(e).(*Worker)
		d.busyWorkers[key] = worker
		worker.submit(task)
		return true
	}

	return false
}

func (d *Dispatcher) drainPendingTasks() int {
	if d.freeWorkers.Len() == 0 {
		return d.pendingTasks.Len()
	}

	for e := d.pendingTasks.Front(); e != nil; {
		task := e.Value.(*Task)
		if d.submitToWorker(task) != true {
			return d.pendingTasks.Len()
		}
		d.pendingTasks.Remove(e)
	}

	return d.pendingTasks.Len()
}

func (d *Dispatcher) process() {
	for {
	Loop:
		d.drainPendingTasks()
		cases := d.updateSelectCases()
		chose, value, ok := reflect.Select(cases)
		if ok != true {
			continue
		}
		if d.isTaskChSelected == true && chose == 0 {
			d.addCount(1)
			task := value.Interface().(*Task)
			if d.submitToWorker(task) == false {
				d.pendingTasks.PushBack(task)
				goto Loop
			}
		} else {
			d.removeCount(1)
			key := value.Interface().(string)
			if worker, found := d.busyWorkers[key]; found {
				delete(d.busyWorkers, key)
				d.freeWorkers.PushBack(worker)
			}
		}
	}
}

func (d *Dispatcher) updateSelectCases() []reflect.SelectCase {
	var i = 0
	d.isTaskChSelected = false
	//log.Printf("updateSelectCases, freeWorkers: %d, pendingTasks: %d\n", d.freeWorkers.Len(), d.pendingTasks.Len())
	if /*d.freeWorkers.Len() > 0 &&*/ d.isFull() != true {
		d.selectedCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(d.taskCh),
		}
		i += 1
		d.isTaskChSelected = true
	}
	for _, consumer := range d.busyWorkers {
		d.selectedCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(consumer.keyCh),
		}
		i += 1
	}
	return d.selectedCases[:i]
}

type ConsistencyExecutor struct {
	ds    []*Dispatcher
	dsNum int
}

func NewConsistencyExecutor(dsNum int) *ConsistencyExecutor {
	if dsNum <= 0 {
		log.Fatal("Number of Dispatchers should be > 0")
	}
	ce := &ConsistencyExecutor{dsNum: dsNum}
	ce.ds = []*Dispatcher{}
	for i := 0; i < dsNum; i += 1 {
		ce.ds = append(ce.ds, newDispatcher(WorkerCount))
	}
	return ce
}

func (ce *ConsistencyExecutor) Submit(r Runner) *Future {
	task := NewTask(r)
	future := task.GetFuture()
	key := []byte(r.GetKey())
	index := int(crc32.ChecksumIEEE(key)) % len(ce.ds)
	ce.ds[index].submit(task)
	return future
}
