package shm

import (
	"bytes"
	"encoding/gob"
	"github.com/overtalk/shm/queue"
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

func NewSingleShm(key, size int, constructor ConstructorFunc) (*SHM, error) {
	s, err := queue.NewRingQueue(key, size)
	if err != nil {
		return nil, err
	}

	return &SHM{
		shmQueue:    s,
		constructor: constructor,
	}, nil
}

func NewMultiShm(key, size int, constructor ConstructorFunc) (*SHM, error) {
	s, err := queue.NewMultiQueue(key, size)
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
