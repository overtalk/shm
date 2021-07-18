package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/kevinu2/shm/ishm"
	"github.com/kevinu2/shm/shmdata"
	"log"
	"math/rand"
	"unsafe"
)

const MAX_SIZE = 1 << 30

func WriteReadSHMI() {
	sm, err := ishm.CreateWithKey(12, MAX_SIZE)
	if err != nil {
		log.Fatal(err)
		sm.Destroy()
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf) // will write to buf

	shmi := shmdata.SHMInfo{}
	lll := shmdata.SizeStruct(shmi)
	fmt.Printf("shmisize:%v ,,,sizof:%v\n", lll, unsafe.Sizeof(shmi))
	shmi.MaxSHMSize = 100
	shmi.MaxContentLen = 64
	shmi.MaxTopicLen = 128
	shmi.Count = 4
	for i, _ := range shmi.Key {
		shmi.Key[i] = rand.Int31()
	}
	fmt.Printf("shm org:%#v\n", shmi)
	encoder.Encode(shmi)
	sm.Write(buf.Bytes())
	shmilen := buf.Len() //unsafe.Sizeof(shmi)

	fmt.Printf("buferlen:%v\n", shmilen)
	od, err := sm.ReadChunk(int64(shmilen), 0)
	if err != nil {
		log.Fatal(err)
	}
	buf.Reset()
	buf.Write(od)
	decoder := gob.NewDecoder(&buf) // will read from buf
	smrd := shmdata.SHMInfo{}
	err = decoder.Decode(&smrd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("shm read:%#v\n", smrd)
	fmt.Printf("sm:%#v\n", sm)
	sm.Destroy()
}
func main() {
	shmi, err := shmdata.GetShareMemoryInfo(999999)
	if err != nil {
		log.Fatal(err)
	}
	for i, k := range shmi.Key {
		if i == int(shmi.Count) {
			break
		}
		fmt.Printf("key:%v\r\n", k)
		if uint64(i) < shmi.Count-1 {
			go func() {
				shmdata.Readtlv(int64(k))
			}()
		} else {
			shmdata.Readtlv(int64(k))
		}
		i++
	}

}
