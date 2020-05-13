package stltest

import (
	"fmt"
	"testing"
)

type Animal interface {
	Spark()
}

type Mammal struct {
	heads int
}

func (m Mammal) Spark() {
	fmt.Println("branimal spark")
}

type Dog struct {
	Mammal
	legs int
}

func (d Dog) Spark() {
	fmt.Println("dog spark")
}

func testAnimal(a Animal) {
	a.Spark()
}

func TestOverride(t *testing.T) {
	m := &Mammal{}
	testAnimal(m)
	d := &Dog{}
	testAnimal(d)
}
