package shm

import (
	"github.com/kevinu2/shm/model"
	"os"
	"syscall"
	"unsafe"
)

func NewMMapMem(path string, size int) (*model.Mem, error) {
	if size > model.MaxCapacity {
		return nil, model.ErrOutOfCapacity
	}

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	// mmap不会更改底层文件的大小，我们要确保访问的映射地址不会超过文件大小，否则会panic
	// 这里设置一下底层文件大小
	if info.Size() != int64(size) {
		if err := fd.Truncate(int64(size)); err != nil {
			return nil, err
		}
	}

	// 使用syscall的mmap接口，创建内存映射
	// mmap接口相比posix接口，少了一个addr参数，如果有需要可以使用syscall.Syscall6接口
	// MAP_SHARED指定映射的类型，该模式下对映射空间的更新对其他进程的映射可见，并且会写回底层文件
	// 映射内存会通过[]byte的形式返回
	buf, err := syscall.Mmap(int(fd.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	// mmap返回之后，底层文件的设备描述符可以立即close掉
	fd.Close() // After the mmap() call has returned, the file descriptor can be closed immediately

	tagBuf := buf[:8]
	data := buf[8:]

	tag := (*model.Tag)(unsafe.Pointer(&tagBuf[0]))

	return &model.Mem{
		Tag:   tag,
		Queue: data,
	}, nil

	//// 这里，直接将映射的内存强制类型转换成一个int指针p
	//p := (*int)(unsafe.Pointer(&buf[0]))
	//// 对指针p的读取就是读取文件内容
	//log.Printf("the value saved on file is %d", *p)
	//// 对指针p的更新就是更新文件内容
	//*p = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000)
	//// 脏页写回并不是及时的，使用msync系统调用强制将更新写回磁盘文件
	//_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(p)), uintptr(size), syscall.MS_SYNC)
	//if errno != 0 {
	//	log.Fatal(syscall.Errno(errno))
	//}
	//// 使用munmap系统调用需求内存映射
	//_, _, errno = syscall.Syscall(syscall.SYS_MUNMAP, uintptr(unsafe.Pointer(p)), uintptr(size), 0)
	//if errno != 0 {
	//	log.Fatal(syscall.Errno(errno))
	//}
}
