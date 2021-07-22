package shmdata

import (
	"errors"
	"fmt"
	"log"
	"github.com/kevinu2/shm/ishm"
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
var Counter int64 = 0

type TagTLV struct {
	Tag   int64
	Len   uint64
	TopicLen uint16
	EventTypeLen uint16
	Topic [30]byte
	EventType [30]byte
	Value [40960]byte
}

type ContentData struct {
	Tag   int64
	Topic  string
	Value string
}
type TagTL struct {
	Tag uint64
	Len uint64
}
type HeadData struct {
	ReadOffSet  uint64
	WriteOffSet uint64
}
var ST1 time.Time
var ST2 time.Time
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

func ReadTLVData(segment *ishm.Segment, offset int64) (*TagTLV, int64, error) {
	tl := TagTL{}
	var retOffset int64 = offset
	od, err := segment.ReadChunk(int64(unsafe.Sizeof(tl)), offset)
	if err != nil {
		log.Fatal(err)
	}
	data := *(*[]byte)(unsafe.Pointer(&od))
	var tll *TagTL = *(**TagTL)(unsafe.Pointer(&data))
	now := time.Now().Unix()
	last := time.Unix(int64(tll.Tag), 0)

	if now > last.Unix() {
		log.Printf("will be wait ,now: %d,last:%d\n", now, last.Unix())
		//time.Sleep(time.Second)
		//return nil,retOffset,nil
	}

	if tll.Len == 0 {
		return nil, 16, errors.New("data is end")
	}
	tlv := TagTLV{}
	datalen := int64(unsafe.Sizeof(tl)) + int64(64) + int64(tll.Len) // SizeStruct(tlv)
	od, err = segment.ReadChunk(datalen, offset)
	if err != nil {
		log.Fatal(err)
	}
	data = *(*[]byte)(unsafe.Pointer(&od))
	var readtlv *TagTLV = &tlv
	readtlv = *(**TagTLV)(unsafe.Pointer(&data))
	retOffset += int64(datalen)

	topic := string(readtlv.Topic[:])
	fmt.Printf("topic:%s\n", topic)
	content := string(readtlv.Value[:])
	//fmt.Printf("content:%s\n", content)
	Counter++
	ST2=time.Now()
	if ST2.Sub(ST1).Seconds() >10.000 {
		log.Printf("data %v per sec\r\n",float64(Counter)/ST2.Sub(ST1).Seconds())
		Counter = 0
		ST1=time.Now()
	}
	contentData:=ContentData{}
	contentData.Tag=readtlv.Tag
	contentData.Topic=topic
	contentData.Value=content

	return readtlv, retOffset, err
}
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func Readtlv(k int64) {
	sm, err := ishm.CreateWithKey(int64(k), 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(sm)
	var offset int64 = 16
	hd, err := GetHeadData(sm)
	if err == nil {
		fmt.Println(hd)
	}
	tlv, retoffset, err := ReadTLVData(sm, offset)
	if tlv != nil {
		fmt.Printf("tlv:Tag %v,Len %v\r\n", tlv.Tag, tlv.Len)
	}
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

	ST1 = time.Now()
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
