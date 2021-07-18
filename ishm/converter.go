package ishm

import (
	"encoding/json"
	"errors"
)

var converters map[string]Converter

func init() {
	converters = make(map[string]Converter, 10)
	RegisterConverter("default", DefaultConverter{})
}

//RegisterConverter ...
func RegisterConverter(codec string, c Converter) error {
	if _, ok := converters[codec]; ok {
		return errors.New("Converter Already Exist")
	}
	converters[codec] = c
	return nil
}

// Encode ...
func Encode(v interface{}, codec ...string) ([]byte, error) {
	c := "default"
	if len(codec) > 0 {
		c = codec[0]
	}
	if _, ok := converters[c]; !ok {
		return []byte{}, errors.New("Unsupport codec")
	}
	return converters[c].Marshal(v)
}

// Decode ...
func Decode(data []byte, v interface{}, codec ...string) error {
	c := "default"
	if len(codec) > 0 {
		c = codec[0]
	}
	if _, ok := converters[c]; !ok {
		return errors.New("Unsupport codec")
	}
	return converters[c].Unmarshal(data, v)
}

// Converter is the interface to convert obj to []byte
type Converter interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

//DefaultConverter is the default Converter using reflect
type DefaultConverter struct{}

// Marshal ...
func (df DefaultConverter) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (df DefaultConverter) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
