package sltest

import (
	"sync/atomic"
	"testing"
)

var n int32 = 0

func BenchmarkLoadInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atomic.LoadInt32(&n)
	}
}
