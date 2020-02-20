package shm_test

import (
	"fmt"
	"oss/model"
	"oss/utils/shm/v2"
	"testing"
)

func testConstructor() interface{} {
	return &model.LogItem{}
}

func TestNewShm(t *testing.T) {
	s, err := shm.NewShm(6, 10000, testConstructor)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 10; i++ {
		item := &model.LogItem{
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
