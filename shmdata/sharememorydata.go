package shmdata

import (
	"fmt"
	"github.com/kevinu2/shm/ishm"
	"log"
	"reflect"
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
	Tag      uint64
	Len      uint64
	Topic    []byte
	Value    []byte
}


//todo  run it use root
func GetShareMemoryInfo(defaultKey int64) (*SHMInfo, error) {
	shmi:=SHMInfo{}
	sm,err:=ishm.CreateWithKey(defaultKey,0)
	if err != nil {
		log.Fatal(err)
		sm.Destroy()
	}
	od,err:=sm.ReadChunk(int64(unsafe.Sizeof(shmi)),0 )
	if err != nil {
		log.Fatal(err)
	}
	data :=  *(*[]byte)(unsafe.Pointer(&od))
	var  readshmi *SHMInfo= *(**SHMInfo)(unsafe.Pointer(&data))
	fmt.Printf("shmiii:%#v\r\n",readshmi)
	fmt.Printf("sm:%#v\n",sm)
	return readshmi,err
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
