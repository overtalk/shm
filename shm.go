package shm

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	maxCapacity = 1024 * 1024 * 1024
	// IpcCreate create if key is nonexistent
	IpcCreate = 00001000
)

var errOutOfCapacity = errors.New("out of capacity")

type shm struct {
	shm  *shmMem
	size int
}

type shmMem struct {
	head int // head
	tail int // tail
	data [maxCapacity]byte
}

func newShm(key, size int) (*shm, error) {
	if size > maxCapacity {
		return nil, errOutOfCapacity
	}

	shmID, _, errCode := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), IpcCreate|0600)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}

	shmAddr, _, errCode := syscall.Syscall(syscall.SYS_SHMAT, shmID, 0, 0)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}

	s := (*shmMem)(unsafe.Pointer(shmAddr))

	return &shm{
		size: size,
		shm:  s,
	}, nil
}

func (s *shm) save(buf []byte) error {
	currentHead := s.shm.head
	var capacity int
	length := len(buf)

	if currentHead > s.shm.tail {
		capacity = currentHead - s.shm.tail
	} else {
		capacity = s.size - (s.shm.tail - currentHead)
	}

	// above capacity
	if capacity < length {
		return errOutOfCapacity
	}

	for index, bit := range buf {
		s.shm.data[(s.shm.tail+index)%s.size] = bit
	}

	s.shm.tail = (s.shm.tail + length) % s.size
	return nil
}

func (s *shm) get() []byte {
	currentTail := s.shm.tail
	var ret []byte

	switch {
	case s.shm.head == currentTail:
		return nil
	case s.shm.head < currentTail:
		ret = append(ret, s.shm.data[s.shm.head:currentTail]...)
	default:
		ret = append(ret, s.shm.data[s.shm.head:s.size]...)
		ret = append(ret, s.shm.data[0:currentTail]...)
	}

	s.shm.head = currentTail
	return ret
}
