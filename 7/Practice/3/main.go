/*
第三关：万能盒子与空接口
知识点：interface{}、切片。
题目描述：
我们知道 interface{} 可以装任何东西。
请定义一个切片，这个切片可以同时存放：一个整数 100，一个字符串 "Hello"，和一个布尔值 true。
使用 for-range 循环遍历这个切片，简单打印出里面的内容。
*/

package main

import "fmt"

type w interface{}

func main() {
	s := []w{100, "Hello", true}
	for _, value := range s {
		fmt.Println(value)
	}
}

//可以直接写 []interface{}，不用先定义了一个类型别名 type w interface{}。
//从 Go 1.18 版本开始，官方觉得每次写 interface{} 太麻烦了，所以内置了一个别名叫做 any
//可以直接写 s := []any{100, "Hello", true}
