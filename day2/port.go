package main

import "fmt"

type PortResult struct {
	Port int
	Open bool
}

// PortResult 实现了 Describer 接口：
// Go 的接口是隐式实现的——只要有 Describe() 方法，不用声明 implements。
func (r PortResult) Describe() string {
	status := "关闭"
	if r.Open {
		status = "开放"
	}
	return fmt.Sprintf("端口 %d %s", r.Port, status)
}

// PrintAll 接受接口切片：任何实现了 Describer 的类型都能放进来，
// 统一调用 Describe()——这就是“参数用接口”的口诀。
func PrintAll(ds []Describer) {
	for _, d := range ds {
		fmt.Println(d.Describe())
	}
}
