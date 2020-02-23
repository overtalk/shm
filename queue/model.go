package queue

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

var ErrOutOfCapacity = errors.New("out of capacity")

type shmMem struct {
	readIndex  int32
	writeIndex int32
	queue      [maxCapacity]byte
}

func newShmMem(key, size int) (*shmMem, error) {
	if size > maxCapacity {
		return nil, ErrOutOfCapacity
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
	return s, nil
}
