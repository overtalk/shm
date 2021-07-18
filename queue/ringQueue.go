package queue

import (
	"errors"

	"github.com/kevinu2/shm/model"
)

// 一个生产者 & 一个消费者的情况之下
// 生产者只修改 writeIndex，消费者只修改 readIndex
// 由此不需要使用 atomic 包操作
type RingQueue struct {
	// 环形队列长度
	// 可以写入的长度为 queueLen-1，这样就可以区分 IsEmpty/IsFull
	// IsEmpty ： readIndex = writeIndex
	// IsFull：readIndex + writeIndex = queueLen
	queueLen int32
	shm      *model.Mem
}

func NewRingQueue(shm *model.Mem, size int) (*RingQueue, error) {
	if len(shm.Queue) != size {
		return nil, errors.New("unmatched size and shm")
	}
	return &RingQueue{
		queueLen: int32(size),
		shm:      shm,
	}, nil
}

// 队列已经使用的空间长度
func (rq *RingQueue) getUsedLen() int32 {
	if rq.shm.WriteIndex >= rq.shm.ReadIndex {
		return rq.shm.WriteIndex - rq.shm.ReadIndex
	}
	return rq.queueLen - rq.shm.ReadIndex + rq.shm.WriteIndex
}

func (rq *RingQueue) getLeftLen() int32 {
	return rq.queueLen - 1 - rq.getUsedLen()
}

// 判断队列是否为空
func (rq *RingQueue) isEmpty() bool {
	return rq.shm.WriteIndex == rq.shm.ReadIndex
}

// 判断队列是否满了
func (rq *RingQueue) isFull() bool {
	return rq.shm.WriteIndex+rq.shm.ReadIndex < rq.queueLen
}

func (rq *RingQueue) Save(data []byte) error {
	buf := newBinaryMessage(data).serialize()

	if len(buf) > int(rq.getLeftLen()) {
		return model.ErrOutOfCapacity
	}

	for index, bit := range buf {
		rq.shm.Queue[int(rq.shm.WriteIndex)+index] = bit
	}

	rq.shm.WriteIndex += int32(len(buf))
	rq.shm.ReadIndex %= rq.queueLen

	return nil
}

func (rq *RingQueue) Get() ([][]byte, error) {
	if rq.getUsedLen() == 0 {
		return nil, nil
	}

	currentWriteIndex := rq.shm.WriteIndex

	var retBytes []byte
	if rq.shm.ReadIndex < currentWriteIndex {
		retBytes = rq.shm.Queue[rq.shm.ReadIndex:currentWriteIndex]
	} else {
		retBytes = append(rq.shm.Queue[rq.shm.ReadIndex:], rq.shm.Queue[:currentWriteIndex]...)
	}

	rq.shm.ReadIndex = currentWriteIndex

	binaryMessages, err := deserializeSlice(retBytes)
	if err != nil {
		return nil, err
	}

	var ret [][]byte
	for _, bm := range binaryMessages {
		ret = append(ret, bm.Body)
	}

	return ret, nil
}
