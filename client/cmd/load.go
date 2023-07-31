package main

import (
	"fmt"
	"time"
)

func NewLine(filename string, total int64, buf int, begin time.Time) chan<- struct{} {
	ch := make(chan struct{}, 10)
	now := 0
	go func(Ch chan struct{}) {
		for {
			select {
			case <-Ch:
				now++
				Line(filename, total, int64(now*buf), begin)
			}
		}
	}(ch)
	return ch
}

func Line(name string, all, now int64, begin time.Time) {
	//清除一整行内容
	fmt.Print("\r\033[K")

	//制作一个空框
	fmt.Print(name + "[" + space(101) + "]")
	//计算进度百分比
	n := now * 100 / all
	fmt.Printf("total:%d%%", n)

	duration := time.Now().UnixMilli() - begin.UnixMilli()
	fmt.Printf(" speed:%.2fM/s", (float64(now-all)/float64((duration)/1000.0))/1024.0)
	//制作一定长度的条
	s := ""
	for i := 0; i < int(n); i++ {
		s += "="
	}
	fmt.Printf("\r\033[%dC%s>", len(name)+1, s)
	//隐藏光标，形成纯享版
	fmt.Print("\033[?25l")

	//下载完显示
	//name Done
	if all <= now {
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
