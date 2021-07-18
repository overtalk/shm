package main

import (
	"fmt"
	"github.com/kevinu2/shm/ishm"
	"log"
)
const MAX_SIZE= 1 << 30
func main() {

	sm,err:=ishm.CreateWithKey(12,MAX_SIZE)
	if err != nil {
		log.Fatal(err)
		sm.Destroy()
	}

	od,err:=sm.ReadChunk(20,10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n",od)
	fmt.Printf("sm:%#v \n",sm)
	sm.Destroy()
}
