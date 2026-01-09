package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 启动一个 goroutine 专门负责：从服务器读数据 -> 显示到屏幕
	go func() {
		// io.Copy 会一直把 conn 的数据拷贝到 Stdout，直到连接断开
		io.Copy(os.Stdout, conn)
		log.Println("连接已断开")
		os.Exit(0)
	}()

	// 主程序负责：从键盘读数据 -> 发送给服务器
	// io.Copy 会一直把 Stdin 的数据拷贝到 conn
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
