package ishm

import (
	"errors"
	"hash/adler32"
)

func calcSum(data []byte) uint32 {
	return adler32.Checksum(data)
}

func checkSum(data []byte, shouldBe uint32) bool {
	return adler32.Checksum(data) == shouldBe
}

// Stack ...
type Stack struct {
	eles   []interface{}
	l      int
	maxLen int
}

// Push ...
func (s *Stack) Push(v interface{}) error {
	if s.l == s.maxLen {
		return errors.New("Stack is full")
	}
	s.eles = append(s.eles, v)
	s.l++
	return nil
}

// Pop ...
func (s *Stack) Pop() interface{} {
	v := s.eles[s.l-1]
	s.eles = s.eles[:s.l-1]
	s.l--
	return v
}

// Len ...
func (s Stack) Len() int {
	return s.l
}

// IsEmpty ...
func (s Stack) IsEmpty() bool {
	return s.l == 0
}

//Clear ...
func (s *Stack) Clear() {
	s.eles = make([]interface{}, 0, s.maxLen)
}

// NewStack ...
func NewStack(maxSize int) *Stack {
	return &Stack{
		eles:   []interface{}{},
		maxLen: maxSize,
	}
}
