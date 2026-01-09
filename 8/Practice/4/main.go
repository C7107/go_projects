/*
练习 4：流水线上的数字 (Pipeline + Close)
考察点： Channel 的串联、关闭 Channel 的重要性、range 循环。
题目描述：
模仿笔记中的“流水线”例子，建立一条有 3 个阶段的流水线：
阶段 1（生成者）： 往管道 A 发送数字 1, 2, 3, 4, 5，发完后关闭管道 A。
阶段 2（处理器）： 从管道 A 读取数字，把它们乘以 10（变成 10, 20...），发送到管道 B。等管道 A 的数据处理完后，关闭管道 B。
阶段 3（打印者/主程序）： 从管道 B 读取结果并打印。
关键点： 确保每个阶段都知道上一阶段什么时候结束（通过 range 和 close）。这个是重点
*/
package main

import "fmt"

func main() {
	Ach := make(chan int)
	Bch := make(chan int)
	go Producer((Ach))
	go Handle(Ach, Bch)

	for data := range Bch {
		fmt.Println(data)
	}
}

func Multiplication(a int) int {
	return a * 10
}

func Producer(a chan<- int) {
	for i := 1; i < 6; i++ {
		a <- i
	}
	defer close(a)
}

func Handle(a <-chan int, b chan<- int) {
	for data := range a {
		b <- Multiplication(data)
	}
	defer close(b)
}
