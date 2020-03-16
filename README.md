# Golang 共享内存

- `golang` 中使用共享内存

## 简介
- 本仓库使用的是 `shm` 共享内存而非 `mmap`
- 本仓库提供了两个共享内存的构造函数
    - NewSingleShm ： 只可以是一个生产者 & 一个消费者
    - NewMultiShm ： 可以 多个生产者 & 一个消费者
- 具体的实现思路是基于[范健的这篇分享](https://cloud.tencent.com/developer/article/1006241)，在此也感谢作者
- 上述分享主要是一个无锁队列的实现，我在此基础上套上了一层共享内存
  
## 安装
```bash
go get github.com/overtalk/shm
```

## 使用
```go
package main

import (
	"fmt"
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
	s, err := shm.NewMultiShm(6, 10000, testConstructor)
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
