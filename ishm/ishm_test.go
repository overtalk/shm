package ishm

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestStruct struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func TestStructWriteRead(t *testing.T) {
	fmt.Println("Test Struct RW")
	orig := &TestStruct{
		Name: "Tencent",
		Desc: "SIXSIXSIX",
	}
	jsonObj, _ := json.Marshal(orig)
	fmt.Println(string(jsonObj))
	var total int64
	total = 1024
	segment, err := Create(total)
	defer segment.Destroy()
	if err != nil {
		t.Errorf("Failed to allocate 1024b segment: %v", err)
		return
	}
	index := len(jsonObj)
	n, err := segment.Write(jsonObj)
	if err != nil {
		return
	}
	fmt.Println("Real size: ", len(jsonObj), ", Wrote size:", n)
	orig.Desc = "666"
	jsonObj, _ = json.Marshal(orig)
	newLen := len(jsonObj)
	n, err = segment.Write(jsonObj)
	if err != nil {
		return
	}
	recoverData := make([]byte, newLen, newLen)
	fmt.Println("Real size: ", newLen, ", Wrote size:", n)
	segment.Seek(int64(index), 0)
	segment.Read(recoverData)
	fmt.Println("Get ", string(recoverData))
	fmt.Println("Real size: ", len(jsonObj), ", Recover size:", len(recoverData))
	recoverObj := &TestStruct{}
	json.Unmarshal(recoverData, recoverObj)
	fmt.Printf("%s is %s\n", recoverObj.Name, recoverObj.Desc)

	recoverData = make([]byte, index, index)
	segment.Seek(int64(0), 0)
	segment.Read(recoverData)
	fmt.Println("Get ", string(recoverData))
	fmt.Println("Real size: ", len(jsonObj), ", Recover size:", len(recoverData))
	recoverObj = &TestStruct{}
	json.Unmarshal(recoverData, recoverObj)
	fmt.Printf("%s is %s\n", recoverObj.Name, recoverObj.Desc)
}
