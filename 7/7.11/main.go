package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// --- 1. 基础定义 ---

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// --- 2. 数据库结构体 (带锁) ---
// 我们不再直接把 database 定义为 map，而是定义为一个结构体
// 这样可以把 map 和 保护它的锁(Mutex) 绑在一起
type database struct {
	store map[string]dollars // 真正存数据的地方
	mu    sync.Mutex         // 互斥锁，用来保护 store
}

// --- 3. CRUD 处理函数 ---

// [R] List: 列出所有商品
func (db *database) list(w http.ResponseWriter, req *http.Request) {
	// 凡是涉及读写 map，都要先加锁
	db.mu.Lock()
	defer db.mu.Unlock() // 函数结束时自动解锁

	for item, price := range db.store {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

// [R] Price: 读取单个商品价格
func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.mu.Lock()
	defer db.mu.Unlock()

	price, ok := db.store[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

// [C] Create: 创建新商品
// URL: /create?item=hat&price=20
func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")

	// 简单的参数校验
	if item == "" || priceStr == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "item and price are required\n")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil || price < 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprint(w, "invalid price\n")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// 检查是否已经存在
	if _, ok := db.store[item]; ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "item %q already exists\n", item)
		return
	}

	db.store[item] = dollars(price)
	fmt.Fprintf(w, "created %s: %s\n", item, dollars(price))
}

// [U] Update: 更新商品价格
// URL: /update?item=socks&price=6
func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil || price < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid price\n")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// 检查是否存在
	if _, ok := db.store[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	db.store[item] = dollars(price)
	fmt.Fprintf(w, "updated %s: %s\n", item, dollars(price))
}

// [D] Delete: 删除商品
// URL: /delete?item=socks
func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.store[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	delete(db.store, item)
	fmt.Fprintf(w, "deleted %s\n", item)
}

// --- 4. 主程序 ---

func main() {
	// 初始化结构体
	db := &database{
		store: map[string]dollars{"shoes": 50, "socks": 5},
	}

	// 注册路由
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)

	fmt.Println("服务器运行在 http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
