package sltest

import (
	"testing"
)

func TestSlice1(t *testing.T) {
	ss := []int{1, 2, 4}
	ss = append(ss, 3)
	t.Logf("%v", ss)
}

func TestSlice2(t *testing.T) {
	ss := make([]int, 1, 2)
	t.Logf("len: %d, cap: %d\n", len(ss), cap(ss))
	ss[0] = 8
	for count := 3; count > 0; count = count - 1 {
		ss = append(ss, count)
	}
	t.Logf("%v", ss)
}
