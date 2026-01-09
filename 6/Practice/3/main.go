/*
第三题：神奇的结构体嵌入
考点： 嵌入结构体后的方法调用（语法糖）。
题目描述：
我们有一个 Engine（引擎）结构体，它有一个 Start() 方法。
我们要造一辆 Car（车），把 Engine 嵌入进去。
*/
package main

import "fmt"

type Engine struct{}

func (e Engine) Start() {
	fmt.Println("引擎启动：轰轰轰！")
}

type Car struct {
	Engine // 嵌入字段，没有名字
}

func main() {
	c := Car{}
	c.Start()
	// 你的任务：在这里写一行代码，让车启动。
	// 要求：不要显式地写出 .Engine 中间层。

}
