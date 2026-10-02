package main

// Day 1：Go 基础语法 —— 变量、函数、多返回值、slice、map、byte/rune、(结果, error)
// 5 个练习：fib / dedup / wordCount / reverse / max
// 运行：cd day1 && go run .

import (
	"errors"
	"fmt"
	"strings"
)

// 1. fib: 返回斐波那契数列前 n 个
func fib(n int) []int {
	if n <= 0 {
		return []int{}
	}
	result := make([]int, 0, n) // 预分配容量，避免反复扩容
	for i := 0; i < n; i++ {
		if i < 2 {
			result = append(result, i)
		} else {
			result = append(result, result[i-1]+result[i-2])
		}
	}
	return result // append 的返回值必须接住
}

// 2. dedup: slice 去重，保留首次出现顺序，O(n)
// 用 map[int]bool 当 set：bool 零值 false 天然表示"没见过"
// 另注：O(n^2) 原地双重循环版删除元素后必须 j--，
// 否则左移顶上来的元素会被 j++ 跳过（如 dedup([]int{1,1,1}) 会漏删）
func dedup(s []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// 3. wordCount: 词频统计
// 查不存在的 key 返回零值 0，所以 counts[w]++ 一行搞定，不用先判存在
func wordCount(s string) map[string]int {
	words := strings.Fields(s) // 按空白切分
	counts := make(map[string]int)
	for _, w := range words {
		counts[w]++
	}
	return counts // map 打印顺序随机，属正常
}

// 4. reverse: 字符串反转
// 必须用 []rune（字符）而非 []byte（字节）：中文 1 字占 3 字节，
// 按字节交换会把字拆碎乱码。英文/二进制场景才用 byte。
func reverse(s string) string {
	r := []rune(s)
	n := len(r)
	if n < 2 {
		return s
	}
	for i := 0; i < n/2; i++ {
		r[i], r[n-1-i] = r[n-1-i], r[i] // Go 多重赋值，一行交换
	}
	return string(r)
}

// 5. max: 返回 slice 最大值，空 slice 返回 error
// Go 惯例：(结果, error)，调用方用 if err != nil 处理
func max(s []int) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("slice is empty") // 静态字符串不用 fmt.Errorf
	}
	maxVal := s[0] // 变量名别叫 max，会遮蔽函数名
	for _, v := range s[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal, nil
}

func main() {
	fmt.Println("fib(10):", fib(10))
	fmt.Println("fib(1):", fib(1))
	fmt.Println("fib(0):", fib(0))

	fmt.Println("dedup:", dedup([]int{1, 2, 2, 3, 1, 4, 5}))

	fmt.Println("wordCount:", wordCount("hello world hello"))

	fmt.Println("reverse(hello):", reverse("hello"))
	fmt.Println("reverse(你好世界):", reverse("你好世界"))

	v, err := max([]int{3, 1, 4})
	fmt.Println("max:", v, err)
	v, err = max([]int{})
	fmt.Println("max(empty):", v, err)
}
