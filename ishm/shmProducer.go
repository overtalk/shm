package ishm

import (
	"log"
	"reflect"
	"unsafe"
)


type CreateSHMParam struct {
	Key int64
	Size int64
}
type UpdateContent struct {
	EventType int16
	Topic string
	Content string
}
func StringToByteArr(s string, arr []byte) {
	src := []rune(s)
	for i, v := range src {
		if i >= len(arr) {
			break
		}
		arr[i] = byte(v)
	}
}

var sizeOfTagTLVStruct = int(unsafe.Sizeof(TagTLV{}))

func TagTLVStructToBytes(s *TagTLV) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfTagTLVStruct
	x.Cap = sizeOfTagTLVStruct
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func BytesToTagTLVStruct(b []byte) *TagTLV {
	return (*TagTLV)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}

func UpdateCtx(shmparam CreateSHMParam, updatectx UpdateContent) (index int, err error){

	log.Printf("UpdateCtx:%#v,%#v",shmparam,updatectx)
	updateSHMInfo(999999,shmparam.Key)

	tlv:=TagTLV{}
	if shmparam.Size < int64(unsafe.Sizeof(tlv)) {
		shmparam.Size =int64(unsafe.Sizeof(tlv))
	}
	sm,err:= CreateWithKey(shmparam.Key,shmparam.Size)
	if err != nil {
		log.Fatal(err)
		return index,err
	}
	tlv.Tag=1
	tlv.Len=uint64(len(updatectx.Content))
	StringToByteArr(updatectx.Topic,tlv.Topic[:])
	StringToByteArr(updatectx.Content,tlv.Value[:])
	wd:=TagTLVStructToBytes(&tlv)
	sm.Write(wd)

	if err != nil {
		log.Fatal(err)
	}
	return int(sm.Id),err
}
func GetCtx(shmparam CreateSHMParam)   (*UpdateContent,error){
	log.Printf("GetCtx:%#v",shmparam)
	sm,err:= CreateWithKey(shmparam.Key,0)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	log.Print(sm)
	data:=make([]byte,sizeOfTagTLVStruct)
	pos ,err:=sm.Read(data)
	if err !=nil {
		log.Fatal(err)
		return nil,err
	}
	log.Println(pos)
	tlv:=BytesToTagTLVStruct(data)
	ctd:=new(UpdateContent)
	ctd.Topic=string(tlv.Topic[:])
	ctd.Content=string(tlv.Value[:])

	return ctd,nil
}
