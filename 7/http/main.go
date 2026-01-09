package main

import (
	"fmt"
	"log"
	"net/http"
)

// --- 1. 定义基础数据结构 ---

// 定义“美元”类型，为了能漂亮地打印出 $50.00 这种格式
type dollars float32

// 给 dollars 类型加个 String 方法，这样 fmt.Printf 的时候会自动调用它
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// 定义数据库类型，它本质上是一个 map
type database map[string]dollars

// --- 2. 定义处理逻辑（具体的业务函数） ---

// 处理 /list 请求：列出所有商品
// 注意：这是一个普通的方法，签名符合 func(w, r)
func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		// Fprintf 会把内容写入 w (也就是发回给浏览器)
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

// 处理 /price 请求：查询单个商品价格
// 例如：/price?item=socks
func (db database) price(w http.ResponseWriter, req *http.Request) {
	// 获取 URL 中的查询参数 "item"
	item := req.URL.Query().Get("item")

	// 在 map 中查找
	price, ok := db[item]

	// 如果没找到
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 设置状态码 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	// 找到了，打印价格
	fmt.Fprintf(w, "%s\n", price)
}

// --- 3. 主程序 ---

func main() {
	// 初始化库存数据
	db := database{"shoes": 50, "socks": 5}

	// 【核心代码】
	// 使用 http.HandleFunc 将处理函数注册到 Go 默认的全局路由表（DefaultServeMux）中。
	// 这里不需要 mux.Handle，也不需要把函数强转为 HandlerFunc，Go 会自动帮我们做适配。
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)

	fmt.Println("服务器已启动，请访问 http://localhost:8000/list")

	// 启动服务
	// 第二个参数是 nil，明确告诉 ListenAndServe：“请使用默认的 DefaultServeMux”
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
