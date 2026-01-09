/*
练习 5：限制并发的爬虫 (带缓存 Channel)
考察点： 利用带缓存 Channel 实现“信号量”（夜店保安模式），防止资源耗尽。
题目描述：
假设你有 100 个任务（打印数字 0 到 99），每个任务耗时 1 秒（用 sleep 模拟）。
如果你直接 go func 启动 100 个，它们会同时开始打印（瞬间刷屏）。
要求： 修改程序，限制同一时间只能有 5 个任务在运行。
效果： 你应该能看到屏幕上每 5 个一组地蹦出数字，而不是一下子全部蹦出来。
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var tokens = make(chan struct{}, 10)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		tokens <- struct{}{}
		wg.Add(1)
		go p(i, tokens, &wg)
	}
	wg.Wait()
}

func p(a int, ch <-chan struct{}, wg *sync.WaitGroup) { //wg必须要传指针
	fmt.Println(a)
	time.Sleep(1 * time.Second)
	<-ch
	wg.Done()
}
