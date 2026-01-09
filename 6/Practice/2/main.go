/*
第二题：值接收器 vs 指针接收器
考点： 什么时候需要加星星 *（修改原值）。
题目描述：
我们需要一个计数器 Counter，它有一个字段 Value。我们希望调用 Inc() 方法时，它的值能加 1。
请看下面的代码，目前的写法有一个严重 Bug（加不上去）。
有 Bug 的代码：
*/

package main

import (
	"fmt"
)

type Counter struct {
	Value int
}

// ❌ 这里的接收器有问题，导致修改无效
func (c *Counter) Inc() {
	c.Value++
}

func main() {
	c := Counter{Value: 0}
	(&c).Inc()
	fmt.Println(c.Value) // 预期输出 1，实际输出了 0，为什么？
}
