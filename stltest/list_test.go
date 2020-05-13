package sltest

import (
	"container/list"
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	l := list.New()
	for i := 0; i < 10; i += 1 {
		l.PushBack(i)
	}

	for l.Len() > 0 {
		fmt.Printf("list size: %d\n", l.Len())
		v := l.Remove(l.Front())
		fmt.Printf("value: %d\n", v.(int))
		fmt.Printf("list size: %d\n", l.Len())
	}
	if e := l.Front(); e != nil {
		l.Remove(e)
	}
}
