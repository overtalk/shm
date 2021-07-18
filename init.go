package shm

import (
	"bytes"
	"encoding/gob"

	"github.com/kevinu2/shm/model"
	"github.com/kevinu2/shm/queue"
)

type ConstructorFunc func() interface{}

type shmQueue interface {
	Save(buf []byte) error
	Get() ([][]byte, error)
}

type SHM struct {
	shmQueue    shmQueue
	constructor ConstructorFunc
}

func NewSingleShm(shm *model.Mem, size int, constructor ConstructorFunc) (*SHM, error) {
	s, err := queue.NewRingQueue(shm, size)
	if err != nil {
		return nil, err
	}

	return &SHM{
		shmQueue:    s,
		constructor: constructor,
	}, nil
}

func NewMultiShm(shm *model.Mem, size int, constructor ConstructorFunc) (*SHM, error) {
	s, err := queue.NewMultiQueue(shm, size)
	if err != nil {
		return nil, err
	}

	return &SHM{
		shmQueue:    s,
		constructor: constructor,
	}, nil
}

func (s *SHM) Save(i interface{}) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(i)
	if err != nil {
		return err
	}

	return s.shmQueue.Save(buf.Bytes())
}

func (s *SHM) Get() ([]interface{}, error) {
	data, err := s.shmQueue.Get()
	if err != nil {
		return nil, err
	}

	var ret []interface{}

	for _, v := range data {
		temp := s.constructor()
		if err = gob.NewDecoder(bytes.NewBuffer(v)).Decode(temp); err != nil {
			return nil, err
		}

		ret = append(ret, temp)
	}

	return ret, nil
}

func (s *SHM) GetByIndex(index int) (interface{}, error) {
	data, err := s.shmQueue.Get()
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		temp := s.constructor()
		err = gob.NewDecoder(bytes.NewBuffer(v)).Decode(temp)
		if err != nil {
			return nil, err
		}

		if index == k {
			return temp, err
		}
	}

	return nil, nil
}
