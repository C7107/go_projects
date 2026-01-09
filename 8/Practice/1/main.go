/*
练习 1：大厨的点名 (Goroutine + sync.WaitGroup)
考察点： 既然“主死全死”，如何让主程序（大厨）优雅地等待所有小工干完活？
题目描述：
编写一个程序，启动 3 个 Goroutine（小工）。
小工 A 打印 "切菜"
小工 B 打印 "烧水"
小工 C 打印 "摆盘"
要求： 主程序必须等待这 3 件事都打印出来后，才能打印 "上菜！" 并退出。
提示： 不要用 time.Sleep 去猜时间，要用你笔记里提到的 sync.WaitGroup。
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3) //必须写在外面，因为当你执行 go func()... 时，主程序（Main）不会等待这个协程启动，而是立刻往下执行 wg.Wait()
	//主程序执行 go func...（小工A被创建，但还没来得及跑）。
	//主程序继续极速往下走，直接运行到了 wg.Wait()。
	//此时，小工A还没来得及运行第一行代码 wg.Add(1)。
	//Wait 发现计数器目前还是 0！（以为没人干活）
	//Wait 直接放行，主程序打印 "上菜"，然后退出。
	go func() {
		fmt.Println("切菜")
		defer wg.Done()
	}()
	go func() {
		fmt.Println("烧水")
		defer wg.Done()
	}()
	go func() {
		fmt.Println("摆盘")
		defer wg.Done()
	}()
	wg.Wait()
	fmt.Println("上菜")
}
