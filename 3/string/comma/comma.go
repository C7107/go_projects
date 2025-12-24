package comma

import (
	"bytes"
	"strings"
)

func Comma(s string) string {
	var buf bytes.Buffer

	// 1. 处理空字符串边界情况
	if len(s) == 0 {
		return ""
	}

	// 2. 处理正负号
	// start 用于记录数字真正开始的位置（跳过符号位）
	start := 0
	if s[0] == '+' || s[0] == '-' {
		buf.WriteByte(s[0]) // 先把符号写入 buffer
		start = 1           // 标记数字从第1位开始
	}

	// 截取不带符号的纯数字部分
	// 例如 "-12345.67" -> "12345.67"
	numberPart := s[start:]

	// 3. 分离整数部分和小数部分
	dot := strings.Index(numberPart, ".")
	var integerPart, fractionalPart string

	if dot >= 0 {
		// 如果有小数点
		integerPart = numberPart[:dot]    // 小数点左边
		fractionalPart = numberPart[dot:] // 小数点右边（包含小数点）
	} else {
		// 如果没有小数点
		integerPart = numberPart
		fractionalPart = ""
	}

	// 4. 处理整数部分（添加逗号的核心逻辑）
	n := len(integerPart)
	if n <= 3 {
		// 如果整数部分小于等于3位，直接写入
		buf.WriteString(integerPart)
	} else {
		// 计算头部长度
		remainder := n % 3
		if remainder == 0 {
			remainder = 3
		}

		// 写入头部
		buf.WriteString(integerPart[:remainder])

		// 循环写入剩余部分
		for i := remainder; i < n; i += 3 {
			buf.WriteByte(',')
			buf.WriteString(integerPart[i : i+3])
		}
	}

	// 5. 把小数部分拼接到最后
	buf.WriteString(fractionalPart)

	return buf.String()
}
