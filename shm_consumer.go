package shm

import (
	"fmt"
	"ndr_forward/src/shm/ishm"
	"ndr_forward/src/shm/shmdata"
	"unsafe"
)

type ShmConsumerStatus int32

const(
	ShmConsumerOk ShmConsumerStatus = 0
	ShmConsumerReadErr ShmConsumerStatus = 1
	ShmConsumerLenErr ShmConsumerStatus = 2
	ShmConsumerInitErr ShmConsumerStatus = 3
	ShmConsumerNoData ShmConsumerStatus = 4
)

//after//shmi, err := shmdata.GetShareMemoryInfo(999999)

type Consumer struct {
	PreTag uint64
	CurTag uint64
	TopicLen uint64
	MaxContentLen uint64
	MaxShmSize uint64
	CurOffset uint64
	PreOffset uint64
	SegLen uint64
	ShmKey int64
	IsRunning bool
	sm *ishm.Segment
}

func (consumer *Consumer)Init(key int64, maxSHMSize uint64, maxContentLen uint64) bool {
	consumer.ShmKey = key
	sm, err := ishm.CreateWithKey(key, 0)
	if err != nil {
		fmt.Printf("Init consume err key-%v\n\n", key)
		return false
	}
	consumer.MaxShmSize = maxSHMSize
	consumer.MaxContentLen = maxContentLen
	consumer.SegLen = maxContentLen + 80
	consumer.CurOffset = 16
	consumer.sm = sm
	consumer.IsRunning = false
	return true
}

func (consumer *Consumer)Reset() {
	consumer.CurOffset = 16
	consumer.PreTag = 0
	consumer.CurTag = 0
	consumer.IsRunning = false
}

func (consumer *Consumer)Next() (*shmdata.TagTLV, ShmConsumerStatus){
	//tl := shmdata.TagTL{}
	//od, err := consumer.sm.ReadChunk(int64(unsafe.Sizeof(shmdata.TagTL)), int64(consumer.CurOffset))
	if consumer.sm == nil{
		return nil, ShmConsumerInitErr
	}
	od, err := consumer.sm.ReadChunk(16, int64(consumer.CurOffset))
	if err != nil {
		return nil, ShmConsumerReadErr
	}
	data := *(*[]byte)(unsafe.Pointer(&od))
	var tll *shmdata.TagTL = *(**shmdata.TagTL)(unsafe.Pointer(&data))
	if tll.Len > consumer.MaxContentLen{
		return nil, ShmConsumerLenErr
	}
	if tll.Len > 0{
		if (tll.Tag > consumer.PreTag) || (tll.Tag == 0 && consumer.PreTag == 18446744073709551615 || consumer.PreTag == 0){
			//copySize := int64(unsafe.Sizeof(tl)) + int64(64) + int64(tll.Len)
			copySize := int64(16) + int64(64) + int64(tll.Len)
			od, err = consumer.sm.ReadChunk(copySize, int64(consumer.CurOffset))
			consumer.PreTag = tll.Tag
			consumer.PreOffset = consumer.CurOffset
			consumer.CurOffset += consumer.SegLen
			if consumer.CurOffset + consumer.SegLen > consumer.MaxShmSize{
				fmt.Printf("Worker-%v new cycle\n", consumer.ShmKey)
				consumer.CurOffset = 16
			}
			consumer.IsRunning = true
			data = *(*[]byte)(unsafe.Pointer(&od))
			//readtlv = *(**shmdata.TagTLV)(unsafe.Pointer(&data))
			return *(**shmdata.TagTLV)(unsafe.Pointer(&data)), ShmConsumerOk
		}
	}else{
		if consumer.CurOffset != 16{
			//od, err := consumer.sm.ReadChunk(int64(unsafe.Sizeof(tl)), int64(consumer.CurOffset))
			od, err := consumer.sm.ReadChunk(int64(16), int64(consumer.CurOffset))
			if err != nil {
				//log.Fatal(err)
				return nil, ShmConsumerReadErr
			}
			data := *(*[]byte)(unsafe.Pointer(&od))
			var headTll *shmdata.TagTL = *(**shmdata.TagTL)(unsafe.Pointer(&data))
			if headTll.Len > 0 && (headTll.Tag > consumer.PreTag || (headTll.Tag == 0 && consumer.PreTag == 18446744073709551615)){
				//new cycle
				consumer.PreOffset = consumer.CurOffset
				fmt.Printf("Worker-%v new cycle headtag-%v\n", consumer.ShmKey, headTll.Tag)
				//copySize := int64(unsafe.Sizeof(tl)) + int64(64) + int64(headTll.Len)
				copySize := int64(16) + int64(64) + int64(headTll.Len)
				od, err = consumer.sm.ReadChunk(copySize, 16)
				consumer.PreTag = headTll.Tag
				consumer.CurOffset = consumer.SegLen + 16
				consumer.IsRunning = true
				data = *(*[]byte)(unsafe.Pointer(&od))
				//readtlv = *(**shmdata.TagTLV)(unsafe.Pointer(&data))
				return *(**shmdata.TagTLV)(unsafe.Pointer(&data)), ShmConsumerOk
			}
		}
	}
	return nil, ShmConsumerNoData
}
