package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// 限制并发数的信号量（令牌桶）
var sema = make(chan struct{}, 20)

// 定义一个结构体来存储目录扫描结果 (用于最后展示)
type RootSize struct {
	Name string
	Size int64
}

// 定义一个结构体，用于在 Channel 中传递文件信息
type FileInfo struct {
	RootIndex int   // 属于第几个根目录
	Size      int64 // 文件大小
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	ticker := time.Tick(5 * time.Second)
	fmt.Printf("开始监控目录 %v ... (按 Ctrl+C 退出)\n", roots)

	printDiskUsage(roots)
	for range ticker {
		printDiskUsage(roots)
	}
}

func printDiskUsage(roots []string) {
	// 使用定义好的 FileInfo 类型
	fileSizes := make(chan FileInfo)
	var n sync.WaitGroup

	for i, root := range roots {
		n.Add(1)
		go walkDir(root, &n, i, fileSizes)
	}

	go func() {
		n.Wait()
		close(fileSizes)
	}()

	totals := make([]int64, len(roots))

	// 收集结果
	for info := range fileSizes {
		totals[info.RootIndex] += info.Size
	}

	// 清屏 (Windows CMD 可能不支持，可删)
	fmt.Print("\033[H\033[2J")
	fmt.Printf("--- Disk Usage (Update: %s) ---\n", time.Now().Format("15:04:05"))

	var results []RootSize
	for i, size := range totals {
		results = append(results, RootSize{Name: roots[i], Size: size})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})

	for _, r := range results {
		fmt.Printf("%-30s %.2f GB\n", r.Name, float64(r.Size)/1e9)
	}
}

// 注意这里 fileSizes 的类型变成了 chan<- FileInfo
func walkDir(dir string, n *sync.WaitGroup, rootIdx int, fileSizes chan<- FileInfo) {
	defer n.Done()

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, rootIdx, fileSizes)
		} else {
			// 发送 FileInfo 结构体
			// 使用 os.FileInfo 接口的 Size() 方法
			info, err := entry.Info()
			if err == nil {
				fileSizes <- FileInfo{RootIndex: rootIdx, Size: info.Size()}
			}
		}
	}
}

// 使用 os.ReadDir 替代 ioutil.ReadDir (Go 1.16+)
// 返回的是 []os.DirEntry 而不是 []os.FileInfo，这更高效
func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}
	defer func() { <-sema }()

	// os.ReadDir 是推荐的新方法，不需要引入 io/ioutil 包了
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	return entries
}
