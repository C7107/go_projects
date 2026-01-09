// 这是一个“阻塞”的服务器，一次只能理一个人
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// 1. 开门营业：监听 8000 端口
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// 2. 坐等客户：Accept 会一直卡在这里，直到有人连接进来
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// 3. 服务客户：这里是关键问题所在！
		// 主程序亲自去服务这个客户，不服务完（客户不断开），
		// 主程序绝不回头去接下一个客。
		handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 每秒钟发一次时间给客户
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		} // 客户断开了才退出
		time.Sleep(1 * time.Second)
	}
}
