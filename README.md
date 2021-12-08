# gonal
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/gonal/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/gonal/branch/main/graph/badge.svg?token=XRI6I1W0C3)](https://codecov.io/gh/czasg/gonal)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/gonal.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/gonal/stargazers)

**Gonal** 通常用于信号的异步通知与回调。

与`os.Signal`（某种特定信号）不同的是，在 gonal 中为每一个元素赋值了标签属性。

这有点类似对象存储，为每一个回调函数绑定属性，然后发送消息时，仅仅需要指定某个属性，就可以通知该属性绑定的所有元素。

# demo
```golang
package main

import (
	"context"
	"fmt"
	"github.com/czasg/gonal"
	"time"
)

func callback(ctx context.Context, payload gonal.Payload) {
	fmt.Println(payload)
}

func main() {
	label1 := gonal.Label{
		"key1": "value1",
		"key2": "value2",
	}
	gonal.Bind(label1, callback) // 绑定属性
	for {
		{
			_ = gonal.Notify(gonal.Payload{
				Label: gonal.Label{"key1": "value1"},
				Body: []byte{},
			})
			time.Sleep(time.Second)
		}
		{
			_ = gonal.Notify(gonal.Payload{
				Label: gonal.Label{"key2": "value2"},
				Body: []byte{},
			})
			time.Sleep(time.Second)
		}
	}
}
```

# more
* SetContext
设置上下文，每一个回调函数都会接收到此上下文。
```go
func init() {
    gonal = &Gonal{
        LabelsMatcher: map[string][]Handler{},
        C:             make(chan struct{}, 1),
    }
    _ = SetContext(context.Background())
}
```
* SetConcurrent
设置最大并发线程
```go
func init() {
    gonal = &Gonal{
        LabelsMatcher: map[string][]Handler{},
        C:             make(chan struct{}, 1),
    }
    _ = SetConcurrent(runtime.NumCPU() * 4)
}
```
* SetQueue
设置消息队列，默认时内存，可以指定为持久化队列，参考 `Queue`。
```go
func init() {
    gonal = &Gonal{
        LabelsMatcher: map[string][]Handler{},
        C:             make(chan struct{}, 1),
    }
    _ = SetQueue(queue.NewFifoMemoryQueue(1024))
}
```