package main

import "fmt"

// 写法二：强制至少传一个参数 (更安全，推荐)
func max(first int, rest ...int) int {
	m := first
	for _, v := range rest {
		if v > m {
			m = v
		}
	}
	return m
}

// 同理 min
func min(first int, rest ...int) int {
	m := first
	for _, v := range rest {
		if v < m {
			m = v
		}
	}
	return m
}

func main() {
	// 1. 正常传多个参数
	fmt.Println(max(10, 5, 88, 3)) // 输出: 88
	fmt.Println(min(10, 5, 88, 3)) // 输出: 3

	// 2. 只传一个参数 (完全合法，rest 为空)
	fmt.Println(max(100)) // 输出: 100

	// 3. 传切片 (需要解包)
	nums := []int{1, 2, 3, 4}
	// 注意：这里要把切片拆开
	// 第一个参数是 nums[0]，剩下的是 nums[1:]...
	// 稍微有点麻烦，但保证了类型安全
	if len(nums) > 0 {
		fmt.Println(max(nums[0], nums[1:]...))
	}

	// 4. 不传参数？
	// fmt.Println(max()) // 编译报错！not enough arguments
	// 这就是在编译阶段帮我们发现了错误
}
