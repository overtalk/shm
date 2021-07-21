# Golang 共享内存

- `golang` 中使用共享内存

## 简介
- 支持的共享内存类型
    - systemV 
    - mmap
    
- 两种应用场景
    - 只可以是一个生产者 & 一个消费者（NewSingleShm 构造函数）
    - 可以 多个生产者 & 一个消费者（NewMultiShm 构造函数）
    
- 具体的实现思路是基于[范健的这篇分享](https://cloud.tencent.com/developer/article/1006241)，在此也感谢作者
- 上述分享主要是一个无锁队列的实现，我在此基础上套上了一层共享内存
  
## 安装
```bash
go get github.com/overtalk/shm
```

## 使用
### 创建共享内存块
```go
// system V
mem, err := shm.NewSystemVMem(6, 10000)
if err != nil {
    log.Fatal(err)
}

// mmap
mem, err := shm.NewMMapMem("./test.txt", 10000)
if err != nil {
    log.Fatal(err)
}
```

### 使用
```go
package main

import (
	"fmt"
	"log"

	"github.com/overtalk/shm"
)

type LogItem struct {
	ProtocolName string
	Fields       []string
	Data         []interface{}
}

func testConstructor() interface{} {
	return &LogItem{}
}

func main() {
    // 构造共享内存块
    // 如果需要使用 mmap 共享内存，使用 NewMMapMem 构造方法即可
	mem, err := shm.NewSystemVMem(6, 10000)
	if err != nil {
		log.Fatal(err)
	}

	s, err := shm.NewMultiShm(mem, 10000, testConstructor)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		item := &LogItem{
			ProtocolName: "1",
			Fields:       []string{fmt.Sprintf("field-%d", i)},
			Data:         []interface{}{i},
		}
		if err := s.Save(item); err != nil {
			fmt.Println(err)
			return
		}
	}

	items, err := s.Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range items {
		fmt.Printf("value : %v, type = %T\n", v, v)
	}
}
```

### 使用订阅使用
```bigquery
package main

import (
	"fmt"
	"log"

	"github.com/overtalk/shm"
)

func callBack(*shmdata.TagTLV){
    //do something
}

func main(){
    if shm.StartSubscribe(999999, callBack) {
        fmt.Println("start alert shm success")
    }else{
        fmt.Println("start alert shm failed")
    }
    if shm.StartSubscribe(888888, callBack) {
        fmt.Println("start log shm success")
    }else{
        fmt.Println("start log shm failed")
    }
}
```
