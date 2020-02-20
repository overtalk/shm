package shm

import (
	"bytes"
	"encoding/gob"
)

type ConstructorFunc func() interface{}

type SHM struct {
	s           *shm
	constructor ConstructorFunc
}

func NewShm(key, size int, constructor ConstructorFunc) (*SHM, error) {
	s, err := newShm(key, size)
	if err != nil {
		return nil, err
	}

	return &SHM{
		s:           s,
		constructor: constructor,
	}, nil
}

func (s *SHM) Save(i ...interface{}) error {
	var toSave []byte
	for _, v := range i {
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(v)
		if err != nil {
			return err
		}

		toSave = append(toSave, newBinaryMessage(buf.Bytes()).serialize()...)
	}

	return s.s.save(toSave)
}

func (s *SHM) Get() ([]interface{}, error) {
	var ret []interface{}

	binaryMessages, err := deserializeSlice(s.s.get())
	if err != nil {
		return nil, err
	}

	for _, bm := range binaryMessages {
		temp := s.constructor()
		if err = gob.NewDecoder(bytes.NewBuffer(bm.Body)).Decode(temp); err != nil {
			return nil, err
		}

		ret = append(ret, temp)
	}

	return ret, nil
}
