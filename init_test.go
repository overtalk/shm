package shm_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/kevinu2/shm"
)

// LogItem defines one log record
type LogItem struct {
	ProtocolName string
	Fields       []string
	Data         []interface{}
}

func LogItemConstructor() interface{} {
	return &LogItem{}
}

func TestSingleShm(t *testing.T) {
	mem, err := shm.NewSystemVMem(6, 10000)
	if err != nil {
		log.Fatal(err)
	}

	s, err := shm.NewSingleShm(mem, 10000, LogItemConstructor)
	if err != nil {
		t.Error(err)
		return
	}

	details(t, s)
}

func TestMultiShm(t *testing.T) {
	mem, err := shm.NewSystemVMem(7, 10000)
	if err != nil {
		log.Fatal(err)
	}

	s, err := shm.NewMultiShm(mem, 10000, LogItemConstructor)
	if err != nil {
		t.Error(err)
		return
	}
	details(t, s)
}

func details(t *testing.T, s *shm.SHM) {
	for i := 0; i < 10; i++ {
		item := &LogItem{
			ProtocolName: "1",
			Fields:       []string{fmt.Sprintf("field-%d", i)},
			Data:         []interface{}{i},
		}
		if err := s.Save(item); err != nil {
			t.Error(err)
			return
		}
	}

	items, err := s.Get()
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range items {
		fmt.Printf("value : %v, type = %T\n", v, v)
	}
}
