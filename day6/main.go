package main

// Day 6：子域名爆破 —— 字典枚举 + DNS 解析
//
// 本机环境限制（2026-10-03 实测）：Clash TUN 劫持全部 53 端口流量，
// UDP/TCP 直连 223.5.5.5 均被转回 fake-ip DNS（198.18.0.0/16）。
// 因此本机只能验证“零误报”；isFake 过滤器兜底，换正常网络可出真实命中。

import (
	"context"
	"fmt"
	"net"
	"time"
)

var fakeNet *net.IPNet

func init() {
	// Clash fake-ip网段:这个网段里的解析结果都是假的,直接丢弃
	_, fakeNet, _ = net.ParseCIDR("198.18.0.0/16")
}

// 只要有一个IP 不是假IP,就当它是真的
func isFake(ips []string) bool {
	for _, s := range ips {
		if ip := net.ParseIP(s); ip == nil || !fakeNet.Contains(ip) {
			return false
		}
	}
	return true
}

func main() {
	resolver := &net.Resolver{
		PreferGo: true, //用Go 原生解析器,不调系统的,这样下面Dial才会被用上
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "tcp", "223.5.5.5:53")
		},
	}

	domain := "baidu.com"
	subs := []string{"notexist12345", "www", "mail", "api", "dev", "test", "blog", "ftp", "vpn"}
	for _, sub := range subs {
		host := sub + "." + domain
		ips, err := resolver.LookupHost(context.Background(), host)
		if err != nil {
			continue //解析失败=不存在,跳过
		}
		if isFake(ips) {
			continue //假 IP,丢弃
		}
		fmt.Println(host, "->", ips)
	}
	fmt.Println("done")
}
