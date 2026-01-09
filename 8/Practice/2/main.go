/*
练习 2：乒乓球接力 (无缓存 Channel)
考察点： 无缓存 Channel 的“同步阻塞”特性（一手交钱一手交货）。
题目描述：
创建两个 Goroutine，一个叫“发球者”，一个叫“接球者”。创建一个无缓存的 Channel 用来传球（int 类型）。
发球者：把数字 1 发送到 Channel 里，然后打印 "球发出了"。
接球者：从 Channel 里拿到数字，打印 "球接到了: 1"。
主程序：等待这两个人都做完（可以用 WaitGroup 或 Channel 通知）。
思考： 如果你把代码里的 channel <- 1 和 <-channel 写在同一个 Goroutine 里（比如都在 main 函数里），
会发生什么？（回顾笔记中 Netcat2 的部分）
答思考：会死锁，pass <- 123 就像你要把接力棒交出去。它有一个强制要求：必须有人在对面伸手接，你才能松手，
导致主程序卡在了上一行，它永远不可能执行到下一行去接球
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	pass := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		pass <- 123
		fmt.Println("发球爽爽爽")
		defer wg.Done()
	}()
	go func() {
		i := <-pass
		fmt.Printf("接球爽爽爽，这个球是%d", i)
		defer wg.Done()
	}()
	wg.Wait()
	fmt.Println("比赛结束")
}

func plan2() { //等待这两个人都做完 Channel 通知
	pass := make(chan int)
	// 1. 创建一个用来通知完成的 channel
	// struct{}{} 是空结构体，不占内存，专门用来当信号（就像拍一下肩膀）
	// 也可以用 chan bool 或 chan int
	done := make(chan struct{})

	go func() {
		pass <- 123
		fmt.Println("发球爽爽爽")

		// 2. 干完活，往 done 里塞一个信号
		done <- struct{}{}
	}()

	go func() {
		i := <-pass
		fmt.Printf("接球爽爽爽，这个球是%d\n", i)

		// 2. 干完活，也往 done 里塞一个信号
		done <- struct{}{}
	}()

	// 3. 主程序在此阻塞等待
	// 因为启动了2个协程，所以要收2次信号
	<-done // 收第一个（不知道是谁先完，反正收一个）
	<-done // 收第二个

	fmt.Println("比赛结束")
}
