/*
练习 8.6： 为并发爬虫增加深度限制。也就是说，
如果用户设置了depth=3，那么只有从首页跳转三次以内能够跳到的页面才能被抓取到。
*/

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/C7107/go_projects/8/8.6/links" // 假设你已经有这个包，或者我会提供简易替代
)

// 定义一个令牌桶，限制并发数为 20
var tokens = make(chan struct{}, 20)

// 包装后的 crawl，增加令牌限制
func crawl(url string) []string {
	fmt.Println(url)     // 打印正在爬取的 URL
	tokens <- struct{}{} // 获取令牌
	list, err := links.Extract(url)
	<-tokens // 释放令牌
	if err != nil {
		log.Print(err)
	}
	return list
}

// 定义任务结构体：包含 URL 列表和当前的深度
type workItem struct {
	urlList []string
	depth   int
}

func main() {
	// 解析命令行参数
	// 比如: go run main.go -depth=2 http://gopl.io
	depthLimit := flag.Int("depth", 3, "depth limit") // 默认深度限制为 3
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"http://gopl.io"} // 默认抓取 gopl.io
	}
	worklist := make(chan workItem)
	var n int // 计数器：等待中的任务数

	// 1. 发送初始任务（深度为 0）
	n++
	go func() {
		worklist <- workItem{urlList: roots, depth: 0}
	}()
	// 记录已访问过的 URL
	seen := make(map[string]bool)
	// 2. 主循环
	for ; n > 0; n-- {
		item := <-worklist // 接收一个任务包

		// 如果当前这一批链接的深度已经达到限制，就不再继续生成新任务了
		// 注意：这里仍然会打印它们（因为是在上一步 crawl 打印的），
		// 但不会去爬它们内部的链接（不会产生 depth+1 的任务）。
		if item.depth >= *depthLimit {
			continue
		}
		for _, link := range item.urlList {
			if !seen[link] {
				seen[link] = true
				n++ // 增加计数器，因为我们即将启动一个新的 goroutine

				// 启动协程爬取，并将结果作为 depth + 1 发送回 worklist
				go func(link string, newDepth int) {
					foundLinks := crawl(link)
					worklist <- workItem{
						urlList: foundLinks,
						depth:   newDepth,
					}
				}(link, item.depth+1) // 传入下一层深度
			}
		}
	}
}
