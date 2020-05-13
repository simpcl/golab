package sltest

import (
	"fmt"
	"testing"
	"time"
)

func watch(ch chan struct{}, name string) {
	select {
	case <-ch:
		fmt.Printf("watcher %s is notified\n", name)
		return
	}
}

func TestCh1(t *testing.T) {
	ch := make(chan struct{})
	go watch(ch, "a")
	go watch(ch, "b")
	time.Sleep(1 * time.Second)
	ch <- struct{}{}
	ch <- struct{}{}
	time.Sleep(2 * time.Second)
}
