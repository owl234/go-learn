package main

import "fmt"

// type 名字 struct { 字段名 类型 }
type Scanner struct {
	Target  string
	Ports   []int
	Timeout int
}

// 方法 = 带接收者的函数, (s Scanner) 就是Python 的self
// Go 习惯用首字母小写单字母做接收者,不叫self

func (s *Scanner) Describe() string {
	return fmt.Sprintf("扫描 %s 的 %d 个端口, 超时 %d 秒", s.Target, len(s.Ports), s.Timeout)
}

func (s *Scanner) AddPort(p int) {
	s.Ports = append(s.Ports, p)
}

// 只要有个Describe() 方法,就是Describer
type Describer interface {
	Describe() string
}

func PrintDescription(d Describer) {
	fmt.Println(d.Describe())
}

type ScanError struct {
	Port int
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("端口 %d 扫描失败", e.Port)
}

func CheckPort(p int) error {
	if p < 0 {
		return &ScanError{Port: p}
	}
	return nil
}

func main() {
	s := Scanner{ //复合字面量初始化,类似dict传参
		Target:  "127.0.0.1",
		Ports:   []int{80, 443},
		Timeout: 5,
	}
	fmt.Println(s)
	fmt.Println(s.Target)

	fmt.Println(s.Describe())
	s.AddPort(8080)
	fmt.Println(s.Ports)

	p := s // 值拷贝：p 是 s 的完整复制，两者是独立变量
	fmt.Printf("%T\n", p)  // %T 打印变量的类型
	p.Timeout = 10         // 改的是拷贝 p
	fmt.Println(s.Timeout) // 仍然是 5：s 不受影响（对比指针接收者的 AddPort）

	PrintDescription(&s)

	if err := CheckPort(-1); err != nil {
		fmt.Println(err)
	}

	items := []Describer{
		&Scanner{Target: "127.0.0.1", Ports: []int{80, 443}, Timeout: 5},
		PortResult{Port: 80, Open: true},
		PortResult{Port: 443, Open: false},
	}
	PrintAll(items)
}
