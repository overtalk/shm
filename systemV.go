package shm

import (
	"fmt"
	"github.com/kevinu2/shm/model"
	"reflect"
	"syscall"
	"unsafe"
)

const ipcCreate = 00001000

func NewSystemVMem(key, size int) (*model.Mem, error) {
	if size > model.MaxCapacity {
		return nil, model.ErrOutOfCapacity
	}

	shmID, _, errCode := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), ipcCreate|0600)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}

	shmAddr, _, errCode := syscall.Syscall(syscall.SYS_SHMAT, shmID, 0, 0)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}

	var data []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = shmAddr + 8
	sh.Len = size
	sh.Cap = size

	return &model.Mem{
		Tag:   (*model.Tag)(unsafe.Pointer(shmAddr)),
		Queue: data,
	}, nil
}
func GetSHMInfo(key, size int) (*reflect.SliceHeader, error) {
	if size > model.MaxCapacity {
		return nil, model.ErrOutOfCapacity
	}

	shmID, _, errCode := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), ipcCreate|0600)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}

	shmAddr, _, errCode := syscall.Syscall(syscall.SYS_SHMAT, shmID, 0, 0)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %d\n", errCode)
	}


	var data []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = shmAddr + 8
	sh.Len = size
	sh.Cap = size

	return sh, errCode
}
