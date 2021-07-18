package queue

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/kevinu2/shm/model"
)

type MultiQueue struct {
	// 环形队列长度
	// 可以写入的长度为 queueLen-1，这样就可以区分 IsEmpty/IsFull
	// IsEmpty ： readIndex = writeIndex
	// IsFull：readIndex + writeIndex = queueLen
	queueLen int32
	shm      *model.Mem
}

func NewMultiQueue(shm *model.Mem, size int) (*MultiQueue, error) {
	if len(shm.Queue) != size {
		return nil, errors.New("unmatched size and shm")
	}

	return &MultiQueue{
		queueLen: int32(size),
		shm:      shm,
	}, nil
}

func (multiArray *MultiQueue) Save(buf []byte) error {
	// cal the displacement
	blockCount := len(buf) / int(blockDataSize)
	lastBlockSize := len(buf) % int(blockDataSize)
	if lastBlockSize != 0 {
		blockCount++
	}
	displacement := int32(blockCount * int(blockSize))

	// modify the writeIndex
	_, currentWriteIndex, usedLen := multiArray.getUsedLen()
	if usedLen+displacement > multiArray.queueLen-1 {
		return model.ErrOutOfCapacity
	}

	if !atomic.CompareAndSwapInt32(&multiArray.shm.WriteIndex, currentWriteIndex, (currentWriteIndex+displacement)%multiArray.queueLen) {
		return errors.New("failed to swap")
	}

	// save data
	for i := blockCount; i > 0; i-- {
		var dataToSave []byte
		if i == blockCount && lastBlockSize != 0 {
			dataToSave = append(dataToSave, buf[(blockCount-1)*int(blockDataSize):]...)
		} else {
			start := (i - 1) * int(blockDataSize)
			dataToSave = append(dataToSave, buf[start:start+int(blockDataSize)]...)
		}

		b, _ := newBlock(uint16(blockCount), uint16(i-1), dataToSave)

		index := (int(currentWriteIndex) + (i-1)*int(blockSize)) % int(multiArray.queueLen)
		multiArray.writeBlock(index, b)
	}

	return nil
}

func (multiArray *MultiQueue) Get() ([][]byte, error) {
	// get block count
	currentReadIndex, _, usedLen := multiArray.getUsedLen()
	if int(usedLen)%int(blockSize) != 0 {
		return nil, errors.New("get error : usedLen is not n * blockSize")
	}
	blockCount := int(usedLen) / int(blockSize)
	// parse all blocks
	blocks := multiArray.getBlocks(currentReadIndex, blockCount)
	return multiArray.parseBlocks(blocks)
}

// getUsedLen return currentReadIndex, currentWriteIndex, usedLen
func (multiArray *MultiQueue) getUsedLen() (int32, int32, int32) {
	currentWriteIndex := atomic.LoadInt32(&multiArray.shm.WriteIndex)
	currentReadIndex := atomic.LoadInt32(&multiArray.shm.ReadIndex)
	var usedLen int32
	if currentWriteIndex >= currentReadIndex {
		usedLen = currentWriteIndex - currentReadIndex
	} else {
		usedLen = multiArray.queueLen - currentReadIndex + currentWriteIndex
	}

	return currentReadIndex, currentWriteIndex, usedLen
}

// writeBlock is to write block details to queue
func (multiArray *MultiQueue) writeBlock(startIndex int, b *block) {
	for index, bit := range b.serialize() {
		multiArray.shm.Queue[startIndex+index] = bit
	}

	multiArray.shm.Queue[startIndex] = 1
}

// writeBlock is to get block details from queue
func (multiArray *MultiQueue) getBlock(startIndex int) *block {
	var data []byte
	endIndex := startIndex + int(blockSize)
	if endIndex < int(multiArray.queueLen) {
		data = multiArray.shm.Queue[startIndex:endIndex]
	} else {
		data = append(data, multiArray.shm.Queue[startIndex:]...)
		data = append(data, multiArray.shm.Queue[:endIndex-int(multiArray.queueLen)]...)
	}

	b, _ := newBlockFromBytes(data)
	return b
}

// getBlocks is to get a bunch of blocks from a startIndex
func (multiArray *MultiQueue) getBlocks(startIndex int32, blockCount int) []*block {
	// parse all blocks
	var blocks []*block
	for i := 1; i <= blockCount; i++ {
		blocks = append(blocks, multiArray.getBlock(int(startIndex)))
		startIndex += int32(blockSize)
	}
	return blocks
}

// parseBlocks is to parse blocks
func (multiArray *MultiQueue) parseBlocks(blocks []*block) ([][]byte, error) {
	var ret [][]byte

	for i := 0; i < len(blocks); {
		if !blocks[i].isCompleted(0) {
			time.Sleep(time.Millisecond * 5)
			if !blocks[i].isCompleted(0) {
				i++
				continue
			}
		}

		// find the first block
		start := i
		i += int(blocks[i].blockCount)

		// get data details
		var data []byte
		for _, b := range blocks[start:i] {
			data = append(data, b.data[:b.blockLen]...)
		}

		ret = append(ret, data)
	}

	return ret, nil
}
