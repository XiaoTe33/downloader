package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i <= 100; i++ {
		mk("test", 100, i)
		time.Sleep(10 * time.Millisecond)
	}
}

func mk(name string, all, now int) {
	//清除一整行内容
	fmt.Print("\r\033[K")

	//制作一个空框
	fmt.Print(name + "[" + space(101) + "]")
	//计算进度百分比
	n := now * 100 / all
	fmt.Printf("total:%d%%", n)

	//制作一定长度的条
	s := ""
	for i := 0; i < n; i++ {
		s += "="
	}
	fmt.Printf("\r\033[%dC%s>", len(name)+1, s)
	//隐藏光标，形成纯享版
	fmt.Print("\033[?25l")

	//下载完显示
	//name Done
	if all == now {
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\r\033[K")
		fmt.Print(name + "\033[32m Done\033[0m\n")
		return
	}
}

func space(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}
