package main

import (
	"fmt"
	"github.com/kevinu2/shm"
	"github.com/kevinu2/shm/shmdata"
	"log"
)

//type LogItem struct {
//	ProtocolName string
//	Fields       []string
//	Data         []interface{}
//}
//go build -o readshm example/shm/main.go
//func testConstructor() interface{} {
//	return &LogItem{}
//}
func testConstructor() interface{} {
	return &shmdata.TagTLV{}
}

//todo run this please run del-shm.sh
func main() {

	//shmi, err := shmdata.GetShareMemoryInfo(999999)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(shmi)

	mem, err := shm.NewSystemVMem(6, 10000)
	if err != nil {
		log.Fatal(err)
	}

	s, err := shm.NewMultiShm(mem, 10000, testConstructor)
	if err != nil {
		fmt.Println(err)
		return
	}

	//for i := 0; i < 10; i++ {
	//	item := &LogItem{
	//		ProtocolName: "1",
	//		Fields:       []string{fmt.Sprintf("field-%d", i)},
	//		Data:         []interface{}{i},
	//	}
	//	if err := s.Save(item); err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}

	items, err := s.Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range items {
		fmt.Printf("value : %v, type = %T\n", v, v)
	}
}
