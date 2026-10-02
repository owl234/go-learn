package main

// Day 3：Go 并发入门 —— goroutine + channel
//
// 模拟端口扫描：偶数端口开放、奇数端口关闭，每个任务耗时 300ms。
// 10 个任务并发执行，总耗时约 300ms（串行则要 3s），且输出顺序不固定。
//
// 运行：cd day3 && go run .

import (
	"fmt"
	"time"
)

type PortResult struct {
	Port int
	Open bool
}

// 每个端口起一个 goroutine，扫完把结果送入 channel
func scanPort(port int, ch chan PortResult) {
	time.Sleep(300 * time.Millisecond) // 模拟网络延迟
	ch <- PortResult{Port: port, Open: port%2 == 0}
}

func main() {
	start := time.Now()

	ports := []int{21, 22, 23, 80, 443, 8080, 3306, 6379, 27017, 9200}
	ch := make(chan PortResult, len(ports)) // 带缓冲：发送方不用等接收方

	for _, port := range ports {
		go scanPort(port, ch) // go 语句非阻塞：创建 goroutine 后立即返回
	}

	openCount, closedCount := 0, 0
	for range ports { // 收满 N 个结果 = 全部任务完成，天然等待
		r := <-ch
		if r.Open {
			openCount++
		} else {
			closedCount++
		}
		fmt.Printf("端口 %d：开放=%v\n", r.Port, r.Open)
	}

	fmt.Printf("开放 %d 个，关闭 %d 个，总耗时 %v\n", openCount, closedCount, time.Since(start))
}
