package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// 1. 创建解码器，连接到标准输入（os.Stdin）
	// 就像接通了传送带电源
	dec := xml.NewDecoder(os.Stdin)

	// 2. 定义一个切片当作“栈”，用来存当前进去的标签名
	var stack []string

	// 3. 死循环，不停地从传送带拿零件
	for {
		// dec.Token() 拿下一个零件
		tok, err := dec.Token()

		// 如果传送带空了（文件读完了），就退出
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "出错啦: %v\n", err)
			os.Exit(1)
		}

		// 4. 【关键】类型分支（Type Switch）
		// 看看拿到的是什么零件
		switch tok := tok.(type) {

		// 情况 A: 是开始标签 (如 <div>)
		case xml.StartElement:
			// 把标签名压入栈
			stack = append(stack, tok.Name.Local)

		// 情况 B: 是结束标签 (如 </div>)
		case xml.EndElement:
			// 把栈顶的元素弹出去（切片长度减1）
			stack = stack[:len(stack)-1]

		// 情况 C: 是文本内容 (如 "Hello")
		case xml.CharData:
			// 5. 检查当前栈里的路径，是否包含用户想要的路径
			// os.Args[1:] 是用户输入的参数，比如 ["div", "h2"]
			if containsAll(stack, os.Args[1:]) {
				// 如果匹配，把栈拼成字符串，打印出文本内容
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}
		}
	}
}

// 辅助函数：检查 stack 是否按顺序包含了 target 中的所有元素
// 比如 stack=["html", "body", "div", "h2"]
// target=["div", "h2"]
// 结果就是 true
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}
