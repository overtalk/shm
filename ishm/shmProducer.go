package ishm

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type CreateSHMParam struct {
	Key    int64
	Size   int64
	Create bool
}
type UpdateContent struct {
	EventType string
	Topic     string
	Content   string
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

//"ipcs -m | tail -n +4 | awk {'print $1'}"

func HasKeyofSHM(key int64) bool {

	err, find := CmdAndChangeDirToShow("ipcs", []string{"-m"}, key)
	if err != nil {

	}
	return find

}
func CmdAndChangeDirToShow(commandName string, params []string, key int64) (error, bool) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDirToFile", cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return err, false
	}
	cmd.Stderr = os.Stderr
	//cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err, false
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	var find bool = false
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		keystr := strconv.FormatInt(key, 16)
		ss := strings.Split(line, " ")
		if strings.Contains(ss[0], keystr) {
			find = true
			fmt.Println("find ......")
			break
		}
	}
	err = cmd.Wait()
	return err, find
}
func encodeBase64(data []byte) string {

	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc) // aGVsbG8gd29ybGQxMjM0NSE/JComKCknLUB+
	return sEnc
}
func decodeBase64(data string) string {

	//sEnc := base64.StdEncoding.EncodeToString([]byte(data))

	// Base64 Standard Decoding
	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return string(sDec)
	}

	fmt.Println(string(sDec)) //hello world12345!?$*&()'-@~
	return string(sDec)
}

func UpdateCtx(shmparam CreateSHMParam, updatectx UpdateContent) (index int, err error) {

	log.Printf("UpdateCtx:%#v,%#v", shmparam, updatectx)
	tlv := TagTLV{}
	if shmparam.Size < int64(unsafe.Sizeof(tlv)) {
		shmparam.Size = int64(unsafe.Sizeof(tlv))
	}

	if shmparam.Create {
		updateSHMInfo(999999, shmparam.Key)
	} else {
		shmparam.Size = 0
	}
	if HasKeyofSHM(shmparam.Key) {
		shmparam.Size = 0
	}
	sm, err := CreateWithKey(shmparam.Key, shmparam.Size)

	if err != nil {
		log.Fatal(err)
		return index, err
	}
	tlv.Tag = 1
	c:=encodeBase64([]byte(updatectx.Content))
	tlv.Len = uint64(len(c))
	e:=encodeBase64([]byte(updatectx.EventType))
	tlv.EventTypeLen = uint16(len(e))

	StringToByteArr(updatectx.Topic, tlv.Topic[:])
	StringToByteArr(c, tlv.Value[:])
	StringToByteArr(e, tlv.EventType[:])
	wd := TagTLVStructToBytes(&tlv)
	sm.Write(wd)

	if err != nil {
		log.Fatal(err)
	}
	return int(sm.Id), err
}
func GetCtx(shmparam CreateSHMParam) (*UpdateContent, error) {
	log.Printf("GetCtx:%#v", shmparam)
	sm, err := CreateWithKey(shmparam.Key, 0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	data := make([]byte, sizeOfTagTLVStruct)
	pos, err := sm.Read(data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(pos)
	tlv := BytesToTagTLVStruct(data)
	ctd := new(UpdateContent)
	t:=byteArrayToString(tlv.Topic[:], int(tlv.TopicLen))
	ctd.Topic = decodeBase64(t)
	c:=byteArrayToString(tlv.Value[:], int(tlv.Len))
	ctd.Content = decodeBase64(c)
	ctd.EventType = byteArrayToString(tlv.EventType[:], int(tlv.EventTypeLen))

	return ctd, nil
}
