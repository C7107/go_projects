/*
第四题：nil 会报错吗？
考点： nil 接收器的安全性。
题目描述：
在 Java 或 Python 中，空对象调用方法通常会直接报错（空指针异常）。但在 Go 里不一定。
请阅读下面的代码，预测运行结果。
代码：
*/
package main

import "fmt"

type Number struct {
	Val int
}

// 方法定义
func (n *Number) SayHello() {
	if n == nil {
		fmt.Println("我是空的，但我没崩！")
		return
	}
	fmt.Println("我的值是", n.Val)
}

func main() {
	var p *Number = nil // 定义一个空指针
	p.SayHello()        // 这里会发生什么？
}
