package main

import (
	"fmt"
	"io" // ✅ 新版：替代 ioutil
	"net/http"
	"os"
	"strings"
)

func main() {
	plan4()
}

func plan1() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(os.Stderr, "fetch:%v\n", err)
			os.Exit(1)
		}
		b, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

func plan2() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		// 修改重点：
		// 使用 io.Copy(dst, src) 替代 io.ReadAll
		// dst (目标) 是 os.Stdout (标准输出)
		// src (来源) 是 resp.Body (HTTP 响应体)
		_, err = io.Copy(os.Stdout, resp.Body)

		// 记得关闭 Body
		resp.Body.Close()

		// 处理 io.Copy 可能返回的错误
		// 注意：这里复用了上面的 err 变量
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

func plan3() {
	for _, url := range os.Args[1:] {
		// 修改重点：检查前缀
		// 如果 url 不以 "http://" 开头，就手动给它加上
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

func plan4() {
	for _, url := range os.Args[1:] {
		// 修改逻辑：同时检查 http:// 和 https://
		// 只有当两个都没有的时候，才默认加上 http://
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		// --- 修改重点：打印 HTTP 状态码 ---
		// resp.Status 是一个字符串，例如 "200 OK" 或 "404 Not Found"
		fmt.Printf("HTTP Status: %s\n", resp.Status)
		// -------------------------------

		//_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		//if err != nil {
		//	fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		//	os.Exit(1)
		//}
	}
}
