package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// --- 1. 基础数据结构 ---
type Track struct {
	Title  string
	Artist string
	Year   int
	Length time.Duration
}

// 模拟数据
var tracks = []*Track{
	{"Go", "Delilah", 2012, length("3m38s")},
	{"Go", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", 2011, length("4m24s")},
}

// 辅助函数：把字符串转为时间
func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

// --- 2. 核心：有状态的排序器 ---

// 定义一个函数类型，代表“某一列的比较逻辑”
type LessFunc func(x, y *Track) bool

// 我们的多重排序器结构体
type customSort struct {
	t       []*Track   // 要排序的数据
	history []LessFunc // 点击历史（比较函数的栈）
}

// Len 和 Swap 是标准实现
func (x *customSort) Len() int      { return len(x.t) }
func (x *customSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

// Less 是核心：遍历历史记录
func (x *customSort) Less(i, j int) bool {
	// 依次取出排序规则（从最新的到最旧的）
	for _, less := range x.history {
		p, q := x.t[i], x.t[j]

		// 规则 A：如果 p < q，返回 true（确定 p 在 q 前）
		if less(p, q) {
			return true
		}
		// 规则 B：如果 p > q (即 q < p)，返回 false（确定 q 在 p 前）
		if less(q, p) {
			return false
		}

		// 规则 C：如果相等，继续比较下一个历史规则
	}
	return false // 所有规则都比完了还相等
}

// 模拟点击列头，把新的规则加到最前面
func (x *customSort) AddSortKey(f LessFunc) {
	// 将新的规则 f 放在切片的第一个位置，原来的往后移
	x.history = append([]LessFunc{f}, x.history...)
}

// --- 3. 打印表格的工具 ---
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\n"
	// 使用 tabwriter 让表格对齐
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Year, t.Length)
	}
	tw.Flush()
}

// --- 4. 主程序 ---
func main() {
	// 初始化排序器
	cs := &customSort{t: tracks}

	fmt.Println("======== 初始状态 (未排序) ========")
	printTracks(cs.t)

	// 1. 模拟点击 "Title"
	fmt.Println("\n======== 第1次点击: Title ========")
	cs.AddSortKey(func(x, y *Track) bool { return x.Title < y.Title })
	sort.Sort(cs)
	printTracks(cs.t)
	// 预期：按标题排，标题相同的(Go)，顺序未定

	// 2. 模拟点击 "Year"
	fmt.Println("\n======== 第2次点击: Year ========")
	// 现在的逻辑应该是：先比 Year，Year 一样才比 Title
	cs.AddSortKey(func(x, y *Track) bool { return x.Year < y.Year })
	sort.Sort(cs)
	printTracks(cs.t)

	// 3. 模拟点击 "Length"
	fmt.Println("\n======== 第3次点击: Length ========")
	// 现在的逻辑：Length > Year > Title
	cs.AddSortKey(func(x, y *Track) bool { return x.Length < y.Length })
	sort.Sort(cs)
	printTracks(cs.t)
}
