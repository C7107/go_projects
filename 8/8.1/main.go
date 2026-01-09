package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. 启动小工：去后台转圈圈（显示动画）
	// 这里的 go 关键字是关键！
	// spinner 函数会一直在那里循环打印 - \ | /，像风车一样转
	go spinner(100 * time.Millisecond)

	const n = 100
	// 2. 大厨（主程序）亲自计算斐波那契数列
	// fib(45) 使用递归计算非常慢，可能要几秒钟到十多秒。
	// 在这期间，主程序被“卡”在这里计算。
	// 但是！刚才启动的小工（spinner）还在旁边干活，所以屏幕上你会看到动画在转。
	fibN := fib(n) // slow

	// 3. 计算完了，大厨打印结果
	// \r 会把光标移回到行首，覆盖掉刚才转圈圈的动画
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
} // 4. main 函数结束，spinner 也会随之立刻停止

// 动画函数：负责显示转动的字符
func spinner(delay time.Duration) {
	for {
		// 循环打印 - \ | /
		for _, r := range `-\|/` {
			// \r 表示回到行首，%c 打印字符
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// 斐波那契数列计算函数：使用低效的递归算法，为了模拟耗时操作
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
