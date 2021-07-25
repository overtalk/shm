package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/kevinu2/shm/ishm"
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

	shmi := ishm.SHMInfo{}
	lll := ishm.SizeStruct(shmi)
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
		//log.Fatal(err)
	}
	buf.Reset()
	buf.Write(od)
	decoder := gob.NewDecoder(&buf) // will read from buf
	smrd := ishm.SHMInfo{}
	err = decoder.Decode(&smrd)
	if err != nil {
		//log.Fatal(err)
	}
	fmt.Printf("shm read:%#v\n", smrd)
	fmt.Printf("sm:%#v\n", sm)
	sm.Destroy()
}
func testReadSHMByDefaultSHMI(){
	shmi, err := ishm.GetShareMemoryInfo(999999,false)
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
				ishm.Readtlv(int64(k))
			}()
		} else {
			ishm.Readtlv(int64(k))
		}
		i++
	}

}
type TestJsonData struct {
	Name string `json:"name"`
	DataLength int `json:"dataLength"`
	Content string `json:"content"`
}
func testProducer()  {
	td:=TestJsonData{"x要针对哪个 proto 文件生成接口代码xx",12,"yy要针对哪个 proto 文件生成接口代码yy"}
	od,err:=json.Marshal(td)
	if err != nil {
		log.Fatal(err)
	}
	shmParam:= ishm.CreateSHMParam{4567,2000,true}
	ctx:= ishm.UpdateContent{EventType: "data-event",Topic: "xxx",Content: string(od)}
	ishm.UpdateCtx(shmParam,ctx)
	readDataFromSHM,err:= ishm.GetCtx(shmParam)



	if err !=nil {

	}else {
		log.Println("read data form shm is:%#v",readDataFromSHM)
	}

	var counter int = 0
	for  {

		shareshminfo,err:=ishm.GetShareMemoryInfo(999999,false)
		if err !=nil {

		}
		log.Print(shareshminfo)
		shmParam.Create=false
		ishm.UpdateCtx(shmParam,ctx)
		readDataFromSHM,err= ishm.GetCtx(shmParam)
		if err !=nil {

		}else {
			log.Println("read data form shm is:%#v",readDataFromSHM)
		}
		counter++
		if counter > 10 {
			 break
		}

	}


}
func main() {

	testProducer()
}
