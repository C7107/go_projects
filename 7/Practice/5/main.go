/*
第五关：智能分拣机（类型分支）
知识点：Type Switch switch x.(type)。
题目描述：
写一个函数 CheckType(v interface{})。
在函数内部使用 switch v.(type) 进行判断：
如果是 int，打印 "这是一个整数"。
如果是 string，打印 "这是一个字符串"。
如果是 bool，打印 "这是一个布尔值"。
其他情况，打印 "未知类型"。
在 main 函数里分别传几个不同类型的值进去测试。
*/
package main

import "fmt"

func CheckType(v interface{}) {
	switch v := v.(type) {
	case nil:
		fmt.Println("是nil")
	case int:
		fmt.Println("是int")
	case string:
		fmt.Println("是string")
	case bool:
		fmt.Println("是bool")
	default:
		fmt.Println("啥也不是")
		_ = v // 明确告诉编译器：我知道你存在，但我不用，可以消除v未被使用的错误
	}
}

func main() {
	CheckType(1)
	CheckType("yg")
	CheckType(true)
}
