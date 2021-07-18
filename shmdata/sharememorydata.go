package shmdata

import (
	"errors"
	"fmt"
	"github.com/kevinu2/shm/ishm"
	"log"
	"reflect"
	"time"
	"unsafe"
)

//todo this will be remove after test ok ,this code form  C
//typedef struct SHMInfo
//{
//unsigned long long max_topic_len;
//unsigned long long max_content_len;
//unsigned long long max_shm_size;
//unsigned long long count;
//key_t key[200];
//} SHMI, *PSHMI;

//todo this will read from  key of 999999  SHM
//typedef struct SHMInfo{
//unsigned long long max_topic_len;
//unsigned long long max_content_len;
//unsigned long long max_shm_size;
//unsigned long long count;
//key_t key[200];
//}SHMI;
type SHMInfo struct {
	MaxTopicLen   uint64
	MaxContentLen uint64
	MaxSHMSize    uint64
	Count         uint64
	Key           [200]int32
}

var MTL uint = 64
var MCL uint = 102400

type TagTLV struct {
	Tag   uint64
	Len   uint64
	Topic []byte
	Value []byte
}
type TagTL struct {
	Tag   uint64
	Len   uint64
}
type HeadData struct {
	ReadOffSet  uint64
	WriteOffSet uint64
}

func GetHeadData(segment *ishm.Segment) (*HeadData, error) {
	h := HeadData{}
	od, err := segment.ReadChunk(int64(unsafe.Sizeof(h)), 0)
	if err != nil {
		log.Fatal(err)
	}
	data := *(*[]byte)(unsafe.Pointer(&od))
	var hd *HeadData = *(**HeadData)(unsafe.Pointer(&data))
	fmt.Printf("hd:%#v\r\n", hd)
	fmt.Printf("sm:%#v\n", segment)
	return hd, err
}
func ReadTLVData(segment *ishm.Segment,offset int64) (*TagTLV, int64,error) {
	tl := TagTL{}
	var retOffset int64=offset
	od, err := segment.ReadChunk(int64(unsafe.Sizeof(tl)), offset)
	if err != nil {
		log.Fatal(err)
	}
	data := *(*[]byte)(unsafe.Pointer(&od))
	var tll *TagTL = *(**TagTL)(unsafe.Pointer(&data))
	fmt.Printf("tll:%#v\r\n", tll)
	if tll.Len==0 {
		return nil,16,errors.New("data is end")
	}
	tlv := TagTLV{}
	tlv.Topic = make([]byte, 64)
	tlv.Value = make([]byte, tll.Len)
	datalen := SizeStruct(tlv)
	od, err = segment.ReadChunk(int64(datalen), offset)
	if err != nil {
		log.Fatal(err)
	}
	data = *(*[]byte)(unsafe.Pointer(&od))
	var readtlv *TagTLV = *(**TagTLV)(unsafe.Pointer(&data))
	retOffset+=int64(datalen)
	fmt.Printf("tlv:T %v Len :%v\r\n",readtlv.Tag,readtlv.Len)
	return readtlv, retOffset,err
}
func Readtlv(k int64)  {
	sm, err := ishm.CreateWithKey(int64(k), 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(sm)
	var offset int64 = 16
	//	for {
	hd, err := GetHeadData(sm)
	if err == nil {
		fmt.Println(hd)
	}
	tlv, retoffset, err := ReadTLVData(sm, offset)
	fmt.Printf("tlv:Tag %v,Len %v\r\n", tlv.Tag, tlv.Len)
	fmt.Printf("offset:%v\r\n", retoffset)
	T1 := time.Now()
	for {

		tlv, retoffset, err = ReadTLVData(sm, retoffset)
		fmt.Printf("offset:%v\r\n", retoffset)
		if err != nil {
			retoffset = 16
			T2 := time.Now()
			log.Printf("key = %v use time %v \r\n", k, T2.Sub(T1).Seconds())
			time.Sleep(time.Second * 2)
			T1 = time.Now()

		}
	}
}

//todo  run it use root
func GetShareMemoryInfo(defaultKey int64) (*SHMInfo, error) {
	shmi := SHMInfo{}
	sm, err := ishm.CreateWithKey(defaultKey, 0)
	if err != nil {
		log.Fatal(err)
		sm.Destroy()
	}
	od, err := sm.ReadChunk(int64(unsafe.Sizeof(shmi)), 0)
	if err != nil {
		log.Fatal(err)
	}
	data := *(*[]byte)(unsafe.Pointer(&od))
	var readshmi *SHMInfo = *(**SHMInfo)(unsafe.Pointer(&data))
	fmt.Printf("shmiii:%#v\r\n", readshmi)
	fmt.Printf("sm:%#v\n", sm)
	return readshmi, err
}

func SizeStruct(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Map:
		sum := 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Slice, reflect.Array:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.String:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Ptr, reflect.Interface:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			return 0
		}
		return sizeof(v.Elem())
	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int(v.Type().Size())

	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}

	return -1
}
