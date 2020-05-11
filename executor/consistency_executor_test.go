package executor

import (
	"fmt"
	"testing"
	"time"
)

type KeyTask struct {
}

func (kt *KeyTask) run() interface{} {
	time.Sleep(time.Second * 2)
	return 1
}

func (kt *KeyTask) getHashCode() int32 {
	return 1
}

func TestConsistencyExecutor(t *testing.T) {
	kt := &KeyTask{}
	ce := NewConsistencyExecutor()
	future := ce.Submit(kt)
	res := future.GetResult()
	switch res.(type) {
	case int:
		fmt.Printf("value %d\n", res.(int))
	}

}
