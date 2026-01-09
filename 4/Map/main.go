// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func plan1() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func plan2() {
	// 创建一个 map 用于统计不同类别的数量
	// Key 是类别名称（如 "Letter", "Digit"），Value 是计数
	categories := make(map[string]int)

	invalid := 0 // 统计无效 UTF-8 字符

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		// 处理非法 UTF-8 编码 (即前面解释的 invalid case)
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		// 使用 switch 分类统计
		switch {
		case unicode.IsLetter(r):
			categories["Letter"]++
		case unicode.IsDigit(r):
			categories["Digit"]++
		case unicode.IsSpace(r):
			categories["Space"]++ // 包含空格、换行、制表符
		case unicode.IsPunct(r):
			categories["Punct"]++ // 标点符号
		case unicode.IsSymbol(r):
			categories["Symbol"]++ // 符号 (如 $, +)
		default:
			categories["Other"]++
		}
	}

	// 打印统计结果
	fmt.Printf("category\tcount\n")
	for cat, n := range categories {
		fmt.Printf("%s\t%d\n", cat, n)
	}

	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func plan3() {
	// 创建一个 map 用于存储单词及其频率
	// Key: 单词 (string), Value: 出现频率 (int)
	counts := make(map[string]int)

	// 创建一个新的 Scanner 来读取标准输入
	input := bufio.NewScanner(os.Stdin)

	// **** 关键一步：设置 Scanner 的分割函数为按单词分割 ****
	// 默认是 bufio.ScanLines (按行分割)
	input.Split(bufio.ScanWords)

	// 循环读取每一个单词
	for input.Scan() {
		word := input.Text() // 获取当前扫描到的单词

		// 改进：将单词转换为小写，以忽略大小写差异
		// 例如："The" 和 "the" 都会被统计为 "the"
		word = strings.ToLower(word)

		// 增加该单词的计数
		counts[word]++
	}

	// 检查是否有错误发生 (除了 io.EOF，EOF不会作为错误返回给Scanner)
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "wordfreq: %v\n", err)
		os.Exit(1)
	}

	// 打印结果
	fmt.Printf("word\tcount\n")
	// Map 的遍历顺序是随机的
	for w, c := range counts {
		fmt.Printf("%s\t%d\n", w, c)
	}
}

func main() {

}
