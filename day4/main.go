package main

// Day 4：实战 —— 真 TCP 端口扫描器（worker 池）
//
// 用 net.DialTimeout 真连端口：连上=开放，connection refused=关闭，超时=无响应。
// 固定 100 个 worker 从 jobs 队列领任务，避免 65535 个端口起 65535 个 goroutine。
// 固定顺序：起工人 → 送任务 → close 打烊 → 收结果。
//
// 运行：cd day4 && go run .
// 本机 1-1024 端口约 10ms 扫完（connection refused 瞬间返回）；
// 扫远程主机时，有超时的端口会慢得多，worker 池的价值才真正体现。

import (
	"fmt"
	"net"
	"time"
)

type PortResult struct {
	Port int
	Open bool
}

// worker 从 jobs 取端口，扫完把结果送入 results；
// jobs 关闭且取空后，for range 自动退出，worker 下班。
func worker(jobs <-chan int, results chan<- PortResult) {
	for port := range jobs {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err != nil {
			results <- PortResult{Port: port, Open: false}
			continue // 这个端口失败就取下一个；return 会结束整个 worker！
		}
		conn.Close() // err == nil 时 conn 才有效，用完随手关闭
		results <- PortResult{Port: port, Open: true}
	}
}

func main() {
	const numWorkers = 100
	const maxPort = 1024

	jobs := make(chan int, numWorkers)
	results := make(chan PortResult, maxPort) // 缓冲 ≥ 任务数，否则发送方阻塞死锁

	start := time.Now()

	for w := 0; w < numWorkers; w++ {
		go worker(jobs, results)
	}
	for p := 1; p <= maxPort; p++ {
		jobs <- p
	}
	close(jobs) // 打烊：不会再有新任务

	openCount := 0
	for i := 0; i < maxPort; i++ {
		if r := <-results; r.Open {
			openCount++
			fmt.Printf("端口 %d 开放\n", r.Port)
		}
	}
	fmt.Printf("扫描 1-%d：开放 %d 个，总耗时 %v\n", maxPort, openCount, time.Since(start))
}
