package model

import "errors"

const MaxCapacity = 8 << 30 // 8G

var ErrOutOfCapacity = errors.New("out of capacity")

type Tag struct {
	ReadIndex  int32
	WriteIndex int32
}

type Mem struct {
	*Tag
	Queue []byte
}
