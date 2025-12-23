// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	plan2()
}

func fetch1(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

func fetch2(urlStr string, ch chan<- string, round int) {
	start := time.Now()

	resp, err := http.Get(urlStr)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close() // 这里的 defer 写法更优雅，确保函数结束一定关闭

	// --- 修改部分：创建文件准备写入 ---

	// 1. 生成一个文件名：把 URL 中的特殊字符换掉，加上轮次后缀
	// 例如: www.baidu.com -> www_baidu_com_round1.html
	safeName := strings.ReplaceAll(urlStr, "https://", "")
	safeName = strings.ReplaceAll(safeName, "http://", "")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	safeName = url.QueryEscape(safeName) // 确保文件名安全
	if len(safeName) > 50 {
		safeName = safeName[:50]
	} // 防止文件名过长
	filename := fmt.Sprintf("%s_round%d.dat", safeName, round)

	// 2. 创建文件
	out, err := os.Create(filename)
	if err != nil {
		ch <- fmt.Sprintf("creating file %s: %v", filename, err)
		return
	}
	defer out.Close()

	// 3. 将网络流 (resp.Body) 拷贝到文件流 (out)
	// io.Copy 会同时完成读取网络数据和写入硬盘数据的工作
	nbytes, err := io.Copy(out, resp.Body)

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", urlStr, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("耗时: %.4fs  大小: %7d 字节  文件: %s  URL: %s", secs, nbytes, filename, urlStr)
}

func plan1() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch1(url, ch) // 告诉 Go 运行时（Runtime），不要在当前的执行流中同步等待 fetch 函数运行结束，而是立即启动一个新的执行流（Goroutine）去运行它
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func plan2() {
	// 如果用户没有提供参数，我们默认提供几个适合测试缓存的大型网站或静态资源
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://www.bilibili.com",                    // 动态内容多，但静态资源有缓存
			"https://code.jquery.com/jquery-3.7.1.min.js", // 纯静态资源（CDN），极其适合测试缓存
			"https://www.github.com",
		}
		fmt.Println("未提供参数，使用默认 URL 进行测试...")
	}

	ch := make(chan string)

	// 执行两轮测试
	for round := 1; round <= 2; round++ {
		fmt.Printf("\n========== 第 %d 轮请求 (Round %d) ==========\n", round, round)

		// 1. 并发启动这一轮的所有请求
		for _, u := range urls {
			go fetch2(u, ch, round)
		}

		// 2. 等待这一轮的所有请求完成
		// 必须在这里接收完所有结果，才能保证 Round 1 彻底结束，Round 2 还没开始
		for range urls {
			fmt.Println(<-ch)
		}
	}
}
