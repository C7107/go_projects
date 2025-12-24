package main

import (
	"fmt"

	"github.com/C7107/go_projects/3/string/comma"
)

func main() {
	examples := []string{
		"12345",       // 普通整数
		"-12345",      // 负整数
		"+1234567",    // 正整数
		"12345.6789",  // 浮点数
		"-1234567.89", // 负浮点数
		".12345",      // 只有小数部分
		"12",          // 短整数
		"",            // 空字符串
	}

	for _, s := range examples {
		fmt.Printf("原始: %-12s -> 处理后: %s\n", s, comma.Comma(s))
	}
}
