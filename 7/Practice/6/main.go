/*
进阶挑战（选做）：排序整理师
知识点：sort.Interface。
题目描述：
定义一个 []int 类型的别名，比如 type IntSlice []int。
为它实现 Len, Less, Swap 三个方法。
注意：在 Less 方法里，我们这次搞个特殊的，从大到小排序（即 p[i] > p[j]）。
在 main 里定义一个切片 nums := IntSlice{3, 1, 4, 2}，调用 sort.Sort(nums)，看结果是不是变成了 [4 3 2 1]。
*/
package main

import (
	"fmt"
	"sort"
)

type IntSlice []int

func (p IntSlice) Len() int {
	return len(p)
}

func (p IntSlice) Less(i, j int) bool {
	return p[i] > p[j]
}

func (p IntSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	nume := IntSlice{3, 1, 4, 2}
	fmt.Printf("原始的：%d\n", nume)
	sort.Sort(nume)
	fmt.Printf("排序后：%d", nume)
}
