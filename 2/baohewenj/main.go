// 文件路径: E:\go_projects\2\main.go
package main

import (
	"fmt"
	// 路径公式：模块名(go.mod里写的) + 文件夹路径
	"github.com/C7107/go_projects/2/baohewenj/tempconv"
)

func main() {
	fmt.Println("正在运行第二章的示例...")
	c := tempconv.BoilingC
	f := tempconv.CToF(c)
	fmt.Printf("沸点: %v 转换为 %v\n", c, f)
}
