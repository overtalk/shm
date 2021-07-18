package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/kevinu2/shm/ishm"
	"github.com/kevinu2/shm/shmdata"
	"log"
)

const MAX_SIZE= 1 << 30

func ReadSHMI()  {
	sm,err:=ishm.CreateWithKey(12,MAX_SIZE)
	if err != nil {
		log.Fatal(err)
		sm.Destroy()
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf) // will write to buf

	shmi:=shmdata.SHMInfo{}
	shmi.MaxSHMSize=100
	shmi.MaxContentLen=64
	shmi.MaxTopicLen=128
	shmi.Count=4
	shmi.Key[0]=1000
	shmi.Key[1]=20000
	shmi.Key[2]=30000
	fmt.Printf("shm org:%#v\n",shmi)
	encoder.Encode(shmi)

	sm.Write(buf.Bytes())
	shmilen:=buf.Len()//unsafe.Sizeof(shmi)
	fmt.Printf("sizeof:%v\n",shmilen)
	od,err:=sm.ReadChunk(int64(shmilen),0 )
	if err != nil {
		log.Fatal(err)
	}
	buf.Reset()
	buf.Write(od)
	decoder := gob.NewDecoder(&buf) // will read from buf
	smrd:=shmdata.SHMInfo{}
	err=decoder.Decode(&smrd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("shm read:%#v\n",smrd)
	fmt.Printf("sm:%#v\n",sm)
	sm.Destroy()
}
func main() {
	ReadSHMI()
}
