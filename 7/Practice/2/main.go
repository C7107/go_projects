/*
第二关：自我介绍（fmt.Stringer）
知识点：常用标准库接口 fmt.Stringer。
题目描述：
定义一个结构体 Person，包含 Name (string) 和 Age (int)。
如果直接 fmt.Println(p)，它会打印 {Name:..., Age:...}，太丑了。
请为 Person 实现 String() string 方法，使得打印时输出格式为："我是 [Name]，今年 [Age] 岁"。
在 main 函数中创建一个 Person 并打印它。
*/
package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (s Person) String() string {
	return fmt.Sprintf("我是%s，今年%d岁", s.Name, s.Age)
}

func main() {
	p := Person{Name: "杨格", Age: 18}
	fmt.Println(p)
}
