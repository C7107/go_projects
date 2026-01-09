package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

func popCount(b byte) int {
	count := 0
	for b != 0 {
		b = b & (b - 1)
		count++
	}
	return count
}

// SHA256BitDiff 计算两个 SHA256 哈希码之间不同 bit 的数目
func SHA256BitDiff(c1, c2 [32]byte) int {
	totalDiff := 0
	// SHA256 结果是固定的 32 个字节
	for i := 0; i < 32; i++ {
		// 1. 先做异或运算：结果中为 1 的位表示原数据在该位不同
		xorVal := c1[i] ^ c2[i]

		// 2. 统计这个字节里有多少个 1，累加到总数
		totalDiff += popCount(xorVal)
	}
	return totalDiff
}

func plan1() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	fmt.Printf("c1: %x\n", c1)
	fmt.Printf("c2: %x\n", c2)

	diff := SHA256BitDiff(c1, c2)
	fmt.Printf("不同 bit 的数目: %d\n", diff)
}

func StringsToSHA384(ss []string) string {
	// 方法1：用空字符串连接（也可用其他分隔符，如 "\n"、"," 等）
	data := strings.Join(ss, "")

	// 创建 SHA-384 哈希器（注意：使用 sha512.New384()）
	hasher := sha512.New384()
	hasher.Write([]byte(data))

	// 返回十六进制字符串
	return hex.EncodeToString(hasher.Sum(nil))

}

func StringsToSHA512(ss []string) string {
	// 方法1：用空字符串连接（也可用其他分隔符，如 "\n"、"," 等）
	data := strings.Join(ss, "")

	// 创建 SHA-384 哈希器（注意：使用 sha512.New384()）
	hasher := sha512.New()
	hasher.Write([]byte(data))

	// 返回十六进制字符串
	return hex.EncodeToString(hasher.Sum(nil))

}

func StringsToSHA256(ss []string) string {
	data := strings.Join(ss, "")

	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func plan2() {
	osha384 := flag.Bool("sha384", false, "Output sha384")
	osha512 := flag.Bool("sha512", false, "Output sha512")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("输入格式错误")
		os.Exit(1)
	}
	// ❗ 互斥选择：只能选一种算法
	if *osha384 && *osha512 {
		fmt.Println("不能同时指定 -sha384 和 -sha512")
		os.Exit(1)
	}
	if *osha384 {
		fmt.Print(StringsToSHA384(flag.Args()))
	} else if *osha512 {
		fmt.Print(StringsToSHA512(flag.Args()))
	} else {
		// 默认：SHA256
		fmt.Print(StringsToSHA256(flag.Args()))
	}
}

func main() {
	plan2()
}
