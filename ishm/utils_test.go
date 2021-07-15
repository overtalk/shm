package ishm

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack(20)
	for i := 0; i < 10; i++ {
		s.Push(i)
	}
	fmt.Printf("%d, %d\n", len(s.eles), cap(s.eles))
	for i := 0; i < 10; i++ {
		fmt.Println(i, ":", s.Pop().(int))
	}
}
