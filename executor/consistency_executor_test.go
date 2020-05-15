package executor

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type KeyTaskRunner struct {
	value int
	name  string
}

func newKeyTaskRunner(v int) *KeyTaskRunner {
	return &KeyTaskRunner{name: fmt.Sprintf("key%d", v%6), value: v}
}

func (ktr *KeyTaskRunner) Run() interface{} {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := int(r.Intn(2))
	time.Sleep(time.Duration(n) * time.Millisecond)
	//fmt.Printf("running %s %d\n", ktr.name, ktr.value)
	return ktr.value
}

func (ktr *KeyTaskRunner) GetKey() string {
	return ktr.name
}

func TestConsistencyExecutor(t *testing.T) {
	ce := NewConsistencyExecutor()
	fmt.Println(ce)
	fs := []*Future{}

	for i := 0; i < 80; i += 1 {
		ktr := newKeyTaskRunner(i)
		future := ce.Submit(ktr)
		fs = append(fs, future)
	}

	for _, future := range fs {
		res := future.GetResult()
		switch res.(type) {
		case int:
			fmt.Printf("value %d\n", res.(int))
		}
	}
}

func TestConsistencyExecutor2(t *testing.T) {
	ce := NewConsistencyExecutor()
	fmt.Println(ce)

	var wg sync.WaitGroup
	for i := 0; i < 100; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ktr := newKeyTaskRunner(i)
			future := ce.Submit(ktr)
			future.GetResult()
		}()
	}
	wg.Wait()
}
