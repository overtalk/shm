package ishm

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// MemBlock ...
type MemBlock struct {
	start, end int64
	prev, next *MemBlock
	once       bool
}

//Size ...
func (mb *MemBlock) Size() int64 {
	return mb.end - mb.start
}

//Equal ...
func (mb *MemBlock) Equal(otherMb *MemBlock) bool {
	return mb.start == otherMb.start && mb.end == otherMb.end
}

//After ...
func (mb *MemBlock) After(otherMb *MemBlock) bool {
	return mb.start >= otherMb.end
}

//NewBlock ...
func NewBlock(start, end int64, once ...bool) *MemBlock {
	b := false
	if len(once) > 0 {
		b = once[0]
	}
	return &MemBlock{
		start: start,
		end:   end,
		prev:  nil,
		next:  nil,
		once:  b,
	}
}

//MemChain ...
type MemChain struct {
	first, last *MemBlock
	length      int64
}

//Insert ...
func (mC *MemChain) Insert(mB *MemBlock) error {
	pointer := mC.first.next
	for !pointer.After(mB) {
		pointer = pointer.next
		if pointer == nil {
			break
		}
	}
	if pointer != nil && pointer.After(mB) {
		pointer.prev.next = mB
		mB.next = pointer
		mB.prev = pointer.prev
		pointer.prev = mB
		return nil
	}
	return errors.New("Can't insert")
}

//Delete ...
func (mC *MemChain) Delete(mB *MemBlock) error {
	pointer := mC.first.next
	for pointer.next != nil {
		if pointer.Equal(mB) {
			mB.prev.next = mB.next
			mB.next.prev = mB.prev
			mB = nil
			return nil
		}
		pointer = pointer.next
	}
	return errors.New("Can't find suck block")
}

//MergeBlocks ...
func (mC *MemChain) MergeBlocks() {
	pointer := mC.first.next
	for pointer.next.next != nil {
		if pointer.end == pointer.next.start {
			mergeBlock := NewBlock(pointer.start, pointer.next.end)
			pointer.prev.next = mergeBlock
			mergeBlock.prev = pointer.prev
			pointer.next.next.prev = mergeBlock
			mergeBlock.next = pointer.next.next
			pointer = mergeBlock
		} else {
			pointer = pointer.next
		}
	}
	pointer = nil
}

// SearchBlock ...
func (mC *MemChain) SearchBlock(needSize int64) (*MemBlock, error) {
	pointer := mC.first.next
	for pointer.next != nil {
		if pointer.Size() >= needSize {
			findBlock := NewBlock(pointer.start, pointer.start+needSize)
			if pointer.start+needSize == pointer.end {
				mC.Delete(pointer)
			} else {
				pointer.start += needSize
			}
			pointer = nil
			return findBlock, nil
		}
		pointer = pointer.next
	}
	return nil, errors.New("Need Size too large")
}

//MergeOnTime ...
func (mC *MemChain) MergeOnTime(gap int) {
	go func() {
		t := time.NewTicker(time.Duration(gap) * time.Second)
		for v := range t.C {
			fmt.Println("Merge Blocks at: ", v)
			mC.MergeBlocks()
		}
	}()
}

//PrintChain ...
func (mC *MemChain) PrintChain() {
	pointer := mC.first.next
	for pointer.next != nil {
		fmt.Print(pointer.start, ", ", pointer.end, " -> ")
		pointer = pointer.next
	}
	fmt.Println()
}

// NewMemChain ...
func NewMemChain() *MemChain {
	first := NewBlock(-2, -2)
	last := NewBlock(math.MaxInt64, math.MaxInt64)
	first.next = last
	last.prev = first
	return &MemChain{first, last, 0}
}
