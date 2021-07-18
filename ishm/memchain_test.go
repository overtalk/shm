package ishm

import (
	"fmt"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	fmt.Println("Test MemChain Insert")
	MC := NewMemChain()
	b1 := NewBlock(0, 10, true)
	MC.Insert(b1)
	b2 := NewBlock(20, 45, true)
	err := MC.Insert(b2)
	if err != nil {
		fmt.Println(err.Error())
	}
	b3, err := MC.SearchBlock(13)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Use Block: ", b3.start, ", ", b3.end)
	MC.PrintChain()
}

func TestSearchMerge(t *testing.T) {
	fmt.Println("Test MemChain Merge")
	MC := NewMemChain()
	b1 := NewBlock(0, 100, true)
	MC.Insert(b1)
	b2, _ := MC.SearchBlock(10)
	b3, _ := MC.SearchBlock(10)
	b4, _ := MC.SearchBlock(10)
	b5, _ := MC.SearchBlock(10)
	fmt.Println("Use Block: ", b2.start, ", ", b2.end)
	fmt.Println("Use Block: ", b3.start, ", ", b3.end)
	fmt.Println("Use Block: ", b4.start, ", ", b4.end)
	fmt.Println("Use Block: ", b5.start, ", ", b5.end)
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	MC.Insert(b3)
	MC.Insert(b4)
	MC.PrintChain()
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	MC.MergeBlocks()
	MC.PrintChain()
}

func TestMergeOnTime(t *testing.T) {
	fmt.Println("Test MemChain MergeOnTime")
	MC := NewMemChain()
	b1 := NewBlock(0, 100, true)
	MC.Insert(b1)
	MC.SearchBlock(10)
	b3, _ := MC.SearchBlock(10)
	b4, _ := MC.SearchBlock(10)
	b5, _ := MC.SearchBlock(10)
	MC.Insert(b3)
	MC.Insert(b4)
	MC.PrintChain()
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	MC.MergeOnTime(1)
	time.Sleep(1 * time.Second)
	MC.PrintChain()
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	MC.Insert(b5)
	time.Sleep(1 * time.Second)
	MC.PrintChain()
}
