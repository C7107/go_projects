package main

import "fmt"

// 1. 先定义类型
type MyInt int

// 2. 定义方法
// 接收器 m 代表"加号"左边的数
// 参数 b 代表"加号"右边的数
// 返回值类型最好也保持一致用 MyInt
func (m MyInt) Add(b MyInt) MyInt {
	return m + b
}

func main() {
	m := MyInt(1) // 第一个数
	b := MyInt(2) // 第二个数

	// 像人类说话一样："m 加上 b"
	result := m.Add(b)

	fmt.Println(result) // 输出 3
}
