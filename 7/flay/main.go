package main

import (
	"flag"
	"fmt"
)

// ---------------------------------------------------------
// 1. 定义基础类型和转换逻辑
// ---------------------------------------------------------

// 定义摄氏度和华氏度
type Celsius float64
type Fahrenheit float64

// 华氏度转摄氏度的公式
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// 给 Celsius 定义 String 方法
// 这样它就实现了 fmt.Stringer 接口，打印时会自动显示 "xx°C"
func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}

// ---------------------------------------------------------
// 2. 定义 flag.Value 接口的实现
// ---------------------------------------------------------

// 定义一个结构体来“包装” Celsius
// 为什么要包装？因为我们需要给它加上 Set 方法，让它满足 flag.Value 接口
type celsiusFlag struct {
	Celsius
}

// 实现 Set 方法：这是 flag.Value 接口的核心
// 它的作用是解析命令行传进来的字符串（如 "100C", "212F"）
func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64

	// 使用 Sscanf 解析数字和单位
	// 比如输入 "100C"，value=100, unit="C"
	fmt.Sscanf(s, "%f%s", &value, &unit)

	switch unit {
	case "C", "°C":
		// 如果是摄氏度，直接存入
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		// 如果是华氏度，先转成摄氏度再存入
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}

	return fmt.Errorf("无效的温度格式: %q (请使用类似 20C 或 100F 的格式)", s)
}

// String 方法：flag.Value 接口也需要这个
// 因为 celsiusFlag 内嵌了 Celsius，而 Celsius 已经有了 String() 方法
// 所以这里不需要重复写，Go 会自动使用 Celsius 的 String() 方法。
// func (f *celsiusFlag) String() string { return f.Celsius.String() }

// ---------------------------------------------------------
// 3. 封装帮助函数 (Glue Code)
// ---------------------------------------------------------

// CelsiusFlag 用于注册一个摄氏度的命令行参数
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	// a. 创建包装对象，填入默认值
	f := celsiusFlag{value}

	// b. 将这个包装对象注册给 flag 系统
	// 关键：f 实现了 flag.Value 接口，所以 flag 系统知道怎么处理它
	flag.CommandLine.Var(&f, name, usage)

	// c. 返回内部真正存储数据的指针
	return &f.Celsius
}

// ---------------------------------------------------------
// 4. 主函数 (使用演示)
// ---------------------------------------------------------

// 定义一个名为 "temp" 的命令行参数，默认值 20.0 度
var temp = CelsiusFlag("temp", 20.0, "设置温度 (例如: 20C, 212F)")

func main() {
	// 解析命令行参数
	flag.Parse()

	// 打印结果
	// 注意：*temp 拿到的是已经被 Set 方法转换过的摄氏度
	fmt.Printf("最终的温度是: %v\n", *temp)
}
