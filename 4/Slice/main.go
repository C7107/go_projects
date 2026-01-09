package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 练习 4.3： 重写reverse函数，使用数组指针代替slice。
func reverse1(s *[6]int) { // 注意：这里必须写死长度 6，否则类型不匹配
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 练习 4.5： 写一个函数在原地完成消除[]string中相邻重复的字符串的操作
func EliminateDuplicatesString(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	out := ss[:0]
	for _, s := range ss {
		// 如果 out 是空的（放入第一个元素）
		// 或者 当前元素 s 不等于 out 里最后一个元素
		if len(out) == 0 || out[len(out)-1] != s {
			out = append(out, s)
		}
	}

	// 4. 返回新的切片（长度已经变短了）
	return out
}

func squashSpaces(b []byte) []byte {
	w := 0             // 写指针（byte 索引）
	prevSpace := false // 上一个 rune 是否是空白

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)

		if unicode.IsSpace(r) {
			if !prevSpace {
				b[w] = ' ' // 统一替换成 ASCII 空格
				w++
				prevSpace = true
			}
		} else {
			// 把当前 rune 原样写回
			w += utf8.EncodeRune(b[w:], r)
			prevSpace = false
		}

		b = b[size:]
	}

	return b[:w]
}

func plan1() {
	a := [6]int{0, 1, 2, 3, 4, 5}
	fmt.Println("反转前:", a)
	// 关键点：传入数组的地址 (&a)
	reverse1(&a)
	fmt.Println("反转后:", a)
}

func plan2() {
	data := []byte("Hello\t\t世界   \n\nGo 语言")
	result := squashSpaces(data)
	fmt.Printf("%q\n", string(result))
}

func main() {
	plan2()
}
