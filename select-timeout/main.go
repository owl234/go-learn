package main

// select 超时模式：同时等 channel 和 time.After，哪个先到走哪个分支。
// Go 里“不等了”的标准写法，真实扫描器里给每个端口加超时就靠它。
//
// 运行：cd select-timeout && go run .
// 试试把 Sleep(2 * time.Second) 改成 500 * time.Millisecond 再跑，看输出怎么变。

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second) // 模拟一个很慢的任务
		ch <- "扫描结果"
	}()

	select {
	case r := <-ch:
		fmt.Println("收到:", r)
	case <-time.After(time.Second):
		fmt.Println("超时了，不等了")
	}
}
