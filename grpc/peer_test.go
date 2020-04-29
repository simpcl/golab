package grpc

import (
	"testing"
	"time"
)

func TestRpcCall(t *testing.T) {
	master := &Master{}
	go master.start("10008")
	time.Sleep(1 * time.Second)

	worker := &Worker{}
	worker.start("127.0.0.1:10008")
	defer worker.destroy()
}
