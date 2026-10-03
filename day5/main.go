package main

// Day 5：实战 —— 端口扫描 + HTTP 指纹识别
//
// 在 Day 4 的 worker 池扫描器上加一层：端口开放后，用 HTTP 抓 Server 头识别服务。
// 指纹函数 grabServer 吃端口号、吐 Server 头；非 HTTP 服务或抓不到时返回空字符串。
// 注意 http.Client 必须设 Timeout（http.Get 默认无超时，会卡住 worker）。
//
// 运行：cd day5 && go run .（先起一个 python3 -m http.server 8080 对照看指纹）

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type PortResult struct {
	Port int
	Open bool
	Server string
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
		server := grabServer(port)
		results <- PortResult{Port: port, Open: true, Server: server}
	}

}

// 指纹函数：吃端口号，吐 Server 头；抓不到就吐空字符串
func grabServer(port int) string {
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := http.Client{Timeout: 2 * time.Second} // http.Get 默认无超时，必须显式设置
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close() // body 是网络连接，不关会泄漏
	return resp.Header.Get("Server")
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
			if r.Server != "" {
				fmt.Printf("端口 %d 开放 [%s]\n", r.Port, r.Server)
			} else {
				fmt.Printf("端口 %d 开放\n", r.Port)
			}
			
		}
	}
	fmt.Printf("扫描 1-%d：开放 %d 个，总耗时 %v\n", maxPort, openCount, time.Since(start))
}
