# go-learn：一个安全工程师的 Go 学习记录

对照 Python 讲、每天一个可运行的例子，从零写到并发端口扫描器。

## 环境

- Go 1.21+（开发时用 1.26.5，`go version` 查看你的版本）
- 每个目录是独立 module，进目录 `go run .` 即可运行

## 学习路线

| 目录 | 主题 | 关键内容 |
|---|---|---|
| `day1` | 基础语法 | 变量、函数、多返回值、slice、map、byte/rune、`(结果, error)` |
| `day2` | struct 与接口 | struct、方法、值/指针接收者、interface 隐式实现、自定义 error |
| `day3` | 并发入门 | goroutine、channel（缓冲/无缓冲）、WaitGroup、select |
| `select-timeout` | select 超时模式 | `select` + `time.After`：Go 里“不等了”的标准写法 |
| `day4` | 实战：并发端口扫描器 | `net.DialTimeout` 真连端口 + worker 池 |

```bash
cd day4 && go run .   # 100 个 worker 并发扫本机 1-1024 端口
```

## 安全场景

Day 4 的端口扫描器是安全研发真实在用的小工具雏形：资产测绘、指纹识别、
漏洞扫描里到处都是“worker 池 + 任务队列”这个模式。

## 学习方法

AI 加速 + 理解验收：每学一小块就亲手跑通，验收标准是“能讲清、能找 bug、
能做对变式题”，不背定义。代码里的注释记录了学习过程中真实踩过的坑。
