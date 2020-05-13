package sltest

import (
	"fmt"
	"testing"
)

//var wg sync.WaitGroup

func task(i int) {
	fmt.Println("task...", i)
	//wg.Done()
}

func TestWG(t *testing.T) {
	for i := 0; i < 10; i++ {
		//wg.Add(1)
		go task(i)
	}
	//wg.Wait()
	fmt.Println("over")
}
