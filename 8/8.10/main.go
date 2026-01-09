package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// client 类型定义为一个只写的字符串通道
// 这一头是往里写数据的，另一头（clientWriter）是读数据的
type client chan<- string

var (
	// 三个全局核心通道
	entering = make(chan client) // 新用户注册通道
	leaving  = make(chan client) // 用户离开注销通道
	messages = make(chan string) // 全局广播消息通道
)

func main() {
	// 1. 启动监听，端口 8000
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("聊天服务器启动，监听 localhost:8000 ...")

	// 2. 启动广播中心（后台管理 Goroutine）
	go broadcaster()

	// 3. 循环等待用户连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// 4. 为每一个连接进来的用户启动一个专属的处理 Goroutine
		go handleConn(conn)
	}
}

// broadcaster 是广播中心，它维护所有在线用户的集合
// 它是唯一能访问 clients map 的 Goroutine，所以不需要锁
func broadcaster() {
	clients := make(map[client]bool) // 记录当前所有在线的客户端通道

	for {
		// select 多路复用，监听三个通道的动静
		select {
		case msg := <-messages:
			// 情况A: 全局消息通道有新消息 -> 广播给所有人
			for cli := range clients {
				cli <- msg // 把消息塞入每个用户的专属通道
			}

		case cli := <-entering:
			// 情况B: 有新用户进来 -> 在名册上登记
			clients[cli] = true

		case cli := <-leaving:
			// 情况C: 有用户离开 -> 删除名单，并关闭他的通道
			delete(clients, cli)
			close(cli)
		}
	}
}

// handleConn 处理单个客户端的生命周期
func handleConn(conn net.Conn) {
	ch := make(chan string)   // 创建该用户的专属消息通道
	go clientWriter(conn, ch) // 启动子协程：专门负责把通道里的消息写回给客户端网络

	who := conn.RemoteAddr().String() // 获取用户 IP:Port
	ch <- "你是: " + who                // 欢迎语（只发给自己）
	messages <- who + " 来了"           // 广播通知所有人
	entering <- ch                    // 向广播中心登记自己

	// 循环读取客户端发送过来的每一行文本
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text() // 将用户说的话放入广播通道
	}
	// 注意：如果客户端断开连接，input.Scan() 会返回 false，循环结束

	leaving <- ch           // 向广播中心注销自己
	messages <- who + " 走了" // 广播通知所有人
	conn.Close()            // 关闭网络连接
}

// clientWriter 专门负责向客户端写数据
// 它遍历 channel，只要里面有数据，就通过网络发给用户
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: 这里忽略了网络写入错误
	}
}
