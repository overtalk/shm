package main

import (
	"fmt"
	"github.com/overtalk/shm"
)

type LogItem struct {
	ProtocolName string
	Fields       []string
	Data         []interface{}
}

func testConstructor() interface{} {
	return &LogItem{}
}

func main() {
	s, err := shm.NewMultiShm(6, 10000, testConstructor)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		item := &LogItem{
			ProtocolName: "1",
			Fields:       []string{fmt.Sprintf("field-%d", i)},
			Data:         []interface{}{i},
		}
		if err := s.Save(item); err != nil {
			fmt.Println(err)
			return
		}
	}

	items, err := s.Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range items {
		fmt.Printf("value : %v, type = %T\n", v, v)
	}
}
