package queue

// 一个生产者 & 一个消费者的情况之下
// 生产者只修改 writeIndex，消费者只修改 readIndex
// 由此不需要使用 atomic 包操作
type RingQueue struct {
	// 环形队列长度
	// 可以写入的长度为 queueLen-1，这样就可以区分 IsEmpty/IsFull
	// IsEmpty ： readIndex = writeIndex
	// IsFull：readIndex + writeIndex = queueLen
	queueLen int32
	shm      *shmMem
}

func NewRingQueue(key, size int) (*RingQueue, error) {
	s, err := newShmMem(key, size)
	if err != nil {
		return nil, err
	}

	return &RingQueue{
		queueLen: int32(size),
		shm:      s,
	}, nil
}

// 队列已经使用的空间长度
func (rq *RingQueue) getUsedLen() int32 {
	if rq.shm.writeIndex >= rq.shm.readIndex {
		return rq.shm.writeIndex - rq.shm.readIndex
	}
	return rq.queueLen - rq.shm.readIndex + rq.shm.writeIndex
}

func (rq *RingQueue) getLeftLen() int32 {
	return rq.queueLen - 1 - rq.getUsedLen()
}

// 判断队列是否为空
func (rq *RingQueue) isEmpty() bool {
	return rq.shm.writeIndex == rq.shm.readIndex
}

// 判断队列是否满了
func (rq *RingQueue) isFull() bool {
	return rq.shm.writeIndex+rq.shm.readIndex < rq.queueLen
}

func (rq *RingQueue) Save(data []byte) error {
	buf := newBinaryMessage(data).serialize()

	if len(buf) > int(rq.getLeftLen()) {
		return ErrOutOfCapacity
	}

	for index, bit := range buf {
		rq.shm.queue[int(rq.shm.writeIndex)+index] = bit
	}

	rq.shm.writeIndex += int32(len(buf))
	rq.shm.readIndex %= rq.queueLen

	return nil
}

func (rq *RingQueue) Get() ([][]byte, error) {
	if rq.getUsedLen() == 0 {
		return nil, nil
	}

	currentWriteIndex := rq.shm.writeIndex

	var retBytes []byte
	if rq.shm.readIndex < currentWriteIndex {
		retBytes = rq.shm.queue[rq.shm.readIndex:currentWriteIndex]
	} else {
		retBytes = append(rq.shm.queue[rq.shm.readIndex:], rq.shm.queue[:currentWriteIndex]...)
	}

	rq.shm.readIndex = currentWriteIndex

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
