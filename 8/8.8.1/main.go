package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 限制并发数的信号量（令牌桶）
// 防止瞬间打开太多文件夹导致系统报错 "too many open files"
var sema = make(chan struct{}, 20)

func main() {
	// 1. 设置你要扫描的目标目录
	// 注意：在 Go 字符串中，反斜杠 \ 需要转义，所以写成 \\ 或者用 /
	root := `F:\新建文件夹`

	fmt.Printf("开始扫描目录: %s ...\n", root)
	startTime := time.Now()

	// 用于接收文件路径的 Channel
	// 为什么用 buffer？稍微缓冲一下，让发送方不至于因为接收方写文件慢了一点点就被卡住
	filePaths := make(chan string, 100)

	// 计数器，用于等待所有 Goroutine 结束
	var n sync.WaitGroup

	// 2. 启动第一个遍历任务
	n.Add(1)
	go walkDir(root, &n, filePaths)

	// 3. 启动一个后台协程，专门负责在任务完成后关闭 Channel
	// 这样主线程的 range 循环才能正常结束
	go func() {
		n.Wait()
		close(filePaths)
	}()

	// 4. 主线程：负责收集结果并写入文件
	outputFile := "file_list.txt"
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer f.Close()

	var count int
	// 不断从 Channel 里拿路径，直到 Channel 被关闭
	for path := range filePaths {
		count++
		// 写入文件，加换行符
		_, err := f.WriteString(path + "\n")
		if err != nil {
			fmt.Printf("写入出错: %v\n", err)
		}

		// 可选：每收集 1000 个打印一下进度，让你知道它活着
		if count%1000 == 0 {
			fmt.Printf("已找到 %d 个文件...\n", count)
		}
	}

	fmt.Printf("\n扫描完成！\n")
	fmt.Printf("耗时: %v\n", time.Since(startTime))
	fmt.Printf("总文件数: %d\n", count)
	fmt.Printf("结果已保存至: %s\n", outputFile)
}

// walkDir 递归并并发地遍历目录
func walkDir(dir string, n *sync.WaitGroup, filePaths chan<- string) {
	defer n.Done() // 这一层任务做完，计数器减1

	// 获取当前目录下的所有内容
	for _, entry := range dirents(dir) {
		fullPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// 如果是子目录，计数器加1，启动一个新的 Goroutine 去钻进去
			n.Add(1)
			go walkDir(fullPath, n, filePaths)
		} else {
			// 如果是文件，发送路径到 Channel
			filePaths <- fullPath
		}
	}
}

// dirents 读取目录内容（带限流保护）
func dirents(dir string) []os.FileInfo {
	// 1. 领令牌
	sema <- struct{}{}
	// 2. 确保函数退出时还令牌
	defer func() { <-sema }()

	// 3. 读取目录
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		// 遇到权限不足等错误，打印一下但不中断程序
		fmt.Fprintf(os.Stderr, "警告: 无法读取目录 %s: %v\n", dir, err)
		return nil
	}
	return entries
}
