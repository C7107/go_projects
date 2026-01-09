/*
练习 3：幸运大转盘 (Select + TimeOut)
考察点： select 的多路复用和超时处理。
题目描述：
模拟一个网络请求。
写一个函数 server()，它会 time.Sleep 一个随机时间（0 到 2 秒之间），然后往一个 Channel 发送 "数据到了"。
在 main 函数里，使用 select 监听这个 Channel。
规则：
如果在 1 秒内收到了数据，打印 "成功收到数据"。
如果超过 1 秒还没收到，打印 "超时了，不玩了" 并退出。
提示： 超时可以用 case <-time.After(1 * time.Second):。
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. 在 main 里创建 channel
	ch := make(chan string)
	go server(ch)
	select {
	case data := <-ch:
		fmt.Printf("收到收到，嘿嘿嘿，数据是%s", data)
	case <-time.After(2 * time.Second):
		fmt.Println("密码的，不玩了，这么久都没来")
	}

}

// chan<- string 表示在这个函数里，我只负责往里“写”数据（单向通道，更安全）
func server(c chan<- string) {
	time.Sleep(3 * time.Second)
	c <- "啊啊啊豪爽"
}
