package queue

import (
	"errors"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

const (
	maxCapacity = 1024 * 1024 * 1024
	// IpcCreate create if key is nonexistent
	IpcCreate = 00001000
)

var ErrOutOfCapacity = errors.New("out of capacity")

type tag struct {
	readIndex  int32
	writeIndex int32
}

type shmMem struct {
	*tag
	queue []byte
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

	tag := (*tag)(unsafe.Pointer(shmAddr))

	var data []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = shmAddr + 8
	sh.Len = size
	sh.Cap = size

	return &shmMem{
		tag:   tag,
		queue: data,
	}, nil
}
