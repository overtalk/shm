package ishm

import (
	"errors"
	"fmt"
	"sync"
)

//NewShmManager ...
func NewShmManager(MemSize int64) *ShmManager {
	return &ShmManager{
		memSize: MemSize,
	}
}

// ShmManager is a tool to manager shm
type ShmManager struct {
	memSize    int64
	usedSize   int64
	freeSize   int64
	usedBlocks map[string]*MemBlock
	memChain   *MemChain
	segment    *Segment
	mu         sync.RWMutex
}

// Init ...
func (sm *ShmManager) Init() (err error) {
	sm.usedSize = 0
	sm.freeSize = sm.memSize
	sm.usedBlocks = make(map[string]*MemBlock, 100)
	sm.memChain = NewMemChain()
	sm.memChain.Insert(NewBlock(0, sm.memSize))
	sm.mu = sync.RWMutex{}
	sm.segment, err = Create(sm.memSize)
	return err
}

// WriteBlock ...
func (sm *ShmManager) WriteBlock(blockName string, data []byte) (int, error) {
	block, err := sm.memChain.SearchBlock(int64(len(data)))
	if err != nil {
		return 0, err
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.usedBlocks[blockName] = block
	pos, _ := sm.segment.Seek(block.start, 0)
	iPrintln("Goto ", pos)
	return sm.segment.Write(data)
}

// DeleteBlock ...
func (sm *ShmManager) DeleteBlock(blockName string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	block, ok := sm.usedBlocks[blockName]
	if !ok {
		return errors.New("No such block: " + blockName)
	}
	sm.memChain.Insert(block)
	delete(sm.usedBlocks, blockName)
	sm.memChain.MergeBlocks()
	return nil
}

// ReadBlock ...
func (sm *ShmManager) ReadBlock(blockName string) ([]byte, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	block, ok := sm.usedBlocks[blockName]
	if !ok {
		return []byte{}, errors.New("No such block: " + blockName)
	}
	iPrintf("Block(%d->%d):%d\n", block.start, block.end, block.Size())
	data := make([]byte, block.Size(), block.Size())
	sm.segment.Seek(block.start, 0)
	_, err := sm.segment.Read(data)
	return data, err
}

// Show ...
func (sm *ShmManager) Show() {
	fmt.Println("Used Blocks")
	for k, v := range sm.usedBlocks {
		fmt.Print(k, ": ", v.start, ",", v.end, "\n")
	}
	fmt.Println("\nFree Blocks")
	sm.memChain.PrintChain()
	fmt.Println()
}
