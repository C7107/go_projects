package main

import (
	"fmt"
	"log"
	"net/http"
)

// 1. 定义数据库类型（本质是个 map）
type database map[string]int

// 2. 定义 list 方法：显示所有商品
func (db database) list(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		// Fprintf 把内容写给 response writer（也就是写回给浏览器）
		fmt.Fprintf(w, "%s: $%d\n", item, price)
	}
}

// 3. 定义 price 方法：查询具体商品价格
func (db database) price(w http.ResponseWriter, r *http.Request) {
	// 从 URL 参数中获取 item，比如 /price?item=socks
	item := r.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 返回 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "$%d\n", price)
}

func main() {
	// 初始化数据
	db := database{"shoes": 50, "socks": 5}
	// ---------------------------------------------------------
	// 重点在这里！
	// 我们没有创建 mux := http.NewServeMux()
	// 而是直接调用 http.HandleFunc。
	// 这会自动把路由注册到全局变量 DefaultServeMux 上。
	// ---------------------------------------------------------
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	fmt.Println("服务启动成功！请访问 http://localhost:8000/list")
	// ---------------------------------------------------------
	// 启动监听
	// 第二个参数传 nil。
	// 意思就是告诉 Server："我没给你传特定的 Handler，你自己去 DefaultServeMux 找吧"
	// ---------------------------------------------------------
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
