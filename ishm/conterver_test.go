package ishm

import (
	"fmt"
	"testing"
)

func TestDefaultConverter(t *testing.T) {
	codec := "default"
	RegisterConverter(codec, DefaultConverter{})
	ts := &TestStruct{
		Name: "Tencent",
		Desc: "666",
	}
	data, _ := Encode(ts)
	fmt.Println(string(data))
	tsRecv := &TestStruct{}
	Decode(data, tsRecv)
	fmt.Printf("%s is %s\n", tsRecv.Name, tsRecv.Desc)
}
