/*
第四关：揭开面具（类型断言）
知识点：安全类型断言 value, ok := x.(T)。
题目描述：
定义一个空接口变量 var x interface{}，并赋值为字符串 "Go语言"。
任务A：尝试断言它是 int 类型。请使用“安全模式”（带 ok 的写法），如果失败，打印 "断言失败，它不是整数"。
任务B：尝试断言它是 string 类型。如果成功，打印 "断言成功，内容是：..."。
*/

package main

import "fmt"

func main() {
	var x interface{}
	x = "Go语言"
	if _, ok := x.(int); ok {
		fmt.Println("哦对的对的,他就是int")
	} else {
		fmt.Println("欧不对不对，他不是int")
	}
	if _, ok := x.(string); ok {
		fmt.Printf("哦对的对的,他就是string,他叫:%s", x)
	} else {
		fmt.Println("欧不对不对，他不是string")
	}
}
