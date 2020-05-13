package sltest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("watcher %s quit\n", name)
			return
		default:
			fmt.Printf("watcher %s is working\n", name)
			time.Sleep(time.Second * 2)
		}
	}
}

//func TestContext1(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	go watch(ctx, "ctx1")
//	go watch(ctx, "ctx2")
//
//	time.Sleep(time.Second * 10)
//	fmt.Println("cancel the watcher")
//	cancel()
//	time.Sleep(time.Second * 5)
//}

func TestContext2(t *testing.T) {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			fmt.Printf("key %s found, value: %s\n", k, v)
			return
		}
		fmt.Printf("key %s not found\n", k)
	}

	ctx := context.WithValue(context.Background(), "k", "v")
	f(ctx, "k")
}

type Task struct {
	n int
}

func TestContext3(t *testing.T) {
	f := func(ctx context.Context, k string) {
		if v := ctx.Value(k); v != nil {
			t, b := v.(*Task)
			if b {
				fmt.Printf("key %s found, value: task %d\n", k, t.n)
			} else {
				fmt.Printf("key %s found, value type dismatch: %v\n", k, v)
			}
			return
		}
		fmt.Printf("key %s not found\n", k)
	}

	ctx := context.WithValue(context.Background(), "k", &Task{n: 3})
	f(ctx, "k")
}
