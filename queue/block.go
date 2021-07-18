package queue

import (
	"encoding/binary"
	"errors"
)

const (
	blockHeadSize uint16 = 7
)

var (
	maxBlockCount       uint16 = 100
	blockDataSize       uint16 = 16
	blockSize                  = blockDataSize + blockHeadSize
	ErrInvalidBlockSize        = errors.New("invalid block size")
	ErrInvalidDataSize         = errors.New("invalid block data size : too long")
)

type block struct {
	completed  uint16
	blockCount uint16
	blockIndex uint16
	blockLen   uint16
	data       []byte
}

func newBlock(blockCount, blockIndex uint16, data []byte) (*block, error) {
	if len(data) > int(blockSize) {
		return nil, ErrInvalidDataSize
	}

	dataSize := uint16(len(data))
	tempData := make([]byte, blockDataSize)
	for k, v := range data {
		tempData[k] = v
	}

	b := &block{
		completed:  0,
		blockCount: blockCount,
		blockIndex: blockIndex,
		blockLen:   dataSize,
		data:       tempData,
	}

	return b, nil
}

func newBlockFromBytes(bytes []byte) (*block, error) {
	if len(bytes) != int(blockDataSize+blockHeadSize) {
		return nil, ErrInvalidBlockSize
	}

	return &block{
		completed:  uint16(bytes[0]),
		blockCount: binary.BigEndian.Uint16(bytes[1:3]),
		blockIndex: binary.BigEndian.Uint16(bytes[3:5]),
		blockLen:   binary.BigEndian.Uint16(bytes[5:7]),
		data:       bytes[blockHeadSize:],
	}, nil
}

func (block *block) isCompleted(expectedBlockIndex uint16) bool {
	var finish bool = false
	if block.completed > 0 {
		finish = true
	}
	return expectedBlockIndex == block.blockIndex &&
		block.blockCount <= maxBlockCount &&
		block.blockIndex < block.blockCount &&
		block.blockLen <= blockDataSize &&
		finish
}

func (block *block) serialize() []byte {
	ret := make([]byte, int(blockDataSize+blockHeadSize))
	if block.completed > 0 {
		ret[0] = 1
	} else {
		ret[0] = 0
	}
	binary.BigEndian.PutUint16(ret[1:3], block.blockCount)
	binary.BigEndian.PutUint16(ret[3:5], block.blockIndex)
	binary.BigEndian.PutUint16(ret[5:7], block.blockLen)
	for index, bit := range block.data {
		ret[int(blockHeadSize)+index] = bit
	}

	return ret
}
