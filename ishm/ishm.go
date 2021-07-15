package ishm

// #include "ishm.h"
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

type SharedMemoryFlags int

const (
	IpcNone                        = 0
	IpcCreate    SharedMemoryFlags = C.IPC_CREAT
	IpcExclusive                   = C.IPC_EXCL
	HugePages                      = C.SHM_HUGETLB
	NoReserve                      = C.SHM_NORESERVE
)

// Segment is a native representation of a SysV shared memory segment
type Segment struct {
	Id     int64
	Size   int64
	offset int64
}

// Create a new shared memory segment with the given size (in bytes).  The system will automatically
// round the size up to the nearest memory page boundary (typically 4KB).
//
func Create(size int64) (*Segment, error) {
	return OpenSegment(size, (IpcCreate | IpcExclusive), 0600)
}

// Open an existing shared memory segment located at the given ID.  This ID is returned in the
// struct that is populated by Create(), or by the shmget() system call.
//
func Open(id int64) (*Segment, error) {
	sz, err := C.sysv_shm_get_size(C.int(id))
	if err == nil {
		return &Segment{
			Id:   id,
			Size: int64(sz),
		}, nil
	}
	return nil, err
}

// OpenSegment creates a shared memory segment of a given size, and also allows for the specification of
// creation flags supported by the shmget() call, as well as specifying permissions.
//
func OpenSegment(size int64, flags SharedMemoryFlags, perms os.FileMode) (*Segment, error) {
	var err error
	if shmid, err := C.sysv_shm_open(C.int(size), C.int(flags), C.int(perms)); err == nil {
		actualSize, err := C.sysv_shm_get_size(shmid)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve SHM size: %v", err)
		}
		return &Segment{
			Id:   int64(shmid),
			Size: int64(actualSize),
		}, nil
	}
	return nil, err
}

// DestroySegment destroy a shared memory segment by its ID
//
func DestroySegment(id int64) error {
	_, err := C.sysv_shm_close(C.int(id))
	return err
}

// ReadChunk reads some or all of the shared memory segment and return a byte slice.
//
func (s *Segment) ReadChunk(length int64, start int64) ([]byte, error) {
	if length < 0 {
		length = s.Size
	}

	buffer := C.malloc(C.size_t(length))
	defer C.free(buffer)

	if _, err := C.sysv_shm_read(C.int(s.Id), buffer, C.int(length), C.int(start)); err != nil {
		return nil, err
	}

	return C.GoBytes(buffer, C.int(length)), nil
}

func (s *Segment) Read(p []byte) (n int, err error) {
	if s.Id == 0 {
		return 0, fmt.Errorf("Cannot read shared memory segment: SHMID not set")
	}

	if s.offset >= s.Size {
		return 0, io.EOF
	}

	length := int64(len(p))

	fmt.Println("Buf size: ", length, ",offset ", s.offset)
	if length > s.Size {
		length = s.Size
	}

	if (length + s.offset) > s.Size {
		length = s.Size - s.offset
	}

	buffer := C.malloc(C.size_t(length))
	defer C.free(buffer)

	if _, err := C.sysv_shm_read(C.int(s.Id), buffer, C.int(length), C.int(s.offset)); err != nil {
		return 0, err
	}
	v := copy(p, C.GoBytes(buffer, C.int(length)))
	if v > 0 {
		s.offset += int64(v)
		return v, nil
	}
	return v, io.EOF
}

// Implements the io.Writer interface for shared memory
//
func (s *Segment) Write(p []byte) (n int, err error) {
	// if the offset runs past the segment size, we've reached the end
	if s.offset >= s.Size {
		return 0, io.EOF
	}

	length := int64(len(p))

	// write length cannot exceed segment size
	if length > s.Size {
		length = s.Size
	}

	// if length+offset would overrun, make length equal (size - offset), which is what remains
	if (length + s.offset) > s.Size {
		length = s.Size - s.offset
	}

	if _, err := C.sysv_shm_write(C.int(s.Id), unsafe.Pointer(&p[0]), C.int(length), C.int(s.offset)); err != nil {
		return 0, err
	} else {
		s.offset += length
		return int(length), nil
	}
}

// Resets the internal offset counter for this segment, allowing subsequent calls
// to Read() or Write() to start from the beginning.
//
func (s *Segment) Reset() {
	s.offset = 0
}

// Implements the io.Seeker interface for shared memory.  Subsequent calls to Read()
// or Write() will start from this position.
//
func (s *Segment) Seek(offset int64, whence int) (int64, error) {
	var computedOffset int64

	switch whence {
	case 1:
		computedOffset = s.offset + offset
	case 2:
		computedOffset = s.Size - offset
	default:
		computedOffset = offset
	}

	if computedOffset < 0 {
		return 0, fmt.Errorf("Cannot seek to position before start of segment")
	}

	s.offset = computedOffset
	return s.offset, nil
}

// Returns the current position of the Read/Write pointer.
//
func (s *Segment) Position() int64 {
	return s.offset
}

// Attaches the segment to the current processes resident memory.  The pointer
// that is returned is the actual memory address of the shared memory segment
// for use with third party libraries that can directly read from memory.
//
func (s *Segment) Attach() (unsafe.Pointer, error) {
	if addr, err := C.sysv_shm_attach(C.int(s.Id)); err == nil {
		return unsafe.Pointer(addr), nil
	} else {
		return nil, err
	}
}

// Detaches the segment from the current processes memory space.
//
func (s *Segment) Detach(addr unsafe.Pointer) error {
	_, err := C.sysv_shm_detach(addr)
	return err
}

// Destroys the current shared memory segment.
//
func (s *Segment) Destroy() error {
	return DestroySegment(s.Id)
}
